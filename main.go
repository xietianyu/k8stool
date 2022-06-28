package main

import (
	"encoding/json"
	"io/ioutil"
	"k8stool/http_router"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

var serverConfig ServerConfig

// ServerConfig ...
type ServerConfig struct {
	HTTPPort string `json:"http_port"`
	IsHTTPS  bool   `json:"is_https"`
}

func main() {
	log.Printf("start main")
	err := LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Printf("config loaded")
	err = http_router.InitK8SClient()
	if err != nil {
		panic(err)
	}
	log.Printf("k8s client init success")
	SetupHttpListen()

}

// LoadConfig ....
func LoadConfig() error {

	ex, err := os.Executable()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	exePath := filepath.Dir(ex)

	fullPath := filepath.Join(exePath, "conf/config.json")

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	err = json.Unmarshal(data, &serverConfig)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	if serverConfig.HTTPPort == "" {
		serverConfig.HTTPPort = "8888"
	}
	return nil
}

func SetupHttpListen() {
	//Http router for rest api
	router := http_router.NewRouter()
	//COR: allow origin settings, for GET request only
	cor := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})
	handler := cor.Handler(router)
	log.Printf("http listen start")
	log.Printf("http port is %s", serverConfig.HTTPPort)
	err := http.ListenAndServe(":"+serverConfig.HTTPPort, handler)
	if err != nil {
		panic(err)
	}
}
