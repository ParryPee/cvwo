package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct { // The AuthMiddleware "class" takes in the jwt key as it needs to verify authentication tokens
	JWTKey []byte
}

type Claims struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // Standard JWT fields like Expiry
}

type contextKey string

const UserIDKey contextKey = "UserID"

func (m *AuthMiddleware) parseUserClaims(r *http.Request) (*Claims, error) { //Takes in a cookie and returns the claims inside it
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
	claims := &Claims{} //Create the claims structure
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return m.JWTKey, nil // Tells the ParseWithClaims function to use the JWTkey passed to the middllware
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	//Checks the signing method, prevents Key Confusion attack
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.NoneSignatureTypeDisallowedError
	}
	return claims, nil
}

func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler { //validates the token and returns the userID in the context
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

// Similar to the Validate token function except it doesnt throw an error,but returns -1 (an invalid or non existent ID)
func (m *AuthMiddleware) OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := m.parseUserClaims(r)
		var userID int64 = -1
		if err == nil {
			userID = claims.UserID
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
