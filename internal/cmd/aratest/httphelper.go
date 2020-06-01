package aratest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const (
	defaultSize = 16 * 1024
)

func getAddr(ctx context.Context, path string) string {
	var u = &url.URL{
		Scheme: "http",
		Host:   hostPortFromCtx(ctx),
		Path:   path,
	}
	return u.String()
}

func getAddrWithParams(ctx context.Context, path, rawQuery string) string {
	var u = url.URL{
		Scheme:   "https",
		Host:     hostPortFromCtx(ctx),
		Path:     path,
		RawQuery: rawQuery,
	}
	return u.String()
}

func httpGet(ctx context.Context, t *testing.T, url string) msi {
	t.Helper()
	var req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenFromCtx(ctx))
	var c = http.Client{Timeout: time.Second}
	var r *http.Response
	r, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			t.Errorf("bad resp closing: %s\n", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes = make([]byte, defaultSize)
	var n int
	n, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Fatalf("can't read response")
	}

	var bodyMap msi
	err = json.Unmarshal(bodyBytes[:n], &bodyMap)
	if err != nil {
		t.Fatalf("can't parse response")
	}

	return bodyMap
}

func httpPost(ctx context.Context, t *testing.T, url string, body ei) msi {
	t.Helper()

	var b, err = json.Marshal(body)
	if err != nil {
		t.Fatalf("can't marshal request: %s", err)
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("can't create request: %s", err)
	}

	var token = tokenFromCtx(ctx)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	var r *http.Response
	var c = http.Client{Timeout: time.Second}
	r, err = c.Do(req)
	if err != nil {
		t.Fatalf("can't do request: %s", err)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			t.Errorf("bad resp closing: %s\n", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes = make([]byte, 16*1024)
	var n int
	n, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Fatalf("can't read response: %s", err)
	}

	var bodyMap msi
	err = json.Unmarshal(bodyBytes[:n], &bodyMap)
	if err != nil {
		t.Fatalf("can't parse response: %s", err)
	}

	return bodyMap
}

func badStatusFatal(t *testing.T, r *http.Response) {
	t.Fatalf("%s %s: bad response status:\nexp: %d\ngot: %d",
		r.Request.Method,
		r.Request.URL,
		r.StatusCode,
		http.StatusOK)
}

func login(ctx context.Context, t *testing.T) context.Context {
	t.Helper()

	var url = getAddr(ctx, "/api/auth/login")
	var body = msi{
		"username": testUsername,
		"password": testPassword,
	}
	var tokens = httpPost(ctx, t, url, body)
	return withToken(ctx, tokens["access_token"].(string))
}

func logout(ctx context.Context, t *testing.T) {
	t.Helper()

	var url = getAddr(ctx, "/api/auth/logout")
	httpPost(ctx, t, url, nil)
}

func getLoc(lat, long float64) string {
	return fmt.Sprintf("%.6f: %.6f", lat, long)
}

func createFile(ctx context.Context, t *testing.T, params msi) {
	t.Helper()

	var err error
	var body = new(bytes.Buffer)
	var writer = multipart.NewWriter(body)
	var filePart io.Writer
	filePart, err = writer.CreateFormFile("gltf", "test.gltf")
	if err != nil {
		t.Fatalf("bad file writer creation: %s", err)
	}

	var fileContents = []byte("12345")
	filePart.Write(fileContents)

	var jsonPart io.Writer
	jsonPart, err = writer.CreateFormField("info")
	if err != nil {
		t.Fatalf("bad json writer creation: %s", err)
	}

	var jsonContents []byte
	jsonContents, err = json.Marshal(params)
	if err != nil {
		t.Fatalf("bad json marshalling: %s", err)
	}

	jsonPart.Write(jsonContents)

	err = writer.Close()
	if err != nil {
		t.Fatalf("writer closing: %s", err)
	}

	var url = getAddr(ctx, "/api/object")
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Fatalf("can't create request: %s", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	var token = tokenFromCtx(ctx)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	var r *http.Response
	var c = http.Client{Timeout: time.Second}
	r, err = c.Do(req)
	if err != nil {
		t.Fatalf("can't do request: %s", err)
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			t.Errorf("bad resp closing: %s\n", err)
		}
	}()

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes = make([]byte, 16*1024)
	var n int
	n, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Fatalf("can't read response: %s", err)
	}

	var bodyMap msi
	err = json.Unmarshal(bodyBytes[:n], &bodyMap)
	if err != nil {
		t.Fatalf("can't parse response: %s", err)
	}

	if res := bodyMap["success"].(bool); !res {
		t.Fatalf("exspected success response")
	}
}
