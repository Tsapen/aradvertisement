package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
)

// Client provides redis functions.
type Client struct {
	*redis.Client
}

// Config is redis.Client cOnfig.
type Config struct {
	Dsn      string
	Password string
}

// NewRedisClient creates new redis client.
func NewRedisClient(c *Config) (*Client, error) {
	if len(c.Dsn) == 0 {
		c.Dsn = "auth_db:6379"
	}

	var client = redis.NewClient(&redis.Options{Addr: c.Dsn, Password: c.Password})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "can't ping redis")
	}

	return &Client{client}, nil
}
