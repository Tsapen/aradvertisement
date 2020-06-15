package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Tsapen/aradvertisement/internal/cmd/aratest"

	"github.com/go-redis/redis"
)

const (
	testPath = "/api/objects_around"
)

func TestMain(m *testing.M) {
	var cPath = os.Getenv("ARA_TEST_CONFIG")
	if cPath == "" {
		log.Printf("skip tests: ARA_TEST_CONFIG doesn't contain path to config\n")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestARA(t *testing.T) {
	var cPath = os.Getenv("ARA_TEST_CONFIG")
	if cPath == "" {
		t.Fatalf("skip tests: ARA_TEST_CONFIG doesn't contain path to config\n")
	}

	var c = openConfig(cPath)
	cleanUpDB(t, c)
	go run(c)
	var addr = "localhost" + c.HTTP.MainPort
	waitRunning(t, addr)

	aratest.TestARA(t, addr)
}

func waitRunning(t *testing.T, addr string) {
	const checkNum = 10
	const maxDelay = 100 * time.Millisecond

	for i := 0; i < checkNum; i++ {
		time.Sleep(maxDelay)

		if _, err := http.Get("http://" + addr + testPath); err == nil {
			return
		}
	}

	t.Fatalf("service could not start")
}

func getDBAddr(c dbCfg) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.UserName, c.Password, c.HostName, c.Port, c.VirtualHost)

}

func cleanUpDB(t *testing.T, c *config) {
	var err error
	var addr = getDBAddr(c.AraDB)

	var db *sql.DB
	db, err = sql.Open("postgres", addr)
	if err != nil {
		t.Fatalf("can't open connection: %s", err)
	}

	var queries = []string{
		`DROP TABLE IF EXISTS migrations`,
		`DROP TABLE IF EXISTS objects`,
		`DROP TABLE IF EXISTS users`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			t.Fatalf("can't do ara db query: %s", err)
		}
	}

	var kvClient = redis.NewClient(&redis.Options{Addr: c.Redis.Dsn})
	_, err = kvClient.Ping().Result()
	if err != nil {
		t.Fatalf("can't ping redis: %s", err)
	}

	queries = []string{
		`flushall`,
	}
	for _, query := range queries {
		if err := kvClient.Do(query).Err(); err != nil {
			t.Fatalf("can't do redis query: %s", err)
		}
	}
}

func clearFileStorage(path string) {
	os.RemoveAll(filepath.Join(path, "gltf"))
}
