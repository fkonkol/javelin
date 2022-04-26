package account

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ContextKey string

const SessionKey ContextKey = "session"

func (acc *AccountHandler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("sid")
		if err != nil {
			fmt.Println("Cookie not present in request")

			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
				SameSite: http.SameSiteStrictMode,
			})
		}

		next(w, r.WithContext(ctx))
	}
}
