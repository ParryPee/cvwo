package handlers

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostHandler struct {
	DB *sql.DB
}

func (m *PostHandler) GetAllTopicPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID := vars["topic_id"]
	if topicID == "" {
		http.Error(w, "Missing topic_id parameter", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}
	topicIDInt, err := strconv.ParseInt(topicID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid topic_id parameter", http.StatusBadRequest)
		return
	}
	posts, err := PostDB.AllByTopicID(topicIDInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching posts: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (m *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		TopicID int64  `json:"topic_id"`
		UserID  int64  `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Title == "" || reqBody.Content == "" || reqBody.TopicID == 0 || reqBody.UserID == 0 {
		http.Error(w, "Title, Content, TopicID, and UserID are required", http.StatusBadRequest)
		return
	}

	PostDB := models.PostDB{DB: m.DB}

	postID, err := PostDB.Create(reqBody.Title, reqBody.Content, reqBody.TopicID, reqBody.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"post_id": postID})
}
func (m *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		currentUserID = 0
	}
	PostDB := models.PostDB{DB: m.DB}
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post_id parameter", http.StatusBadRequest)
		return
	}
	post, err := PostDB.GetByID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching post: %v", err), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
func (m *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get variables from the request
	postID := vars["post_id"]

	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	PostDB := models.PostDB{DB: m.DB}

	postIDInt, err := strconv.ParseInt(postID, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	post, err := PostDB.GetByID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching post: %v", err), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if post.UserID != currentUserID {
		http.Error(w, "Forbidden: You can only delete your own posts", http.StatusForbidden)
		return
	}
	res := PostDB.Delete(postIDInt)

	if res != nil {
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	var reqBody struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if reqBody.Content == "" || reqBody.Title == "" {
		http.Error(w, "Content and Title is required", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	PostDB := models.PostDB{DB: m.DB}

	postIDInt, err := strconv.ParseInt(postID, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := PostDB.GetByID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching post: %v", err), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	if post.UserID != currentUserID {
		http.Error(w, "Forbidden: You can only update your own posts", http.StatusForbidden)
		return
	}
	res := PostDB.Update(postIDInt, reqBody.Title, reqBody.Content)
	if res != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}
	if err := PostDB.LikePost(postIDInt, currentUserID); err != nil {
		http.Error(w, fmt.Sprintf("Error liking comment: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *PostHandler) SearchPost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}
	posts, err := PostDB.SearchPost(query)
	if err != nil {
		http.Error(w, "Error searching for posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
