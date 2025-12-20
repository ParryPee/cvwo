package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Content string `json:"content"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	TopicID int64 `json:"topic_id"`

	UserID int64 `json:"user_id"`

	CreatedByUsername string `json:"created_by_username"`
}

type PostDB struct {
	DB *sql.DB
}

func (m *PostDB) AllByTopicID(topicID int64) ([]Post, error) {
	rows, err := m.DB.Query("SELECT p.id, p.title, p.content, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username FROM posts p join users u on p.user_id = u.id WHERE p.topic_id = ?", topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
func (m *PostDB) Create(title, content string, topicID, userID int64) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO posts (title, content, created_at, updated_at, topic_id, user_id) VALUES (?, ?, ?, ?, ?, ?)", title, content, time.Now(), time.Now(), topicID, userID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *PostDB) Delete(postID int64) error {
	_, err := m.DB.Exec("DELETE FROM posts WHERE id = ?", postID)
	return err
}
func (m *PostDB) GetByID(postID int64) (*Post, error) {
	row := m.DB.QueryRow("SELECT p.id, p.title, p.content, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username FROM posts p join users u on p.user_id = u.id WHERE p.id = ?",
		postID)
	var p Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
func (m *PostDB) Update(postID int64, content string) error {
	_, err := m.DB.Exec("UPDATE posts SET content = ?, updated_at = ? WHERE id = ?", content, time.Now(), postID)
	return err
}
