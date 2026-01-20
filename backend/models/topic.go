package models

import (
	"database/sql"
	"time"
)

type Topic struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	UserID int64 `json:"created_by"`

	//Additional fields
	CreatedByUsername string `json:"created_by_username"`
	PostCount         int64  `json:"post_count"`
}

type TopicDB struct {
	DB *sql.DB
}

func (m *TopicDB) All() ([]Topic, error) {
	rows, err := m.DB.Query(`
		SELECT t.id, t.title, t.description, t.created_at, t.user_id, u.username, COUNT(p.id) as post_count
		FROM topics t
		JOIN users u ON t.user_id = u.id
		LEFT JOIN posts p ON t.id = p.topic_id
		GROUP BY t.id, t.title, t.description, t.created_at, t.user_id, u.username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserID, &t.CreatedByUsername, &t.PostCount); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}
func (m *TopicDB) GetByID(topicID int64) (*Topic, error) {
	row := m.DB.QueryRow(`
		SELECT t.id, t.title, t.description, t.created_at, t.user_id, u.username
		FROM topics t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = ?`, topicID)
	var t Topic
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserID, &t.CreatedByUsername); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
func (m *TopicDB) Create(title, description string, createdBy int64) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO topics (title, description, created_at, user_id) VALUES (?, ?, ?, ?)",
		title, description, time.Now().UTC(), createdBy)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *TopicDB) Delete(topicID int64) error {
	_, err := m.DB.Exec("DELETE FROM topics WHERE id = ?", topicID)
	return err
}
func (m *TopicDB) Update(topicID int64, title, description string) error {
	_, err := m.DB.Exec("UPDATE topics SET title = ?, description = ? WHERE id = ?", title, description, topicID)
	return err
}
func (m *TopicDB) GetByBatch(batch_size, offset int) ([]Topic, error) {
	rows, err := m.DB.Query(`SELECT t.id, t.title, t.description, t.created_at, t.user_id, u.username
		FROM topics t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
		LIMIT  ? OFFSET  ?`, batch_size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserID, &t.CreatedByUsername); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	if topics == nil {
		topics = []Topic{}
	}
	return topics, nil
}
func (m *TopicDB) SearchTopic(query string) ([]Topic, error) {
	sql_qry := `SELECT t.id, t.title, t.description,  t.created_at,t.user_id, u.username
	FROM topics t
	JOIN users u ON t.user_id = u.id
	WHERE t.title LIKE ? OR t.description LIKE ?
	ORDER BY t.created_at DESC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.Query(sql_qry, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.UserID, &t.CreatedByUsername); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	if topics == nil {
		topics = []Topic{}
	}
	return topics, nil
}
