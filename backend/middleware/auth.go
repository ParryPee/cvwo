package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	JWTKey []byte
}

type Claims struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // Standard JWT fields like Expiry
}

type contextKey string

const UserIDKey contextKey = "UserID"

func (m *AuthMiddleware) parseUserClaims(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	var tokenString string

	if err == nil {
		tokenString = cookie.Value
	} else {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return nil, http.ErrNoCookie
		}
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return m.JWTKey, nil // Use the key stored in the struct
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.parseUserClaims(r)
		if err != nil {
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (m *AuthMiddleware) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.parseUserClaims(r)
		var userID int64 = 0
		if err == nil {
			userID = claims.UserID
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
