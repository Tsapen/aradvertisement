package aratest

import (
	"context"
	"testing"
)

type tc struct {
	name     string
	testFunc func(ctx context.Context, t *testing.T)
}

// TestARA does integration testing.
func TestARA(t *testing.T, hostPort string) {
	var ctx = withHostPort(context.Background(), hostPort)

	var testcases = []tc{
		{name: "test user registration", testFunc: testRegisterUser},
		{name: "test user login", testFunc: testLoginUser},
		{name: "test user objects CRUD", testFunc: testCRUDObjects},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.testFunc(ctx, t)
		})
	}
}
