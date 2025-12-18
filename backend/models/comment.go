package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID int64 `json:"id"`

	Content string `json:"content"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	PostID int64 `json:"post_id"`

	UserID int64 `json:"user_id"`
}

type CommentDB struct {
	DB *sql.DB
}

func (m *CommentDB) All() ([]Comment, error) {
	rows, err := m.DB.Query("SELECT id, content, created_at, updated_at, post_id, user_id FROM comments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (m *CommentDB) AllByPostID(postID int64) ([]Comment, error) {
	rows, err := m.DB.Query("SELECT id, content, created_at, updated_at, post_id, user_id FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (m *CommentDB) Create(postID int64, userID int64, content string) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO comments (content, created_at, updated_at, post_id, user_id) VALUES (?, ?, ?, ?, ?)",
		content, time.Now(), time.Now(), postID, userID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *CommentDB) Delete(commentID int64) error {
	_, err := m.DB.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}

func (m *CommentDB) Update(commentID int64, content string) error {
	_, err := m.DB.Exec("UPDATE comments SET content = ?, updated_at = ? WHERE id = ?", content, time.Now(), commentID)
	return err
}
func (m *CommentDB) GetByID(commentID int64) (*Comment, error) {
	row := m.DB.QueryRow("SELECT id, content, created_at, updated_at, post_id, user_id FROM comments WHERE id = ?", commentID)
	var c Comment
	if err := row.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.PostID, &c.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
