package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Tsapen/aradvertisement/internal/cmd/aratest"
)

const (
	testPath = "/login"
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
	go run(c)
	waitRunning(t)

	var hostPort = net.JoinHostPort("localhost", c.HTTP.Port)
	aratest.TestARA(t, hostPort)
}

func waitRunning(t *testing.Tg, addr string) {
	const checkNum = 10
	const maxDelay = 100 * time.Millisecond

	for i := 0; i < checkNum; i++ {
		time.Sleep(maxDelay)

		if _, err := http.Get(addr + testPath); err == nil {
			return
		}
	}

	t.Fatalf("service could not start")
}
