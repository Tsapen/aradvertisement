package redis

import (
	"time"

	"github.com/Tsapen/aradvertisement/internal/auth"
)

// FetchAuth returns username by token.
func (client *Client) FetchAuth(authD *auth.AccessDetails) (string, error) {
	username, err := client.Get(authD.AccessUUID).Result()
	if err != nil {
		return "", err
	}

	return username, nil
}

// CreateAuth put token with username in redis.
func (client *Client) CreateAuth(username string, td *auth.TokenDetails) error {
	var at = time.Unix(td.AtExpires, 0)
	var rt = time.Unix(td.RtExpires, 0)
	var now = time.Now()

	if err := client.Set(td.AccessUUID, username, at.Sub(now)).Err(); err != nil {
		return err
	}

	if err := client.Set(td.RefreshUUID, username, rt.Sub(now)).Err(); err != nil {
		return err
	}

	return nil
}

// DeleteAuth deletes user token.
func (client *Client) DeleteAuth(givenUUID string) error {
	_, err := client.Del(givenUUID).Result()
	if err != nil {
		return err
	}
	return nil
}
