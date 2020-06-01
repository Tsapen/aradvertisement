package postgres

import (
	"github.com/Tsapen/aradvertisement/internal/ara"
)

// CreateUser creates user.
func (s *DB) CreateUser(user ara.UserCreationInfo) error {
	var q = `INSERT INTO users(username, password, email) VALUES ($1, $2, $3);`
	var _, err = s.Exec(q, user.Username, user.Password, user.Email)
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

// CheckLogin returns true if user exists.
func (s *DB) CheckLogin(ul ara.UserLogin) (bool, error) {
	var q = `SELECT 1 FROM users WHERE username = $1 AND password = $2;`
	var err error
	var dummy int
	err = s.QueryRow(q, ul.Username, ul.Password).Scan(&dummy)
	if err != nil {
		return false, err
	}

	return true, nil
}
