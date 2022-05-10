package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContextKey string

const SessionKey ContextKey = "session"

func (acc *AccountHandler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		sessionCookie, err := r.Cookie("sid")
		var ctx context.Context
		persist := false
		if err != nil {
			log.Println("Session cookie not present in request")

			if err != http.ErrNoCookie {
				log.Printf("Auth middleware error: %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// If session cookie is not present, but persist cookie is, we can
			// validate the latter in cache and renew the session accordingly.

			persistCookie, err := r.Cookie("persist")
			if err != nil {
				log.Printf("Persist cookie error: %v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			split := strings.Split(persistCookie.Value, ".")
			username := split[0]
			key := split[1]

			// If error potential cookie hijacking - TODO: Log out user and notify him by email
			cached, err := acc.sessions.Get(context.Background(), key).Result()
			if err != nil {
				log.Printf("Auth middleware persistence cache error: %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Potential cookie hijacking
			if username != cached {
				log.Printf("Auth middleware persistence usernames do not match")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			var uid int
			query := "SELECT id FROM users WHERE username=$1"
			err = acc.db.QueryRow(context.Background(), query, username).Scan(&uid)
			if err != nil {
				log.Printf("Auth middleware db query error: %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			buffer := make([]byte, 32)
			_, err = rand.Read(buffer)
			if err != nil {
				log.Printf("Auth middleware buffer read error: %v\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sessionID := hex.EncodeToString(buffer)

			_, err = acc.sessions.SetNX(context.Background(), sessionID, uid, time.Duration(SESSION_TIME)*time.Second).Result()
			if err != nil {
				log.Printf("Auth middleware session store set error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "sid",
				Value:    sessionID,
				Path:     "/",
				MaxAge:   SESSION_TIME,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})

			ctx = context.WithValue(r.Context(), SessionKey, uid)
			persist = true
			log.Println("Cookie persisted")
		}

		if !persist {
			res, err := acc.sessions.Get(context.Background(), sessionCookie.Value).Result()
			if err != nil {
				fmt.Println("Session token not cached")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			uid, err := strconv.Atoi(res)
			if err != nil {
				fmt.Printf("Auth middleware string to int error: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), SessionKey, uid)

			// Sliding expiration: reset the expiration time for a valid session cookie
			// if a request is made and more than half of the timeout interval has elapsed.
			// TODO: Allow constant session cookie slide without it expiring for a limited time
			ttl, err := acc.sessions.TTL(context.Background(), sessionCookie.Value).Result()
			if err != nil {
				return
			}

			// Check if more than half of session timeout interval has elapsed
			halfTimeoutElapsed := ttl.Seconds()+float64(SESSION_TIME/2) < float64(SESSION_TIME)

			if halfTimeoutElapsed {
				// Reset the expiration time
				_, err = acc.sessions.Expire(context.Background(), sessionCookie.Value, time.Duration(SESSION_TIME)*time.Second).Result()
				if err != nil {
					log.Printf("Auth middleware renew session error: %v\n", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "sid",
					Value:    sessionCookie.Value,
					Path:     "/",
					MaxAge:   SESSION_TIME,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
				})
			}

			next(w, r.WithContext(ctx))
		} else {
			next(w, r.WithContext(ctx))
		}
	}
}
