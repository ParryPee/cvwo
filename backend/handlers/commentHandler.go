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
	vars := mux.Vars(r) //Get all variables in the URL
	postID := vars["post_id"]
	if postID == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}
	currentUserID, ok := getUserIDFromContext(r.Context()) // Get current user id, if unable, set current id to 0.
	if !ok {
		currentUserID = 0
	}
	CommentDB := models.CommentDB{DB: m.DB}
	//Convert string to integer
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post_id parameter", http.StatusBadRequest)
		return
	}
	//Get all comments under a post, this also returns a liked_by column that is boolean using the currentUserID.
	//Note: No user will have the user ID 0, so if it is 0 all comments returned should be false in the liked_by column.
	//Thus the reason for setting the userID to 0 is to handle unregistered users who are simply browsing the website.
	//Versus throwing and error preventing ALL unregistered users from using the application.
	comments, err := CommentDB.AllByPostID(postIDInt, currentUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching comments: %v", err), http.StatusInternalServerError)
		return
	}
	var comments_cleaned []models.Comment
	//Loop through all the returned comments, if the deleted column is True, set the content and Username to "deleted" and "redacted"
	//This ensures privacy and provides BACKEND censoring versus just censoring it in the frontend where
	// Malicious users might still be able to  Look at deleted content.
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
	var reqBody struct { //The response body we expect to receive
		PostID    int64  `json:"post_id"`
		UserID    int64  `json:"user_id"`
		Content   string `json:"content"`
		CreatedBy int64  `json:"created_by"`
		ParentID  *int64 `json:"parent_id,omitempty"`
	}
	//Decode the JSON body using the structure defined above, if a parameter is missing or doesn't match the type specified,
	//an error will be thrown.
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	//Check if the response body contains a parent ID, i.e. the user created a sub-reply.
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
	//Return the created commentID
	json.NewEncoder(w).Encode(map[string]int64{"id": commentID})
}
func (m *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get all variables in the URL
	commentID := vars["comment_id"]
	//Throw an error if missing the params
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	//Get the current user id from the context, if unable throw an authentication error.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	//Convert the commentID to an integer
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	//Fetch the commentID requested to be deleted, this will be used to compare the comment.userID to the current userID in the JWT.
	comment, err := CommentDB.GetByID(commentIDInt)
	if err != nil {
		http.Error(w, "Error fetching comment in delete", http.StatusInternalServerError)
		return
	}
	//Check if the user requesting to delete the comment is authorized to do so, if not throw a forbidden error.
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
	commentID := vars["comment_id"] // Get comment_id from the URL
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	// Convert commentID to integer
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment_id parameter", http.StatusBadRequest)
		return
	}
	var reqBody struct { //Request body we expect to receive
		Content string `json:"content"`
	}
	//Throw an error if the request we receive is not what we expected as per above.
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	CommentDB := models.CommentDB{DB: m.DB}
	//Get the current user id from the context. IF unable to do so or is empty, throw an authentication error.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	c, err := CommentDB.GetByID(commentIDInt) //Get the comment ID so that we can check if the creator of the comment is the requester.

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
	commentID := vars["comment_id"] //Get the comment ID from the URL
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	//Convert commentID to integer
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
	commentID := vars["comment_id"] //Get commentID from URL
	if commentID == "" {
		http.Error(w, "Missing comment_id parameter", http.StatusBadRequest)
		return
	}
	//Get user ID from the context, if unable throw an error.
	currentUserID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Auth error. Please ensure you are logged in.", http.StatusBadRequest)
		return
	}
	commentIDInt, err := strconv.ParseInt(commentID, 10, 64) // Convert CommentID to integer
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
