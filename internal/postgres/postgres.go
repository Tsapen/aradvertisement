package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Config contains settings for db
type Config struct {
	UserName    string
	Password    string
	HostName    string
	Port        string
	VirtualHost string
}

// DB contains db connection.
type DB struct {
	*sql.DB
}

func (c *Config) dbAddr() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.UserName, c.Password, c.HostName, c.Port, c.VirtualHost)
}

// NewDBConnection create new storage.
func NewDBConnection(c *Config) (*DB, error) {
	var err error
	var dbAddr = c.dbAddr()
	var db *sql.DB
	db, err = sql.Open("postgres", dbAddr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't open connection %s", dbAddr))
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't ping with connection %s", dbAddr))
	}

	return &DB{db}, nil
}
