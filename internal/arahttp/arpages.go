package arahttp

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tsapen/aradvertisement/internal/ara"
	"github.com/gorilla/mux"
)

func (h *handler) arPage(w http.ResponseWriter, r *http.Request) {
	var err error
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s\n", m, url)

	var v = mux.Vars(r)
	var idStr, okID = v["id"]

	var id int
	id, err = strconv.Atoi(idStr)
	if !okID || err != nil {
		processError(w, url, m, "id not found", nil)
		return
	}

	var obj ara.ObjectSelectByID
	obj, err = h.araDB.SelectObjectByID(id)

	var tmpFunc, ok = h.templMethods[obj.Type]
	if !ok {
		logError(url, m, "can't create template: unrecognized object type", err)
		return
	}

	if err = tmpFunc(w, obj.Username, id); err != nil {
		logError(url, m, "can't send message", err)
		return
	}
}
