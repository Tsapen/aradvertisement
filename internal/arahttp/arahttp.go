package arahttp

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Tsapen/aradvertisement/internal/ara"
	"github.com/Tsapen/aradvertisement/internal/auth"
	"github.com/Tsapen/aradvertisement/internal/filestore"
	"github.com/gorilla/mux"
)

// API contains configured http.Server.
type API struct {
	*http.Server
}

// Config for running http server.
type Config struct {
	Addr         string
	Port         string
	ReadTimeout  string
	WriteTimeout string
	AraDB        ara.DB
	AuthDB       auth.DB
	Storage      *filestore.Storage
}

type handler struct {
	port    string
	araDB   ara.DB
	authDB  auth.DB
	router  *mux.Router
	storage *filestore.Storage
	tmps    templates
}

type templates struct {
	Text *template.Template
	Img  *template.Template
	GLTF *template.Template
}

// NewAPI creates new api.
func NewAPI(config *Config) (*API, error) {
	r := mux.NewRouter()
	var h = handler{
		port:    config.Port,
		araDB:   config.AraDB,
		authDB:  config.AuthDB,
		router:  r,
		storage: config.Storage,
	}

	// Object handlers API.
	r.HandleFunc("/api/objects_around", h.objectsByLocation).
		Queries("latitude", "{latitude}", "longitude", "{longitude}").Methods(http.MethodGet)
	r.HandleFunc("/api/user_objects", h.userObjects).Methods(http.MethodGet)
	r.HandleFunc("/api/object", h.newObject).Methods(http.MethodPost)
	r.HandleFunc("/api/object/upd", h.updateObject).Methods(http.MethodPost)
	r.HandleFunc("/api/object/del", h.deleteObject).Methods(http.MethodPost)

	// Auth handlers API.
	r.HandleFunc("/api/auth/registration", h.registration).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", h.login).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/logout", h.logout).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/refresh", h.refresh).Methods(http.MethodPost)

	r.HandleFunc("/api/ar", h.arPage).
		Queries("id", "{id}").Methods(http.MethodGet)

	var wrappedHandler = withHeaders(h.router)

	var readTimeout, err = time.ParseDuration(config.ReadTimeout)
	if err != nil {
		return nil, err
	}

	var writeTimeout time.Duration
	writeTimeout, err = time.ParseDuration(config.WriteTimeout)
	if err != nil {
		return nil, err
	}

	var s = newServer(config.Port, wrappedHandler, readTimeout, writeTimeout)

	return &API{s}, nil
}

func newServer(addr string, h http.Handler, rt, wt time.Duration) *http.Server {

	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("192.168.1.52"),
	// 	Cache:      autocert.DirCache("certs"),
	// }

	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  rt,
		WriteTimeout: wt,
		// TLSConfig: &tls.Config{
		// 	GetCertificate:     certManager.GetCertificate,
		// 	InsecureSkipVerify: true,
		// 	ClientAuth:         tls.RequireAndVerifyClientCert,
		// },
	}
}

// Start run server.
func (a *API) Start(crt, key string) {
	if err := a.ListenAndServeTLS(crt, key); err != nil {
		// if err := a.ListenAndServe(); err != nil {
		log.Printf("api.Start error: %s\n", err)
	}
}
