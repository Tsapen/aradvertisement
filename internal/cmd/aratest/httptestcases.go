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
	lat3     = 51.000001
	long1    = 50.000000
	long2    = 51.000000
	long3    = 51.000001
	comment1 = "1st obj"
	comment2 = "2st obj"
	comment3 = "2st obj"
	// degree   = 113200
)

type ei = interface{}
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
	ctx = withToken(ctx, tokens["access_token"].(string))

	logout(ctx, t)
}

func testLoginUser(ctx context.Context, t *testing.T) {
	ctx = login(ctx, t)
	logout(ctx, t)
}

func getDefaultObjects() smsi {
	return smsi{
		msi{"latitude": lat1, "longitude": long1, "comment": comment1, "glTF": []byte{}},
		msi{"latitude": lat2, "longitude": long2, "comment": comment2, "glTF": []byte{}},
		msi{"latitude": lat3, "longitude": long3, "comment": comment3, "glTF": []byte{}},
	}
}

func testCRUDObjects(ctx context.Context, t *testing.T) {
	var url string
	// var err error
	ctx = login(ctx, t)
	defer logout(ctx, t)

	// 1. Create testing.
	var objs = getDefaultObjects()
	for _, obj := range objs {
		createFile(ctx, t, obj)

	}

	// 2. Read testing.
	// 2.1. User objects reading.
	url = getAddr(ctx, "/api/user_objects")
	var r = httpGet(ctx, t, url)
	var exp = smsi{
		msi{"location": getLoc(lat1, long1), "comment": comment1},
		msi{"location": getLoc(lat2, long2), "comment": comment2},
		msi{"location": getLoc(lat3, long3), "comment": comment3},
	}
	var objInfo = r["response"].([]ei)

	if len(objInfo) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	}

	for num, obj := range objInfo {
		exp[num]["id"] = obj.(msi)["id"]
	}

	for num, obj := range objInfo {
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}

	// // 2.2. Objects around reading.
	// var rawQuery = fmt.Sprintf("latitude=%f&longitude=%f", lat2, long2)
	// url = getAddrWithParams(ctx, "/api/objects_around", rawQuery)
	// r = httpGet(ctx, t, url)
	// var objAround = r["response"].(smsi)

	// var expObjAround = objs[1:]
	// if len(objAround) != len(expObjAround) {
	// 	t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	// }

	// for num, obj := range objAround {
	// 	expObjAround[num]["id"] = obj["id"]
	// }

	// for num, obj := range objAround {
	// 	if !reflect.DeepEqual(expObjAround[num], obj) {
	// 		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", expObjAround[num], obj)
	// 	}
	// }

	// 3. Update testing.
	url = getAddr(ctx, "/api/object/upd")

	var updComment = exp[2]["comment"].(string) + " updated"
	exp[2]["comment"] = updComment
	var updated = exp[2]

	r = httpPost(ctx, t, url, updated)
	var success = r["success"].(bool)
	if !success {
		t.Fatalf("unexpected result:\nexp: %v\ngot: %v\n", true, success)
	}

	url = getAddr(ctx, "/api/user_objects")
	r = httpGet(ctx, t, url)
	objInfo = r["response"].([]ei)

	if len(objInfo) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, r)
	}

	for num, obj := range objInfo {
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}

	// 4. Delete object.
	url = getAddr(ctx, "/api/object/del")
	var delReqBody = msi{
		"id": exp[2]["id"],
	}

	r = httpPost(ctx, t, url, delReqBody)
	success = r["success"].(bool)
	if !success {
		t.Fatalf("unexpected result:\nexp: %v\ngot: %v\n", true, success)
	}

	exp = exp[:2]
	url = getAddr(ctx, "/api/user_objects")

	r = httpGet(ctx, t, url)
	objInfo = r["response"].([]ei)

	if len(objInfo) != len(exp) {
		t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp, objInfo)
	}

	for num, obj := range objInfo {
		if !reflect.DeepEqual(exp[num], obj) {
			t.Errorf("unexpected result:\nexp: %v\ngot: %v\n", exp[num], obj)
		}
	}
}
