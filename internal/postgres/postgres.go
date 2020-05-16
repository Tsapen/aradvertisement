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
		"postgres://%s:%s@%s:%s%s",
		c.UserName,
		c.Password,
		c.HostName,
		c.Port,
		c.VirtualHost)
}

// NewDBConnection create new storage.
func NewDBConnection(c *Config) (*DB, error) {
	var dbAddr = c.dbAddr()

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		return nil, errors.Wrap(err, "can't open connection")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
