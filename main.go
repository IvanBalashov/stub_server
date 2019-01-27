package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Version string             `json:"version"`
	Host    string             `json:"host"`
	Port    string             `json:"port"`
	Urls    map[string][]Query `json:"urls"`
}

type Query struct {
	Method        string              `json:"method"`
	Queries       map[string]string   `json:"queries"`
	HttpStatus    int                 `json:"http_status"`
	Headers       map[string]string   `json:"headers"`
	Data          string              `json:"data"`
	MimeType      string              `json:"mime_type"`
	PostArguments map[string]string   `json:"post_arguments,omitempty"`
	File          string              `json:"file,omitempty"`
	Img           string              `json:"img,omitempty"`
	Cookies       map[string][]Cookie `json:"cookies"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	MaxAge   int    `json:"max_age,omitempty"`
	Path     string `json:"path"`
	Domain   string `json:"domain,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HttpOnly bool   `json:"http_only,omitempty"`
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.SetFlags(log.LstdFlags)
}

func ReadConfig() (Config, error) {
	conf := Config{}
	data, err := ioutil.ReadFile("config.json")
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

func main() {
	config, err := ReadConfig()
	if err != nil {
		os.Exit(2)
	}
	handler := gin.Default()

	for key, url := range config.Urls {
		for _, val := range url {
			switch val.Method {
			case "GET":
				handler.GET(key, func(ctx *gin.Context) {
					if val.Cookies != nil {
						setCookies(val.Cookies, ctx)
					}
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "POST":
				handler.POST(key, func(ctx *gin.Context) {
					if val.Cookies != nil {
						setCookies(val.Cookies, ctx)
					}
					setHeaders(val.Headers, ctx)
					setPostForm(val.PostArguments, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "PUT":
				handler.PUT(key, func(ctx *gin.Context) {
					if val.Cookies != nil {
						setCookies(val.Cookies, ctx)
					}
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "DELETE":
				handler.DELETE(key, func(ctx *gin.Context) {
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "PATCH":
				handler.PATCH(key, func(ctx *gin.Context) {
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "HEAD":
				handler.HEAD(key, func(ctx *gin.Context) {
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
			case "OPTIONS":
				handler.OPTIONS(key, func(ctx *gin.Context) {
					setHeaders(val.Headers, ctx)
					ctx.Data(http.StatusOK, val.MimeType, bytes.NewBufferString(val.Data).Bytes())
				})
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

func setHeaders(headers map[string]string, ctx *gin.Context) {
	for key, val := range headers {
		ctx.Header(key, val)
	}
}

func setPostForm(args map[string]string, ctx *gin.Context) {
	for key, val := range args {
		if val != ctx.PostForm(key) {
			ctx.JSON(http.StatusBadRequest, gin.H{})
		}
	}
}

func setCookies(cookies map[string][]Cookie, ctx *gin.Context) {
	for _, val := range cookies {
		for i := range val {
			ctx.SetCookie(val[i].Name,
				val[i].Value,
				val[i].MaxAge,
				val[i].Path,
				val[i].Domain,
				val[i].Secure,
				val[i].HttpOnly)
		}
	}
}
