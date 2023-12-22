package apiserver

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

func rateLimiterMiddleware(h http.Handler) http.Handler {
	limiter := rate.NewLimiter(5, 10)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func httpLoggingMiddleware(l *slog.Logger, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method

		l.Info("http request", "method", method, "uri", uri)

	}
	return http.HandlerFunc(fn)
}

func jwtAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == apiPrefix+"/generate-jwt" {
			h.ServeHTTP(w, r)
			return
		}
		if len(r.URL.Path) > 7 && r.URL.Path[:8] == "/swagger" {
			h.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// @Summary Generate JWT
// @Description Generating JWT Token for API Authorization.
// @Tags JWT
// @Accept json
// @Produce json
// @Success 200 {object} string "Token Generating Successfully."
// @Router /generate-jwt [get]
func generateJWT(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Error while signing the token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"` + tokenString + `"}`))
}
