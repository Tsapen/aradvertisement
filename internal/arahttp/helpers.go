package arahttp

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Tsapen/aradvertisement/internal/jwt"
	"github.com/pkg/errors"
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
	// 1. get object by id

	// 2. get template with it type

	// 3. execute template

	// var tmpl = h.tmps.Text

	// var params = struct {
	// 	Text string
	// }{}

	// tmpl.Execute(w, params)

	// var m = r.Method
	// var url = r.URL
	// log.Printf("%s %s", m, url)

	// var response, err = h.storage.GetARPage()
	// if err != nil {
	// 	logError(url, m, "open ar page file error", err)
	// 	return
	// }

	// w.Header().Set("Content-Type", "text/html")
	// if _, err := w.Write(response); err != nil {
	// 	logError(url, m, "can't send message", err)
	// 	return
	// }
}

func getFileWithPars(r *http.Request, pars interface{}) (file []byte, err error) {
	var mr *multipart.Reader
	mr, err = r.MultipartReader()
	if err != nil {
		return nil, errors.Wrap(err, "bad multipart reader creation")
	}

	for {
		var part *multipart.Part
		part, err = mr.NextPart()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = nil
			return
		}

		if err != nil {
			fmt.Println(part)
			return nil, errors.Wrap(err, "bad multipart iteration")
		}

		if part.FormName() == "gltf" {
			file = make([]byte, defaultSize)
			var n int
			n, err = io.ReadFull(part, file)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return nil, errors.Wrap(err, "bad file reading")
			}

			file = file[:n]
		}

		if part.FormName() == "info" {
			var buf = make([]byte, defaultSize)
			var n int
			n, err = io.ReadFull(part, buf)
			if err != nil && err != io.ErrUnexpectedEOF {
				return nil, errors.Wrap(err, "bad text form reading")
			}

			buf = buf[:n]
			if err = json.Unmarshal(buf, pars); err != nil {
				return nil, errors.Wrap(err, "bad json multipart marshalling")
			}
		}
	}
}
