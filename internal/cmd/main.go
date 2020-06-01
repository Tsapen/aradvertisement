package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Tsapen/aradvertisement/internal/arahttp"
	"github.com/Tsapen/aradvertisement/internal/filestore"
	"github.com/Tsapen/aradvertisement/internal/jwt"
	"github.com/Tsapen/aradvertisement/internal/postgres"
	"github.com/Tsapen/aradvertisement/internal/redis"
)

const (
	certDir = `certs`
)

type config struct {
	Radius           int      `json:"radius"`
	StorageDirectory string   `json:"storage_directory"`
	Auth             authCfg  `json:"auth"`
	HTTP             httpCfg  `json:"http"`
	AraDB            dbCfg    `json:"ara_db"`
	Redis            redisCfg `json:"redis"`
}

type authCfg struct {
	SetManually   bool   `json:"set_manually"`
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
}

type httpCfg struct {
	Port         string `json:"port"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

type dbCfg struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	HostName    string `json:"hostname"`
	Port        string `json:"port"`
	VirtualHost string `json:"virtual_host"`
}

type redisCfg struct {
	Dsn      string `json:"dsn"`
	Password string `json:"password"`
}

func main() {
	var cPath = os.Getenv("ARA_CONFIG")
	if cPath == "" {
		panic(fmt.Sprintf("config path should be set in environment variable ARA_CONFIG"))
	}

	run(openConfig(cPath))
}

func run(c *config) {
	var err error
	var pwd string
	pwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("can't prepare auth db: %s", err))
	}

	var sec = jwt.Secrets(c.Auth)
	if err := jwt.PrepareAuthEnvironment(sec); err != nil {
		panic(fmt.Sprintf("can't prepare jwt environment: %s", err))
	}

	var s *filestore.Storage
	s, err = filestore.PrepareStorage(c.StorageDirectory)
	if err != nil {
		panic(fmt.Sprintf("can't connect with storage: %s", err))
	}

	// ara db
	var cAraDB = postgres.Config(c.AraDB)
	var araDB *postgres.DB
	araDB, err = postgres.NewDBConnection(&cAraDB)
	if err != nil {
		panic(fmt.Sprintf("can't connect with ara db: %s", err))
	}

	// auth db
	var redisClient *redis.Client
	var cRedis = &redis.Config{Dsn: c.Redis.Dsn, Password: c.Redis.Password}
	redisClient, err = redis.NewRedisClient(cRedis)
	if err != nil {
		panic(fmt.Sprintf("can't prepare auth db: %s", err))
	}

	if err := araDB.AraMigrate(); err != nil {
		panic(fmt.Sprintf("can't prepare ara db: %s", err))
	}

	var cHTTP = arahttp.Config{
		Port:         c.HTTP.Port,
		ReadTimeout:  c.HTTP.ReadTimeout,
		WriteTimeout: c.HTTP.WriteTimeout,
		AraDB:        araDB,
		AuthDB:       redisClient,
		Storage:      s,
	}

	var api *arahttp.API
	api, err = arahttp.NewAPI(&cHTTP)
	if err != nil {
		panic(fmt.Sprintf("can't start api: %s", err))
	}

	log.Printf("main: start listening %s", c.HTTP.Port)

	api.Start(getCertPath(pwd, "/server.crt"), getCertPath(pwd, "/server.key"))
}

func getCertPath(projectDir, fileName string) string {
	return filepath.Join(projectDir, certDir, fileName)
}

func openConfig(cPath string) *config {
	var cFile, err = os.Open(cPath)
	if err != nil {
		panic(fmt.Sprintf("can't open config: %s", err))
	}

	var c = &config{}
	err = json.NewDecoder(cFile).Decode(c)
	if err != nil {
		panic(fmt.Sprintf("can't encode config: %s", err))
	}

	return c
}
