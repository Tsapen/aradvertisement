package arahttp

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Tsapen/aradvertisement/internal/ara"
)

const (
	errUnauthorized ara.Error = "unauthorized"
)

type errWithStatus struct {
	status int
	err    error
}

func (err errWithStatus) Error() string {
	return err.Error()
}

func processError(w http.ResponseWriter, url *url.URL, method, message string, err error) {
	http.Error(w, message, errStatus(err))
	logError(url, method, message, err)
}

func logError(url *url.URL, method, message string, err error) {
	log.Printf("%s %s: %s: %s\n", method, url, message, err)
}

func errStatus(err error) int {
	switch errtype := err.(type) {
	case errWithStatus:
		return errtype.status

	case ara.BadParametersError:
		return http.StatusUnprocessableEntity

	default:

	}

	switch err {
	case ara.ErrBadParameters:
		return http.StatusUnprocessableEntity

	case errUnauthorized:
		return http.StatusUnauthorized

	default:

	}

	return http.StatusInternalServerError
}
