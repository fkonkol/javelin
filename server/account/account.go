package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Time in seconds after which session times out.
const SESSION_TIME int = 5

// Time in seconds after which user has to log in manually again.
// Value set to 7 days.
const PERSIST_SESSION_TIME int = 604800

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

	query := "SELECT id, username, password FROM users WHERE email=$1"

	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		json.NewDecoder(r.Body).Decode(&request)

		// TODO: Change to DBQueryResponse struct
		var userID int
		var username string
		var hashedPassword string
		err := acc.db.QueryRow(context.Background(), query, request.Email).Scan(&userID, &username, &hashedPassword)
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

		bytes := make([]byte, 32)
		_, err = rand.Read(bytes)
		if err != nil {
			log.Printf("User login random bytes error: %v\n", err)
			return
		}

		sessionID := hex.EncodeToString(bytes)

		_, err = acc.sessions.SetNX(context.Background(), sessionID, userID, time.Duration(SESSION_TIME)*time.Second).Result()
		if err != nil {
			log.Printf("User login session store error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes = make([]byte, 32)
		_, err = rand.Read(bytes)
		if err != nil {
			log.Printf("User login persist random bytes error: %v\n", err)
			return
		}

		key := hex.EncodeToString(bytes)
		persistID := fmt.Sprintf("%s.%s", username, key)

		// Map username to buffer value
		_, err = acc.sessions.SetNX(context.Background(), key, username, time.Duration(PERSIST_SESSION_TIME)*time.Second).Result()
		if err != nil {
			log.Printf("User login persist store error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Cookie used for authentication and authorization
		http.SetCookie(w, &http.Cookie{
			Name:     "sid",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   SESSION_TIME,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode, // Frontend runs on different port, so we can't use strict mode
		})

		// Cookie used for login persistence
		http.SetCookie(w, &http.Cookie{
			Name:     "persist",
			Value:    persistID,
			Path:     "/",
			MaxAge:   PERSIST_SESSION_TIME,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode, // Frontend runs on different port, so we can't use strict mode
		})
	}
}

func (acc *AccountHandler) GetUserByUsername() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
	}

	type User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}

	type Response struct {
		Users []User `json:"users"`
	}

	query := `SELECT id, username FROM users WHERE username LIKE '%'||$1||'%'`

	return func(w http.ResponseWriter, r *http.Request) {
		var request Request
		request.Username = r.URL.Query().Get("username")

		rows, err := acc.db.Query(context.Background(), query, request.Username)
		if err != nil {
			log.Printf("Get user by username query error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer rows.Close()

		var response Response
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username)
			if err != nil {
				log.Printf("Rows scan error: %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			response.Users = append(response.Users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&response)
	}
}

func (acc *AccountHandler) SendFriendRequest() http.HandlerFunc {
	type FriendshipStatus string
	//Friendship status as declared in database schema
	const senderStatus FriendshipStatus = "PENDING_SENT"
	const receiverStatus FriendshipStatus = "PENDING_RECEIVED"

	query := `
		INSERT INTO friends (uuid, status, user_id, friend_id) 
		VALUES ($1, $2::text::friends_status, $3, $4)
	`

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(SessionKey).(int)

		friendParam := r.URL.Query().Get("friend_id")
		friendID, err := strconv.Atoi(friendParam)
		if err != nil {
			log.Printf("Send friend request invalid parameter: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		friendshipUUID := uuid.NewString()
		friendshipUUID = strings.Replace(friendshipUUID, "-", "", -1)

		// Create pending friendship between users
		_, err = acc.db.Exec(context.Background(), query, friendshipUUID, senderStatus, userID, friendID)
		if err != nil {
			log.Printf("Send friend request db query error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = acc.db.Exec(context.Background(), query, friendshipUUID, receiverStatus, friendID, userID)
		if err != nil {
			log.Printf("Send friend request db query error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (acc *AccountHandler) AcceptFriendRequest() http.HandlerFunc {
	query := `
		UPDATE friends
		SET status = 'ACCEPTED'
		WHERE 
		user_id = $1 AND 
		friend_id = $2 AND 
		status = 'PENDING_RECEIVED'
	`

	query2 := `
		UPDATE friends
		SET status = 'ACCEPTED'
		WHERE 
		user_id = $1 AND 
		friend_id = $2 AND 
		status = 'PENDING_SENT'
	`

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(SessionKey).(int)

		friendParam := r.URL.Query().Get("friend_id")
		friendID, err := strconv.Atoi(friendParam)
		if err != nil {
			log.Printf("Accept friend request invalid parameter: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = acc.db.Exec(context.Background(), query, userID, friendID)
		if err != nil {
			log.Printf("Accept friend request db exec error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = acc.db.Exec(context.Background(), query2, friendID, userID)
		if err != nil {
			log.Printf("Accept friend request db exec error 2: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
