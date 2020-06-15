package arahttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Tsapen/aradvertisement/internal/filestore"
	"github.com/Tsapen/aradvertisement/internal/jwt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	defaultJSONSize = 64 * 1024
	defaultFileSize = 128 * 1024 * 1024
)

func extractBody(r *http.Request, i interface{}) error {
	var buf = make([]byte, defaultJSONSize)
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

type textPars struct {
	Text string
}

func (h *handler) getTextTemplate(w http.ResponseWriter, username string, name int) error {
	var err error
	var tmpl = h.tmps.Text

	var text []byte
	text, err = h.storage.ReadFile(filestore.NewFM(username, name))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, textPars{string(text)})
	if err != nil {
		return err
	}

	return nil
}

type image struct {
	User string
	File string
}

type gltf struct {
	User string
	File string
}

func (h *handler) getImageTemplate(w http.ResponseWriter, username string, name int) error {
	var err error
	var tmpl = h.tmps.Img

	err = tmpl.Execute(w, image{User: username, File: strconv.Itoa(name)})
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) getGLTFTemplate(w http.ResponseWriter, username string, name int) error {
	var err error
	var tmpl = h.tmps.GLTF

	err = tmpl.Execute(w, gltf{User: username, File: strconv.Itoa(name)})
	if err != nil {
		return err
	}

	return nil
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

		if part.FormName() == "object" {
			file = make([]byte, defaultFileSize)
			var n int
			n, err = io.ReadFull(part, file)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return nil, errors.Wrap(err, "bad file reading")
			}

			file = file[:n]
		}

		if part.FormName() == "info" {
			var buf = make([]byte, defaultJSONSize)
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

func newTemplates(s *filestore.Storage) (*templates, error) {
	var tmplText, tmplImage, tmplGLTF *template.Template
	var path = s.GetTemplatePath("text_template.html")
	tmplText = template.Must(template.ParseFiles(path))

	path = s.GetTemplatePath("image_template.html")
	tmplImage = template.Must(template.ParseFiles(path))

	path = s.GetTemplatePath("gltf_template.html")
	tmplGLTF = template.Must(template.ParseFiles(path))

	var t = &templates{
		Text: tmplText,
		Img:  tmplImage,
		GLTF: tmplGLTF,
	}
	return t, nil
}

func (h *handler) arObject(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, url)
	var err error

	var v = mux.Vars(r)
	var user, okUser = v["user"]
	var file, okFile = v["file"]
	if !okUser || !okFile {
		processError(w, url, m, "user or file not found", nil)
		return
	}

	var path = filepath.Join("objects_storage", user, file)

	var f []byte
	f, err = ioutil.ReadFile(path)
	if err != nil {
		logError(url, m, "open ar page file error", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(f); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}
