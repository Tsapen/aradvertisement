package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
)

type dbMigration struct {
	number  int
	command string
}

func araDBMigrations() []dbMigration {
	return []dbMigration{
		dbMigration{
			number: 1,
			command: `
			CREATE TABLE IF NOT EXISTS users(
				id			SERIAL NOT NULL PRIMARY KEY
			,	username 	VARCHAR(100) NOT NULL UNIQUE
			,	email		VARCHAR(100) NOT NULL
			);`,
		},

		dbMigration{
			number: 2,
			command: `
			CREATE UNIQUE INDEX ON users(username) USING hash;
			`,
		},

		dbMigration{
			number: 3,
			command: `
			CREATE TABLE IF NOT EXISTS objects(
				id			SERIAL NOT NULL PRIMARY KEY
			,	user_id		INT NOT NULL REFERENCES users(id)
			,	latitude	VARCHAR(10) NOT NULL
			,	longitude	VARCHAR(10) NOT NULL
			,	comment		TEXT
			);`,
		},

		dbMigration{
			number: 4,
			command: `
			CREATE INDEX ON objects(latitude, longitude) USING hash;
			`,
		},
	}
}

func authDBMigrations() []dbMigration {
	return []dbMigration{
		dbMigration{
			number: 1,
			command: `
			CREATE TABLE IF NOT EXISTS tokens(
				uuid		VARCHAR(230) NOT NULL UNIQUE
			,	username 	VARCHAR(100) NOT NULL 
			, 	exp			BIGINT NOT NULL 
			);`,
		},

		dbMigration{
			number: 2,
			command: `
			CREATE UNIQUE INDEX ON tokens(uuid) USING hash;
			`,
		},

		dbMigration{
			number: 3,
			command: `
			CREATE TABLE users(
				username	VARCHAR(100) NOT NULL UNIQUE
			,	password	VARCHAR(100) NOT NULL
			);`,
		},

		dbMigration{
			number: 4,
			command: `
			CREATE UNIQUE INDEX ON users(username) USING hash;
			`,
		},
	}
}

func (db *DB) applied(num int) (bool, error) {
	var ex = 0
	var q = `SELECT 1 FROM migrations WHERE num = $1;`
	var err = db.QueryRow(q, num).Scan(&ex)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) apply(num int) error {
	var q = `INSERT INTO migrations(num) VALUES ($1);`
	var _, err = db.Exec(q, num)
	if err != nil {
		return err
	}

	return nil
}

// AraMigrate prepares ara db to work.
func (db *DB) AraMigrate() error {
	return db.migrate(araDBMigrations())
}

// AuthMigrate prepares auth db to work.
func (db *DB) AuthMigrate() error {
	return db.migrate(authDBMigrations())
}

func (db *DB) migrate(migs []dbMigration) error {
	var q = `CREATE TABLE IF NOT EXISTS migrations(
				num 		INT NOT NULL PRIMARY KEY,
				created_at	TIMESTAMP NOT NULL DEFAULT LOCALTIMESTAMP
			);`
	var _, err = db.Exec(q)
	if err != nil {
		return errors.Wrap(err, "can't create migrations table: ")
	}

	for _, mig := range migs {
		var migrated, err = db.applied(mig.number)
		if err != nil {
			return errors.Wrapf(err, "can't check %d migration: ", mig.number)
		}

		if migrated {
			continue
		}

		if _, err := db.Exec(mig.command); err != nil {
			errors.Wrapf(err, "can't apply %d migration: ", mig.number)
		}

		if err := db.apply(mig.number); err != nil {
			errors.Wrapf(err, "can't set %d migration as applied: ", mig.number)
		}
	}
	return nil
}
