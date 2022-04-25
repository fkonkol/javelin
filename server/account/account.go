package account

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *AccountHandler {
	return &AccountHandler{db}
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
