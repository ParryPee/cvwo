package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Content string `json:"content"`

	Likes int `json:"likes"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	TopicID int64 `json:"topic_id"`

	UserID int64 `json:"user_id"`

	CreatedByUsername string `json:"created_by_username"`
	LikedByUser       bool   `json:"liked_by_user"`
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
func (m *PostDB) GetByID(postID, userID int64) (*Post, error) {
	query := `SELECT p.id, p.title, p.content, p.likes, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username, 
						EXISTS (SELECT 1 FROM post_likes pl where pl.post_id = p.id AND pl.user_id = ?) AS liked_by_user
	FROM posts p join users u on p.user_id = u.id WHERE p.id = ?`
	row := m.DB.QueryRow(query, userID, postID)
	var p Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Likes, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername, &p.LikedByUser); err != nil {
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
func (m *PostDB) LikePost(postID, userID int64) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	var exists bool

	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = ? AND user_id = ?)", postID, userID).Scan(&exists)
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists {
		_, err := tx.Exec("DELETE FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err := tx.Exec("INSERT INTO post_likes (post_id,user_id) VALUES (?,?)", postID, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()

}
