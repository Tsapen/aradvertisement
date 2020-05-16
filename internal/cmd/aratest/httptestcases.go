package aratest

import (
	"context"
	"reflect"
	"testing"
)

const (
	testUsername = "test_user"
	testEmail    = "test@email.com"
	testPassword = "test_password"

	lat1     = 50.000000
	lat2     = 51.000000
	lat3     = 51.000000
	long1    = 50.000000
	long2    = 51.000000
	long3    = 51.000000
	comment1 = "1st obj"
	comment2 = "2st obj"
	comment3 = "2st obj"
	degree   = 113200
)

type ei = interface{}
type f = float32
type msi = map[string]ei
type smsi = []msi

func testRegisterUser(ctx context.Context, t *testing.T) {
	var url = getAddr(ctx, "/api/auth/registration")
	var body = msi{
		"username": testUsername,
		"email":    testEmail,
		"password": testPassword,
	}
	var tokens = httpPost(ctx, t, url, body)
	ctx = withToken(ctx, tokens["access_token"])

	logout(ctx)
}

func testLoginUser(ctx context.Context, t *testing.T) {
	ctx = login(ctx)
	logout(ctx)
}

func testCRUDObjects(ctx context.Context, t *testing.T) {
	ctx = login(ctx)
	defer logout(ctx)

	// Create testing.
	var objs = smsi{
		msi{"latidude": lat1, "longitude": long1, comment: comment1, "glTF": []byte{}},
		msi{"latidude": lat2, "longitude": long2, comment: comment2, "glTF": []byte{}},
		msi{"latidude": lat3, "longitude": long3, comment: comment3, "glTF": []byte{}},
	}

	var url = getAddr(ctx, "/api/object")
	for _, obj := range objs {
		var resp = httpPost(t, url, obj)
		var success = resp["success"].(bool)
		if success != true {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", success, true)
		}
	}

	// Read testing.
	url = getAddr(ctx, "/api/user_objects")
	var r = httpGet(t, ctx, url)
	var exp = smsi{
		msi{"location": getLoc(lat1, long1), "comment": comment1},
		msi{"location": getLoc(lat2, long2), "comment": comment2},
		msi{"location": getLoc(lat3, long3), "comment": comment3},
	}
	if len(r) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	}

	for num, obj := range r {
		exp[num]["id"] = obj["id"]
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}

	// Update testing.
	url = getAddr(ctx, "/api/object/upd")
	var updComment = exp[2]["comment"] + " updated"
	var updated = msi{
		"id":      exp[2]["id"],
		"comment": updComment,
	}
	r = httpPost(t, ctx, url, upd)
	success = r["success"].(bool)
	if !success {
		t.Fatalf("unexpected result:\nexp: %v\ngot: %v\n", true, success)
	}

	exp[2]["comment"] = updComment
	url = getAddr(ctx, "/api/user_objects")
	var r = httpGet(t, ctx, url)

	if len(r) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	}

	for num, obj := range r {
		exp[num]["id"] = obj["id"]
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}

	// Delete object.
	url = getAddr(ctx, "/api/object/del")
	var del = msi{
		"id": exp[2]["id"],
	}
	r = httpPost(t, ctx, url, del)
	success = r["success"].(bool)
	if !success {
		t.Fatalf("unexpected result:\nexp: %v\ngot: %v\n", true, success)
	}

	exp = exp[:2]
	url = getAddr(ctx, "/api/user_objects")
	var r = httpGet(t, ctx, url)

	if len(r) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	}

	for num, obj := range r {
		exp[num]["id"] = obj["id"]
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}
}

func testSelectObjectsAroundLocation(ctx context.Context, t *testing.T) {

}

func testDeletenUser(ctx context.Context, t *testing.T) {

	var url = getAddr(ctx, "/api/auth/login")
	var body = msi{
		"username": testUsername,
		"password": testPassword,
	}
	var tokens = httpPost(ctx, t, url, body)
	ctx = withToken(ctx, tokens["access_token"].(string))
}
