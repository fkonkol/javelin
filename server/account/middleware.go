package account

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
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

		next(w, r.WithContext(ctx))
	}
}
