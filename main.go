package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"stub_server/methods"
	"time"
)

type Config struct {
	Version string             `json:"version"`
	Host    string             `json:"host"`
	Port    string             `json:"port"`
	Urls    map[string][]methods.Query `json:"urls"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.SetFlags(log.LstdFlags)
}

func JsonConfig(path string) (Config, error) {
	conf := Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("error while read .config.json - %s\n", err.Error())
		return Config{}, err
	}
	if err := json.Unmarshal(data, &conf); err != nil {
		log.Printf("error while parse .config.json - %s\n", err.Error())
		return Config{}, err
	}
	return conf, nil
}

func YamlConfig(path string) (Config, error) {
	conf := Config{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("error while read config.yaml - %s\n", err.Error())
		return Config{}, err
	}
	if err := yaml.Unmarshal(data, &conf); err != nil {
		log.Printf("error while parse .config.yaml - %s\n", err.Error())
		return Config{}, err
	}
	return conf, nil
}

func main() {
	var path = "example_config.json"
	config := Config{}
	var err error
	configJson := flag.String("--json", "config.json", "setup server with .json config")
	configYaml := flag.String("--yaml", "", "setup server with .yaml config")
	flag.Parse()
	if *configYaml != "" {
		config, err = YamlConfig(*configYaml)
		if err != nil {
			os.Exit(2)
		}
	}
	if *configJson != "" {
		config, err = JsonConfig(path)
		if err != nil {
			os.Exit(2)
		}
	}

	handler := gin.New()
	handler.Use(Logger())
	for key, urls := range config.Urls {
		for _, url := range urls {
			for _, val := range url.Answers {
				switch url.Method {
				case "GET":
					handler.GET(key, func(ctx *gin.Context) {
						methods.Get(val, ctx)
					})
				case "POST":
					handler.POST(key, func(ctx *gin.Context) {
						methods.Post(val, ctx)
					})
				case "PUT":
					handler.PUT(key, func(ctx *gin.Context) {
						methods.Put(val, ctx)
					})
				case "DELETE":
					handler.DELETE(key, func(ctx *gin.Context) {
						methods.Delete(val, ctx)
					})
				case "PATCH":
					handler.PATCH(key, func(ctx *gin.Context) {
						methods.Patch(val, ctx)
					})
				case "HEAD":
					handler.HEAD(key, func(ctx *gin.Context) {
						methods.Head(val, ctx)
					})
				case "OPTIONS":
					handler.OPTIONS(key, func(ctx *gin.Context) {
						methods.Options(val, ctx)
					})
				}
			}
		}
	}

	routes := handler.Routes()
	log.Printf("Got %v handlers\n", len(routes))
	for key, val := range routes {
		log.Printf("%v - %+v\n", key, val)
	}

	server := &http.Server{
		Addr:              net.JoinHostPort(config.Host, config.Port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 10,
	}
	log.Printf("Server starting on - %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error while serving - \n\t%s", err.Error())
		return
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		url := c.Request.URL
		host := c.Request.Host
		log.Printf("Http_Server: Status - %3v |Latency - %6v |Host - %10v |Url - %40v ", status, latency, host, url)
	}
}
