package postgres

import (
	"github.com/Tsapen/aradvertisement/internal/ara"
)

// CreateUser creates user.
func (s *DB) CreateUser(user ara.UserCreationInfo) error {
	var q = `INSERT INTO users(username, email) VALUES ($1, $2);`

	var _, err = s.Exec(q, user.Username, user.Email)
	return err
}

// DeleteUser deletes user.
func (s *DB) DeleteUser(username string) error {
	var q = `DELETE FROM objects 
				WHERE user_id = (SELECT id FROM users WHERE username = $1);`
	if _, err := s.Exec(q, username); err != nil {
		return err
	}

	q = `DELETE FROM users WHERE username = $1;`
	var _, err = s.Exec(q, username)
	return err
}
