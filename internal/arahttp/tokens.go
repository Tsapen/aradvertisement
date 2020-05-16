package arahttp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Tsapen/aradvertisement/internal/ara"
	"github.com/Tsapen/aradvertisement/internal/auth"
	"github.com/Tsapen/aradvertisement/internal/jwt"
)

type user struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *handler) registration(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var u user
	if err := extractBody(r, &u); err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "can't extract body", err)
		return
	}

	var authUser = auth.User{Username: u.Username, Password: u.Password}
	var err = h.authDB.InsertUser(authUser)
	if err != nil {
		processError(w, url, m, "can't create user", err)
		return
	}

	var ts *auth.TokenDetails
	ts, err = jwt.CreateToken(u.Username)
	if err != nil {
		err = errWithStatus{http.StatusUnprocessableEntity, err}
		processError(w, url, m, "can't create token", err)
		return
	}

	if err := h.authDB.CreateAuth(u.Username, ts); err != nil {
		processError(w, url, m, "can't save token", err)
		return
	}

	if err := h.araDB.CreateUser(ara.UserCreationInfo{Username: u.Username}); err != nil {
		processError(w, url, m, "can't create user", err)
		return
	}

	if err := h.storage.CreateUserDir(u.Username); err != nil {
		processError(w, url, m, "can't create user dir", err)

		defer func() {
			if err = h.storage.DeleteUserDir(u.Username); err != nil {
				logError(url, m, "can't delete directory", err)
			}
		}()

		defer func() {
			if err = h.araDB.DeleteUser(u.Username); err != nil {
				logError(url, m, "can't delete user from db", err)
			}
		}()

		return
	}

	tokens := loginTokens{ts.AccessToken, ts.RefreshToken}
	var cookie = &http.Cookie{
		Name:     "token",
		Value:    ts.AccessToken,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var u loginData
	if err := extractBody(r, &u); err != nil {
		err = errWithStatus{http.StatusUnprocessableEntity, err}
		processError(w, url, m, "can't extract body", err)
		return
	}

	var ok, err = h.authDB.CheckLogin(auth.User(u))
	if err != nil {
		processError(w, url, m, "user not found", err)
		return
	}

	if !ok {
		processError(w, url, m, "can't login data is not valid", nil)
		return
	}

	var ts *auth.TokenDetails
	ts, err = jwt.CreateToken(u.Username)
	if err != nil {
		err = errWithStatus{http.StatusUnprocessableEntity, err}
		processError(w, url, m, "can't create token", err)
		return
	}

	saveErr := h.authDB.CreateAuth(u.Username, ts)
	if saveErr != nil {
		processError(w, url, m, "can't save token", err)
		return
	}

	tokens := loginTokens{ts.AccessToken, ts.RefreshToken}
	var cookie = &http.Cookie{
		Name:     "token",
		Value:    ts.AccessToken,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

type username struct {
	Username string `json:"username"`
}

func (h *handler) username(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var tokenstring = extractToken(r)
	var tokenAuth, err = jwt.ExtractTokenMetadata(tokenstring)
	if err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	var uname string
	uname, err = h.authDB.FetchAuth(tokenAuth)
	if err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	var u = username{uname}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&u); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var token = extractToken(r)
	var au, err = jwt.ExtractTokenMetadata(token)
	if err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	delErr := h.authDB.DeleteAuth(au.AccessUUID)
	if delErr != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	var resp = successResponse{true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

func (h *handler) refresh(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var mapToken map[string]string
	if err := extractBody(r, &mapToken); err != nil {
		err = errWithStatus{http.StatusUnprocessableEntity, err}
		processError(w, url, m, "can't extract body", err)
		return
	}

	var refreshToken = mapToken["refresh_token"]
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	var claims, err = jwt.Parse(refreshToken)
	if err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "refresh token expired", err)
		return
	}

	var refreshUUID, ok = claims["refresh_uuid"].(string)
	if !ok {
		err = ara.ErrBadParameters
		processError(w, url, m, err.Error(), err)
		return
	}

	var username string
	username, ok = claims["username"].(string)
	if !ok {
		err = ara.ErrBadParameters
		processError(w, url, m, "username isn't string", err)
		return
	}

	var delErr = h.authDB.DeleteAuth(refreshUUID)
	if delErr != nil {
		processError(w, url, m, "unauthorized", err)
		return
	}

	var ts, createErr = jwt.CreateToken(username)
	if createErr != nil {
		err = errWithStatus{http.StatusForbidden, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	var saveErr = h.authDB.CreateAuth(username, ts)
	if saveErr != nil {
		err = errWithStatus{http.StatusForbidden, err}
		processError(w, url, m, "unauthorized", err)
		return
	}

	var tokens = loginTokens{ts.AccessToken, ts.RefreshToken}
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

func (h *handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, url)

	var uname username
	if err := extractBody(r, &uname); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	if err := h.araDB.DeleteUser(uname.Username); err != nil {
		processError(w, url, m, "can't delete user", err)
		return
	}

	if err := h.storage.DeleteUserDir(uname.Username); err != nil {
		processError(w, url, m, "delete dir error", err)
		return
	}

	var resp = successResponse{Success: true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}
