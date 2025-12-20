package models

import (
	"database/sql"
	"time"
)

type Topic struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`

	CreatedBy int64 `json:"created_by"`
	// Denormalized convenience field for the creator's username
	CreatedByUsername string `json:"created_by_username"`
}

type TopicDB struct {
	DB *sql.DB
}

func (m *TopicDB) All() ([]Topic, error) {
	rows, err := m.DB.Query(`
		SELECT t.id, t.title, t.description, t.created_at, t.created_by, u.username
		FROM topics t
		JOIN users u ON t.created_by = u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.CreatedBy, &t.CreatedByUsername); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}

	return topics, nil
}
func (m *TopicDB) GetByID(topicID int64) (*Topic, error) {
	row := m.DB.QueryRow(`
		SELECT t.id, t.title, t.description, t.created_at, t.created_by, u.username
		FROM topics t
		JOIN users u ON t.created_by = u.id
		WHERE t.id = ?`, topicID)
	var t Topic
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt, &t.CreatedBy, &t.CreatedByUsername); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
func (m *TopicDB) Create(title, description string, createdBy int64) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO topics (title, description, created_at, created_by) VALUES (?, ?, ?, ?)", title, description, time.Now(), createdBy)
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
