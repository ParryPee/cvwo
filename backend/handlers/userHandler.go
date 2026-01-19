package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// Our user handler class is a little different and takes in the JWT key, this is the secret key used to sign and verify the
// legitimacy of the tokens.
type UserHandler struct {
	DB     *sql.DB
	JWTKey []byte
}

// This stores the information that we want to keep in our JWT
type Claims struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // Standard JWT fields like Expiry
}

func containsWhitespace(s string) bool {
	return strings.IndexFunc(s, unicode.IsSpace) != -1
}
func (m *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody struct { //Request body we expect to receive
		Username string `json:"username"`
	}
	//Validate the request body
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	//check username for white spaces and length
	if len(reqBody.Username) > 15 || len(reqBody.Username) < 7 {
		http.Error(w, "Username must be between 7 and 15 characters", http.StatusBadRequest)
		return
	}
	if containsWhitespace(reqBody.Username) {
		http.Error(w, "Username cannot contain whitespace", http.StatusBadRequest)
		return
	}
	UserDB := models.UserDB{DB: m.DB}
	row, err := UserDB.DB.Query("SELECT id FROM users WHERE username = ?", reqBody.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking existing user: %v", err), http.StatusInternalServerError)
		return
	}
	if row.Next() { // Check if user already exists
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	userID, err := UserDB.Create(reqBody.Username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//Return the userID of the newly created user
	json.NewEncoder(w).Encode(map[string]int64{"id": userID})
}

// This function isn't actually implemented as there is no account deletion feature as of yet
func (m *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var reqBody struct { //Request body we expect to receive
		UserID int64 `json:"user_id"`
	}
	// Validate the request body
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//TODO: Add authorization if this feature is intended
	UserDB := models.UserDB{DB: m.DB}
	if err := UserDB.Delete(reqBody.UserID); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqBody struct { //Request body we expect to receive
		UserName string `json:"username"`
	}
	//validate request body
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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
	//Check if user exists in the database
	if user == nil {
		http.Error(w, "User not found, please register", http.StatusNotFound)
		return
	} else {
		userID = user.ID
	}
	//Create a DateTime object for 24 hours from the login time.
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours

	// Create the Claims to be stored inside the cookie, in this case its just the
	// user id and the expiration time
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //Creates the token with the specified claims above
	tokenString, err := token.SignedString(m.JWTKey)           //Signs the token with the secret passphrase stored in the .env
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//sends the cookie back to the client
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) { // Returns the user object from the userid in the context
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
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
func (m *UserHandler) Logout(w http.ResponseWriter, r *http.Request) { //Logout function
	//MaxAge and Epires both "tell" the browser to delete the cookie but just added both for completeness.
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Logged out successfully"))
}
func (m *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) { // Returns the user object by ID
	vars := mux.Vars(r)
	userIDParam := vars["user_id"]
	var uid int64
	uid, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	UserDB := models.UserDB{DB: m.DB}
	user, err := UserDB.GetByID(uid)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
