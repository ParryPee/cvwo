package models

import (
	"database/sql"
	"time"
)

// Comment class, stores important attributes to send to the frontend
type Comment struct {
	ID int64 `json:"id"`

	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	PostID int64 `json:"post_id"`

	UserID  int64 `json:"user_id"`
	Deleted bool  `json:"deleted"`

	ParentCommentID sql.NullInt64 `json:"parent_comment_id,omitempty"`
	//Added field to indicate if the comment is liked by the user making the request
	LikedByUser bool `json:"liked_by_user"`
	//Username of the comment creator
	CreatedByUsername string `json:"created_by_username"`
}

// DB instance to make queries to
type CommentDB struct {
	DB *sql.DB
}

func (m *CommentDB) AllByPostID(postID, userID int64) ([]Comment, error) { //Gets all the comments under a certain post
	//Gets the respective comment columns, together with the username that matches the user id of the comment row
	//Also searches the comment_likes table for an entry where the both the user id and comment id match the row entry
	//This is returned in a separate boolean column liked_by_user
	query := `SELECT c.id, c.content, c.likes, c.created_at, c.updated_at,
		 c.post_id, c.user_id, c.parent_id,c.deleted, u.username, 
		 EXISTS (SELECT 1 FROM comment_likes cl where cl.comment_id = c.id AND cl.user_id = ?) AS liked_by_user
	
	FROM comments c join users u on c.user_id = u.id WHERE c.post_id = ? `
	rows, err := m.DB.Query(query, userID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	//Scans through the rows returned and checks that the returned fields matches the Comment class.
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.Likes, &c.CreatedAt, &c.UpdatedAt,
			&c.PostID, &c.UserID, &c.ParentCommentID, &c.Deleted, &c.CreatedByUsername, &c.LikedByUser); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (m *CommentDB) Create(postID int64, userID int64, content string, parentCommentID sql.NullInt64) (int64, error) {
	//Inserts a new comment
	result, err := m.DB.Exec("INSERT INTO comments (content, created_at, updated_at, post_id, user_id, parent_id) VALUES (?, ?, ?, ?, ?, ?)",
		content, time.Now().UTC(), time.Now().UTC(), postID, userID, parentCommentID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *CommentDB) Delete(commentID int64) error {
	//Sets the deleted flag in a comment to true(1)
	_, err := m.DB.Exec("UPDATE comments SET deleted = 1 WHERE id = ?", commentID)
	return err
}

func (m *CommentDB) Update(commentID int64, content string) error {
	//Updates a comment content
	_, err := m.DB.Exec("UPDATE comments SET content = ?, updated_at = ? WHERE id = ?", content, time.Now().UTC(), commentID)
	return err
}

// Get all comments under a parent comment, useful for sub-replies
func (m *CommentDB) GetByParentID(commentID int64) (*[]Comment, error) {
	rows, err := m.DB.Query("SELECT c.id, c.content, c.likes, c.created_at, c.updated_at, c.post_id, c.user_id, c.parent_id, u.username FROM comments c join users u on c.user_id = u.id WHERE c.parent_id = ?", commentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.Likes, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Deleted, &c.CreatedByUsername); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return &comments, nil

}

// Get comment by ID
func (m *CommentDB) GetByID(commentID int64) (*Comment, error) {
	row := m.DB.QueryRow(`SELECT c.id, c.content, c.likes, c.created_at, c.updated_at, c.post_id, c.user_id, c.parent_id,c.deleted, u.username 
	FROM comments c join users u on c.user_id = u.id WHERE c.id = ?`, commentID)

	var c Comment
	if err := row.Scan(&c.ID, &c.Content, &c.Likes, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Deleted, &c.CreatedByUsername); err != nil {
		return nil, err
	}
	return &c, nil

}

// Likes a comment
func (m *CommentDB) LikeComment(commentID, userID int64) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	//Checks if the comment is already liked, if it is liked, means that the user intends to unlike it, so delete it from the table.
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = ? AND user_id = ?)", commentID, userID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		_, err = tx.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("UPDATE comments SET likes = likes - 1 WHERE id = ?", commentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else { // If the like entry does NOT exist, it means the user intends to like the comment, so insert the like entry into the table.
		_, err = tx.Exec("INSERT INTO comment_likes (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
