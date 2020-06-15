package postgres

import (
	"database/sql"

	"github.com/Tsapen/aradvertisement/internal/ara"
)

// CreateObject puts new glTF object path in db.
func (s *DB) CreateObject(obj ara.ObjectCreationInfo) (int, error) {
	var q = `
	INSERT INTO objects(user_id, latitude, longitude, comment, type) VALUES ((
		SELECT id FROM users WHERE users.username = $1), 
		$2, $3, $4, $5) RETURNING id;`

	var id int
	var err = s.QueryRow(q, obj.Username, obj.Latitude, obj.Longitude, obj.Comment, obj.Type).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SelectObjectsAround selects glTF object paths from db.
func (s *DB) SelectObjectsAround(params ara.ObjectSelectInfo) (res []ara.ObjectAroundResp, err error) {
	// 0.00045 degrees gps are approximately equal to 50 meters
	var q = `SELECT u.username, o.latitude, o.longitude FROM objects o
				JOIN users u ON u.id = o.user_id
				WHERE |/((o.latutude-$1)^2 + (o.longitude-$2)^2) <= 0.00045;`

	var rows *sql.Rows
	rows, err = s.Query(q, params.Latitude, params.Longitude)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = ara.HandleErrPair(rows.Close(), err)
	}()

	for rows.Next() {
		var obj ara.ObjectAroundResp

		if err = rows.Scan(&obj.Username, &obj.Latitude, &obj.Longitude); err != nil {
			return nil, err
		}

		res = append(res, obj)
	}

	return res, nil
}

// SelectObjectByID selects object by id.
func (s *DB) SelectObjectByID(id int) (ara.ObjectSelectByID, error) {
	var q = `SELECT o.type, u.username 
			FROM objects o JOIN users u ON o.user_id=u.id 
			WHERE o.id = $1
			`
	var objType, username = "", ""

	if err := s.QueryRow(q, id).Scan(&objType, &username); err != nil {
		return ara.ObjectSelectByID{}, err
	}

	return ara.ObjectSelectByID{Type: rune(objType[0]), Username: username}, nil
}

// SelectUsersObjects selects objects of user.
func (s *DB) SelectUsersObjects(username string) (res []ara.UserObjectSelectResp, err error) {
	var q = `SELECT o.id, o.comment, o.latitude, o.longitude FROM objects o
				JOIN users u ON u.id = o.user_id
				WHERE u.username = $1;`

	var rows *sql.Rows
	rows, err = s.Query(q, username)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = ara.HandleErrPair(rows.Close(), err)
	}()

	for rows.Next() {
		var obj ara.UserObjectSelectResp

		if err = rows.Scan(&obj.ID, &obj.Comment, &obj.Latitude, &obj.Longitude); err != nil {
			return nil, err
		}

		res = append(res, obj)
	}

	return res, nil
}

// UpdateObject deletes object info from db.
func (s *DB) UpdateObject(obj ara.ObjectUpdateInfo) error {
	var q = `UPDATE objects SET
				comment = $1 WHERE
				id = $2;`

	var _, err = s.Exec(q, obj.Comment, obj.ID)
	return err
}

// DeleteObject deletes object info from db.
func (s *DB) DeleteObject(id int) error {
	var q = `DELETE FROM objects o WHERE id = $1;`

	var _, err = s.Exec(q, id)
	return err
}
