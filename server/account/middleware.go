package account

import (
	"context"
	"fmt"
	"net/http"
)

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

		_, err = acc.sessions.Get(context.Background(), sessionCookie.Value).Result()
		if err != nil {
			fmt.Println("Session token not cached")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
