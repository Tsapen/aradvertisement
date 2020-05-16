package aratest

import (
	"context"
	"testing"
)

// TestARA does integration testing.
func TestARA(t *testing.T, hostPort string) {
	var ctx = WithHostPort(context.Background())

	var testcases = []struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T)
	}{
		tc{name: "test user registration", testFunc: testRegisterUser},
		tc{name: "test user login", testFunc: testLoginUser},
		tc{name: "test user objects CRUD", testFunc: testCRUDObjects},
		tc{name: "test objects around location select", testFunc: testSelectObjectsAroundLocation},
		tc{name: "test user logout", testFunc: testLogoutUser},
		tc{name: "test user delete", testFunc: testDeletenUser},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, testcase.testFunc(ctx))
	}
}
