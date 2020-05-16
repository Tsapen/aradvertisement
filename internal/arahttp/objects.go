package arahttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tsapen/aradvertisement/internal/ara"
	"github.com/Tsapen/aradvertisement/internal/filestore"
	"github.com/gorilla/mux"
)

type objectResp struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	GLTF      []byte  `json:"gltf"`
}

type objectsResponse struct {
	Response []objectResp `json:"response"`
}

func (h *handler) objectsByLocation(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var v = mux.Vars(r)
	var latStr, okLat = v["latitude"]
	var longStr, okLong = v["longitude"]
	if !okLat || !okLong {
		processError(w, url, m, "latitude or longitude not found", nil)
		return
	}

	var lat, errLat = strconv.ParseFloat(latStr, 64)
	var long, errLong = strconv.ParseFloat(longStr, 64)
	if errLat != nil || errLong != nil {
		processError(w, url, m, "latitude or longitude are not float", nil)
		return
	}

	var selectPars = ara.ObjectSelectInfo{
		Latitude:  lat,
		Longitude: long,
	}
	var objectsInfo, err = h.araDB.SelectObjectsAround(selectPars)
	if err != nil {
		processError(w, url, m, "db error", err)
		return
	}

	var objects = make([]objectResp, len(objectsInfo))
	for num, objInfo := range objectsInfo {
		var fm = filestore.NewFM(objInfo.Username, objInfo.ID)
		var b, err = h.storage.ReadFile(fm)
		if err != nil {
			logError(url, m, "can't read file", err)
			continue
		}

		objects[num].Username = objInfo.Username
		objects[num].Latitude = objInfo.Latitude
		objects[num].Longitude = objInfo.Longitude
		objects[num].GLTF = b
	}

	var resp = objectsResponse{Response: objects}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		processError(w, url, m, "sending response error", err)
		return
	}
}

type userObjectInfoResponse struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
	Comment  string `json:"comment"`
}

type userObjectsInfoResponse struct {
	Response []userObjectInfoResponse `json:"response"`
}

func (h *handler) userObjects(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, r.URL)

	var username, err = getUsername(r)
	if err != nil {
		processError(w, url, m, "username not found", err)
		return
	}

	var objectsInfo []ara.UserObjectSelectResp
	objectsInfo, err = h.araDB.SelectUsersObjects(username)
	if err != nil {
		processError(w, url, m, "db error", err)
		return
	}
	var objects = make([]userObjectInfoResponse, len(objectsInfo))
	for num, objInfo := range objectsInfo {
		objects[num].ID = objInfo.ID
		objects[num].Location = fmt.Sprintf("%.6f: %.6f", objInfo.Latitude, objInfo.Longitude)
		objects[num].Comment = objInfo.Comment
	}

	var resp = userObjectsInfoResponse{Response: objects}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		processError(w, url, m, "sending response error", err)
		return
	}
}

type newObjectReq struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	GLTF      []byte  `json:"gltf"`
}

func (h *handler) newObject(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var obj newObjectReq
	if err := extractBody(r, &obj); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	var username, err = getUsername(r)
	if err != nil {
		processError(w, url, m, "username not found", err)
		return
	}

	var creationInfo = ara.ObjectCreationInfo{username, obj.Latitude, obj.Longitude}
	var pars = ara.ObjectCreationInfo(creationInfo)

	var id int
	id, err = h.araDB.CreateObject(pars)
	if err != nil {
		processError(w, url, m, "can't create object", err)
		return
	}

	var fm = filestore.NewFM(username, id)
	if err := h.storage.WriteFile(fm, obj.GLTF); err != nil {
		processError(w, url, m, "can't write file", err)

		defer func() {
			if err = h.storage.DeleteFile(fm); err != nil {
				logError(url, m, "can't delete file from file system", err)
			}
		}()

		defer func() {
			if err = h.araDB.DeleteObject(id); err != nil {
				logError(url, m, "can't delete file from db", err)
			}
		}()

		return
	}

	var resp = successResponse{Success: true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

type objectUpdateInfo struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

type successResponse struct {
	Success bool `json:"success"`
}

func (h *handler) updateObject(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var info objectUpdateInfo
	if err := extractBody(r, &info); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	if err := h.araDB.UpdateObject(ara.ObjectUpdateInfo(info)); err != nil {
		processError(w, url, m, "can't update object", err)
		return
	}

	var resp = successResponse{Success: true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}

type objectDeleteInfo struct {
	ID int `json:"id"`
}

func (h *handler) deleteObject(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var username, err = getUsername(r)
	if err != nil {
		err = errWithStatus{http.StatusUnauthorized, err}
		processError(w, url, m, "can't get username", err)
		return
	}

	var info objectDeleteInfo
	if err := extractBody(r, &info); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	err = h.araDB.DeleteObject(info.ID)
	if err != nil {
		processError(w, url, m, "can't delete object", err)
		return
	}

	var fm = filestore.NewFM(username, info.ID)
	if err := h.storage.DeleteFile(fm); err != nil {
		processError(w, url, m, "delete file error", err)
		return
	}

	var resp = successResponse{Success: true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}
