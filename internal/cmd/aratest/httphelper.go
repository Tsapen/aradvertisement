package aratest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func getAddr(ctx context.Context, path string) string {
	var u = url.URL{
		Scheme: "http",
		Host:   AddrFromCtx(ctx),
		Path:   path,
	}
	return u.String()
}

func httpGet(ctx context.Context, t *testing.T, url string) smsi {
	t.Helper()
	var req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenFromCtx(ctx))
	var c = http.Client{Timeout: time.Second}
	var r *http.Response
	r, err = c.do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes []byte
	_, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil {
		t.Fatalf("can't read response")
	}

	var bodyMap smsi
	err = json.Unmarshal(bodyBytes, &bodyMap)
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

	req.Header.Set("Authorization", "Bearer "+tokenFromCtx(ctx))
	var c = http.Client{Timeout: time.Second}
	var r *http.Response
	r, err = c.Do(req)
	if err != nil {
		t.Fatalf("can't do request: %s", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes []byte
	_, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil {
		t.Fatalf("can't read response")
	}

	var bodyMap msi
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err != nil {
		t.Fatalf("can't parse response")
	}

	return bodyMap, nil
}

func badStatusFatal(t *testing.T, r *http.Response) {
	t.Fatalf("%s %s: bad response status:\nexp: %d\ngot: %d",
		r.Request.Method,
		r.Request.URL,
		r.StatusCode,
		http.StatusOK)
}

func login(ctx context.Context) context.Context {
	t.Helper()

	var url = getAddr(ctx, "/api/auth/login")
	var body = msi{
		"username": testUsername,
		"password": testPassword,
	}
	var tokens = httpPost(ctx, t, url, body)
	return withToken(ctx, tokens["access_token"].(string))
}

func logout(ctx context.Context) {
	t.Helper()

	url = getAddr(ctx, "/api/auth/logout")
	httpPost(ctx, t, url, nil)
}

func getLoc(lat, long string) string {
	return fmt.Sprintf("%.6f: %.6f", lat, long)
}
