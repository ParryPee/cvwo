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

type CommentHandler struct {
	DB *sql.DB
}

func (m *CommentHandler) GetAllPostComments(w http.ResponseWriter, r *http.Request) {
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
	CommentDB := models.CommentDB{DB: m.DB}
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post_id parameter", http.StatusBadRequest)
		return
	}
	comments, err := CommentDB.AllByPostID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching comments: %v", err), http.StatusInternalServerError)
		return
	}
	var comments_cleaned []models.Comment
	for _, c := range comments {
		if c.Deleted {
			comments_cleaned = append(comments_cleaned,
				models.Comment{ID: c.ID, Content: "[Deleted]", Likes: 0, CreatedAt: c.CreatedAt,
					UpdatedAt: c.UpdatedAt, PostID: c.PostID,
					UserID: c.UserID, CreatedByUsername: "[Redacted]",
					ParentCommentID: c.ParentCommentID, LikedByUser: c.LikedByUser, Deleted: c.Deleted})
			continue
		}
		comments_cleaned = append(comments_cleaned, c)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments_cleaned)
}
func (m *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		PostID    int64  `json:"post_id"`
		UserID    int64  `json:"user_id"`
		Content   string `json:"content"`
		CreatedBy int64  `json:"created_by"`
		ParentID  *int64 `json:"parent_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	var parentCommentID sql.NullInt64
	if reqBody.ParentID != nil {
		parentCommentID = sql.NullInt64{Int64: *reqBody.ParentID, Valid: true}
	} else {
		parentCommentID = sql.NullInt64{Valid: false}
	}
	commentID, err := CommentDB.Create(reqBody.PostID, reqBody.UserID, reqBody.Content, parentCommentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating comment: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": commentID})
}
func (m *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		currentUserID = 0
	}
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	comment, err := CommentDB.GetByID(commentIDInt)
	if err != nil {
		http.Error(w, "Error fetching comment in delete", http.StatusInternalServerError)
		return
	}
	if comment.UserID != currentUserID {
		http.Error(w, "Not allowed to delete comments other than your own!", http.StatusForbidden)
		return
	}

	if err := CommentDB.Delete(commentIDInt); err != nil {
		http.Error(w, fmt.Sprintf("Error deleting comment: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	var reqBody struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	//Ensure that the user is authorized to edit this comment
	currentUserID, ok := getUserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	c, err := CommentDB.GetByID(commentIDInt)

	if err != nil {
		http.Error(w, "Error validating comment ownership", http.StatusInternalServerError)
		return
	}
	if c.UserID != currentUserID { // IF user is not authorized throw an error.
		http.Error(w, "You can only edit your own comments.", http.StatusUnauthorized)
	}

	if err := CommentDB.Update(commentIDInt, reqBody.Content); err != nil {
		http.Error(w, fmt.Sprintf("Error updating comment: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (m *CommentHandler) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	comment, err := CommentDB.GetByID(commentIDInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching comment: %v", err), http.StatusInternalServerError)
		return
	}
	if comment == nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
func (m *CommentHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	if err := CommentDB.LikeComment(commentIDInt, currentUserID); err != nil {
		http.Error(w, fmt.Sprintf("Error liking comment: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
