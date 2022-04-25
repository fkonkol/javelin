package account

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	db       *pgxpool.Pool
	sessions *redis.Client
}

func NewHandler(db *pgxpool.Pool, sessions *redis.Client) *AccountHandler {
	return &AccountHandler{db, sessions}
}

// Handles user registration. Validates values given in request body,
// then hashes password and inserts user into database.
// TODO: Email code verification
func (acc *AccountHandler) Register() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,alphanum"`
		Password string `json:"password" validate:"required,min=8"`
	}

	query := "INSERT INTO users(email, username, password) VALUES($1, $2, $3)"

	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request body
		var request Request
		json.NewDecoder(r.Body).Decode(&request)

		// Basic input validation
		err := validator.New().Struct(request)
		if err != nil {
			log.Printf("User register input validation error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Hash password with bcrypt
		hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error during password hashing: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Insert user into database
		_, err = acc.db.Exec(context.Background(), query, request.Email, request.Username, hash)
		if err != nil {
			log.Printf("User registration database query error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Registered user with email %s\n", request.Email)

		w.WriteHeader(http.StatusCreated)
	}
}

func (acc *AccountHandler) Login() http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	query := "SELECT password FROM users WHERE email=$1"

	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		json.NewDecoder(r.Body).Decode(&request)

		var hashedPassword string
		err := acc.db.QueryRow(context.Background(), query, request.Email).Scan(&hashedPassword)
		if err != nil {
			log.Printf("User login query row error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if given password and hashed password are equal
		// And therefore given credentials are valid
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(request.Password))
		if err != nil {
			log.Printf("User login bcrypt compare error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionToken := uuid.NewString()

		_, err = acc.sessions.SetNX(context.Background(), sessionToken, request.Email, 60*time.Second).Result()
		if err != nil {
			log.Printf("User login session store error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sid",
			Value:    sessionToken,
			MaxAge:   60,
			HttpOnly: true,
		})
	}
}
