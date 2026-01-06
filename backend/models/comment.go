package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID int64 `json:"id"`

	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	PostID int64 `json:"post_id"`

	UserID int64 `json:"user_id"`

	CreatedByUsername string `json:"created_by_username"`

	ParentCommentID sql.NullInt64 `json:"parent_comment_id,omitempty"`
	LikedByUser     bool          `json:"liked_by_user"`
	Deleted         bool          `json:"deleted"`
}

type CommentDB struct {
	DB *sql.DB
}

func (m *CommentDB) AllByPostID(postID, userID int64) ([]Comment, error) {
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
	result, err := m.DB.Exec("INSERT INTO comments (content, created_at, updated_at, post_id, user_id, parent_id) VALUES (?, ?, ?, ?, ?, ?)",
		content, time.Now(), time.Now(), postID, userID, parentCommentID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *CommentDB) Delete(commentID int64) error {
	_, err := m.DB.Exec("UPDATE comments SET deleted = 1 WHERE id = ?", commentID)
	return err
}

func (m *CommentDB) Update(commentID int64, content string) error {
	_, err := m.DB.Exec("UPDATE comments SET content = ?, updated_at = ? WHERE id = ?", content, time.Now(), commentID)
	return err
}
func (m *CommentDB) GetByParentID(commentID int64) (*[]Comment, error) { // Get all comments under a parent comment
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
func (m *CommentDB) GetByID(commentID int64) (*Comment, error) { // Get comment by ID
	row := m.DB.QueryRow(`SELECT c.id, c.content, c.likes, c.created_at, c.updated_at, c.post_id, c.user_id, c.parent_id,c.deleted, u.username 
	FROM comments c join users u on c.user_id = u.id WHERE c.id = ?`, commentID)

	var c Comment
	if err := row.Scan(&c.ID, &c.Content, &c.Likes, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID, &c.ParentCommentID, &c.Deleted, &c.CreatedByUsername); err != nil {
		return nil, err
	}
	return &c, nil

}

func (m *CommentDB) LikeComment(commentID, userID int64) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
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
	} else {
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
