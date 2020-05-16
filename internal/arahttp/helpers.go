package arahttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Tsapen/aradvertisement/internal/jwt"
)

const defaultSize = 16 * 1024

func extractBody(r *http.Request, i interface{}) error {
	var buf = make([]byte, defaultSize)
	var n, err = io.ReadFull(r.Body, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil
	}

	buf = buf[:n]
	if err := json.Unmarshal(buf, i); err != nil {
		return err
	}

	return nil
}

func getUsername(r *http.Request) (string, error) {
	var tokenstring = extractToken(r)
	var tokenAuth, err = jwt.ExtractTokenMetadata(tokenstring)
	if err != nil {
		return "", err
	}
	return tokenAuth.Username, nil
}

func extractToken(r *http.Request) string {
	var bearToken = r.Header.Get("Authorization")
	var strArr = strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func (h *handler) arPage(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, url)

	var response, err = h.storage.GetARPage()
	if err != nil {
		logError(url, m, "open ar page file error", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(response); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}
