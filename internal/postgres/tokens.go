package postgres

import (
	"log"
	"time"

	"github.com/Tsapen/aradvertisement/internal/auth"
)

// RunCleaning starts every minute deleting out-of-date tokens.
func (db *DB) RunCleaning() error {
	var q = "DELETE FROM tokens WHERE localtimestamp >= exp"
	var stmtCleaning, err = db.Prepare(q)
	if err != nil {
		return err
	}

	go func() {
		var interval = time.Minute
		for {
			time.Sleep(interval)
			if _, err = stmtCleaning.Exec(); err != nil {
				log.Printf("auth db: cleaning tokens error: %s\n", err)
			}
		}
	}()

	return nil
}

// CreateAuth creates authentiction in db.
func (db *DB) CreateAuth(username string, td *auth.TokenDetails) error {
	var q = `INSERT INTO tokens (uuid, username, exp) VALUES ($1, $2, $3);`

	if _, err := db.Exec(q, td.AccessUUID, username, td.AtExpires); err != nil {
		return err
	}

	if _, err := db.Exec(q, td.RefreshUUID, username, td.RtExpires); err != nil {
		return err
	}

	return nil
}

// DeleteAuth deletes tokens.
func (db *DB) DeleteAuth(givenUUID string) error {
	var q = `DELETE FROM tokens WHERE uuid = $1;`
	if _, err := db.Exec(q, givenUUID); err != nil {
		return err
	}

	return nil
}

// FetchAuth returns username.
func (db *DB) FetchAuth(authD *auth.AccessDetails) (string, error) {
	var q = `SELECT username FROM tokens WHERE uuid = $1`
	var username string
	if err := db.QueryRow(q, authD.AccessUUID).Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}

// CheckLogin returns username.
func (db *DB) CheckLogin(user auth.User) (bool, error) {
	var q = `SELECT TRUE FROM users WHERE username = $1 AND password = $2`
	var existsUser bool
	if err := db.QueryRow(q, user.Username, user.Password).Scan(&existsUser); err != nil {
		return false, err
	}

	if !existsUser {
		return false, nil
	}

	return true, nil
}

// InsertUser creates user in db.
func (db *DB) InsertUser(user auth.User) error {
	var q = `INSERT INTO users (username, password) VALUES ($1, $2);`

	if _, err := db.Exec(q, user.Username, user.Password); err != nil {
		return err
	}

	return nil
}
