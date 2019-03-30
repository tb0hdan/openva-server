package auth

import (
	"log"
	"net/http"
)

func VerifyToken(token string) bool {
	log.Println("Token: ", token)
	return true
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		token_values := values["token"]
		tokenValid := false
		if len(token_values) > 0 {
			token := token_values[0]
			tokenValid = VerifyToken(token)
		}
		if !tokenValid {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
