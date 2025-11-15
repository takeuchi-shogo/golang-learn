package middleware

import (
	"log"
	"net/http"
)

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("JwtVerify", token)
		next.ServeHTTP(w, r)
	})
}
