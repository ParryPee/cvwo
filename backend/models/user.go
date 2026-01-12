package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int64 `json:"id"`

	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserDB struct {
	DB *sql.DB
}

func (m *UserDB) All() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, username, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User

		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
func (m *UserDB) Create(username string) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO users (username, created_at) VALUES (?, ?)", username, time.Now().UTC())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func (m *UserDB) Delete(userID int64) error {
	_, err := m.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}
func (m *UserDB) GetByID(userID int64) (*User, error) {
	row := m.DB.QueryRow("SELECT * FROM users WHERE id = ?", userID)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
func (m *UserDB) GetByUsername(username string) (*User, error) {
	row := m.DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
