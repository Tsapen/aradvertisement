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
	"github.com/pkg/errors"
)

// API contains configured http.Server.
type API struct {
	*http.Server
}

// Servers contains all APIs.
type Servers struct {
	mainAPI      *API
	templatesAPI *API
}

// Config for running http server.
type Config struct {
	Addr         string
	MainPort     string
	ArPort       string
	ReadTimeout  string
	WriteTimeout string
	AraDB        ara.DB
	AuthDB       auth.DB
	Storage      *filestore.Storage
}

type handler struct {
	mainPort     string
	ArPort       string
	araDB        ara.DB
	authDB       auth.DB
	router       *mux.Router
	storage      *filestore.Storage
	tmps         *templates
	templMethods map[rune]func(http.ResponseWriter, string, int) error
}

type templates struct {
	Text *template.Template
	Img  *template.Template
	GLTF *template.Template
}

// NewServers creates new servers.
func NewServers(config *Config) (*Servers, error) {
	var err error
	r := mux.NewRouter()
	var h = handler{
		mainPort: config.MainPort,
		araDB:    config.AraDB,
		authDB:   config.AuthDB,
		router:   r,
		storage:  config.Storage,
	}

	h.tmps, err = newTemplates(h.storage)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse templates")
	}

	h.templMethods = map[rune]func(http.ResponseWriter, string, int) error{
		't': h.getTextTemplate,
		'i': h.getImageTemplate,
		'g': h.getGLTFTemplate,
	}

	// Object handlers API.
	r.HandleFunc("/api/nearest_objects", h.objectsByLocation).Methods(http.MethodGet)
	r.HandleFunc("/api/user_objects", h.userObjects).Methods(http.MethodGet)
	r.HandleFunc("/api/object", h.newObject).Methods(http.MethodPost)
	r.HandleFunc("/api/object/upd", h.updateObject).Methods(http.MethodPost)
	r.HandleFunc("/api/object/del", h.deleteObject).Methods(http.MethodPost)

	// Auth handlers API.
	r.HandleFunc("/api/auth/registration", h.registration).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/login", h.login).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/logout", h.logout).Methods(http.MethodPost)
	r.HandleFunc("/api/auth/refresh", h.refresh).Methods(http.MethodPost)

	r.HandleFunc("/api/objects_around", h.objectsByLocation).
		Queries("latitude", "{latitude}", "longitude", "{longitude}").Methods(http.MethodGet)
	r.HandleFunc("/api/ar/page/{id}", h.arPage).Methods(http.MethodGet)
	r.HandleFunc("/api/ar/object", h.arObject).
		Queries("user", "{user}", "file", "{file}").Methods(http.MethodGet)

	var wrappedHandler = withHeaders(h.router)

	var readTimeout time.Duration
	readTimeout, err = time.ParseDuration(config.ReadTimeout)
	if err != nil {
		return nil, err
	}

	var writeTimeout time.Duration
	writeTimeout, err = time.ParseDuration(config.WriteTimeout)
	if err != nil {
		return nil, err
	}

	var mainServer = newHTTPServer(config.MainPort, wrappedHandler, readTimeout, writeTimeout)
	var templatesServer = newHTTPSServer(config.ArPort, wrappedHandler, readTimeout, writeTimeout)

	var s = &Servers{
		mainAPI:      &API{mainServer},
		templatesAPI: &API{templatesServer},
	}
	return s, nil
}

func newHTTPServer(addr string, h http.Handler, rt, wt time.Duration) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  rt,
		WriteTimeout: wt,
	}
}

func newHTTPSServer(addr string, h http.Handler, rt, wt time.Duration) *http.Server {
	// var certManager = autocert.Manager{
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
func (s *Servers) Start(crt, key string) {
	go func(crt, key string) {
		if err := s.templatesAPI.ListenAndServeTLS(crt, key); err != nil {
			log.Printf("api.Start error: %s\n", err)
		}
	}(crt, key)

	if err := s.mainAPI.ListenAndServe(); err != nil {
		log.Printf("api.Start error: %s\n", err)
	}
}
