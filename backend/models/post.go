package models

import (
	"database/sql"
	"time"
)

type Post struct { //Post Class, holds important attributes in our Post object
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Content string `json:"content"`

	Likes int `json:"likes"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	TopicID    int64  `json:"topic_id"`
	TopicTitle string `json:"topic_title"`

	UserID int64 `json:"user_id"`

	CreatedByUsername string `json:"created_by_username"`
	LikedByUser       bool   `json:"liked_by_user"`
}

type PostDB struct {
	DB *sql.DB
}

func (m *PostDB) AllByTopicID(topicID int64) ([]Post, error) { //Selects all the posts under a specific topic
	rows, err := m.DB.Query("SELECT p.id, p.title, p.content, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username FROM posts p join users u on p.user_id = u.id WHERE p.topic_id = ?", topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() { //Ensure that the rows retrieved from the backend matches our Post object attributes
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
func (m *PostDB) Create(title, content string, topicID, userID int64) (int64, error) { //Creates a new Post
	result, err := m.DB.Exec("INSERT INTO posts (title, content, created_at, updated_at, topic_id, user_id) VALUES (?, ?, ?, ?, ?, ?)",
		title, content, time.Now().UTC(), time.Now().UTC(), topicID, userID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *PostDB) Delete(postID int64) error {
	_, err := m.DB.Exec("DELETE FROM posts WHERE id = ?", postID)
	return err
}

// Returns a Post by ID together with an additional column of whether the post is liked by the user
func (m *PostDB) GetByID(postID, userID int64) (*Post, error) {
	query := `SELECT p.id, p.title, p.content, p.likes, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username, t.title,
						EXISTS (SELECT 1 FROM post_likes pl where pl.post_id = p.id AND pl.user_id = ?) AS liked_by_user
	FROM posts p 
	JOIN users u ON p.user_id = u.id 
	JOIN topics t ON p.topic_id = t.id
	WHERE p.id = ?`
	row := m.DB.QueryRow(query, userID, postID)
	var p Post
	if err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Likes, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername, &p.TopicTitle, &p.LikedByUser); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Updates the post
func (m *PostDB) Update(postID int64, title, content string) error {
	_, err := m.DB.Exec("UPDATE posts SET title = ?, content = ?, updated_at = ? WHERE id = ?", title, content, time.Now().UTC(), postID)
	return err
}

// Like post
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
func (m *PostDB) SearchPost(query string) ([]Post, error) {
	sql_qry := `SELECT p.id, p.title, p.content, p.created_at, p.updated_at, p.topic_id, p.user_id,
	u.username FROM posts p 
	JOIN users u ON p.user_id = u.id
	WHERE MATCH(p.title,p.content) AGAINST (? IN NATURAL LANGUAGE MODE)
	ORDER BY p.created_at DESC
	`

	rows, err := m.DB.Query(sql_qry, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt,
			&p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
func (m *PostDB) GetAll(userID int64, limit int64, offset int64) ([]Post, error) {
	query := `SELECT p.id, p.title, p.content, p.likes, p.created_at, p.updated_at, p.topic_id, p.user_id, u.username, t.title,
						EXISTS (SELECT 1 FROM post_likes pl where pl.post_id = p.id AND pl.user_id = ?) AS liked_by_user
	FROM posts p 
	JOIN users u ON p.user_id = u.id
	JOIN topics t ON p.topic_id = t.id
	ORDER BY p.created_at DESC
	LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(query, userID, limit, offset)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Likes, &p.CreatedAt, &p.UpdatedAt, &p.TopicID, &p.UserID, &p.CreatedByUsername, &p.TopicTitle, &p.LikedByUser); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
