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
	vars := mux.Vars(r)         // Get all variables in the request
	topicID := vars["topic_id"] // Get the topic Id from the URL/request
	if topicID == "" {
		http.Error(w, "Missing topic_id parameter", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}
	//Convert topic ID into integer
	topicIDInt, err := strconv.ParseInt(topicID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid topic_id parameter", http.StatusBadRequest)
		return
	}
	posts, err := PostDB.AllByTopicID(topicIDInt) // Returns all the posts within a topic
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching posts: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (m *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody struct { // The request body we expect to receive
		TopicID int64  `json:"topic_id"`
		UserID  int64  `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	//Decode the request body that is in JSON,, if the request body is not what we expected throw an error
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Second check to ensure that the request body has all the necessary required fields.
	if reqBody.Title == "" || reqBody.Content == "" {
		http.Error(w, "Title, Content are required", http.StatusBadRequest)
		return
	}

	PostDB := models.PostDB{DB: m.DB}

	postID, err := PostDB.Create(reqBody.Title, reqBody.Content, reqBody.TopicID, reqBody.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating post: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	//Return the created postID
	json.NewEncoder(w).Encode(map[string]int64{"post_id": postID})
}
func (m *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)       //Get variables from the request
	postID := vars["post_id"] // Get post ID from the request
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	//Gets the current userID from the context, if unable or the user is not logged in returns 0.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		currentUserID = 0
	}
	PostDB := models.PostDB{DB: m.DB}
	//converts PostID to integer
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post_id parameter", http.StatusBadRequest)
		return
	}
	//Get the specified post by ID together with a boolean column "liked_by_user" this column will help to determine if the post is
	// liked by the user. It is important to note that no user will ever have the user id of 0, so if the requester is a visitor
	//without an account, this function willl still work even without a user ID.
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
	vars := mux.Vars(r)
	postID := vars["post_id"] //Get post ID from params

	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	//Get current user id, if not throw an error.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	PostDB := models.PostDB{DB: m.DB}

	postIDInt, err := strconv.ParseInt(postID, 10, 64)

	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	//Get the Post by ID.
	post, err := PostDB.GetByID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching post: %v", err), http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	//Verify that the requester is authorized to delete the post.
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
	var reqBody struct { // Request body that we expect to receive
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
	//Get user ID and throw authentication error if unable to get user id.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
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
	//Verify that the creator of the post is the same as the current user
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
	//Get user from context and throw authentication error if unable
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
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
func (m *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	PostDB := models.PostDB{DB: m.DB}

	size_str := r.URL.Query().Get("size")
	offset_str := r.URL.Query().Get("offset")
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		currentUserID = 0
	}

	var posts []models.Post

	var err error

	if size_str != "" && offset_str != "" {
		size, _ := strconv.ParseInt(size_str, 10, 64)
		offset, _ := strconv.ParseInt(offset_str, 10, 64)
		posts, err = PostDB.GetAll(currentUserID, size, offset)
	} else {
		// If no size or offset is specified just return the first 10
		posts, err = PostDB.GetAll(currentUserID, 10, 0)
	}
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}
