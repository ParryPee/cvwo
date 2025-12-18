package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	DB     *sql.DB
	JWTKey []byte
}

type Claims struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // Standard JWT fields like Expiry
}

func (m *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	UserDB := models.UserDB{DB: m.DB}
	userID, err := UserDB.Create(reqBody.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": userID})
}
func (m *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	UserDB := models.UserDB{DB: m.DB}
	if err := UserDB.Delete(reqBody.UserID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserName string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username is provided
	if reqBody.UserName == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	UserDB := models.UserDB{DB: m.DB}
	user, err := UserDB.GetByUsername(reqBody.UserName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	var userID int64
	if user == nil {
		http.Error(w, "User not found, please register", http.StatusNotFound)
		return
	} else {
		userID = user.ID
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})
}
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		// Should catch this just in case middleware fails
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDB := models.UserDB{DB: h.DB}
	user, err := userDB.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
func (m *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Logged out successfully"))
}
