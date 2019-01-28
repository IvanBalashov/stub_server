package methods

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"time"
)

func Get(val Answers, ctx *gin.Context) {
	errMsg := ErrMsq{
		Error: "error",
	}
	noErr := true
	args := ResponseArgs{
		Code:     200,
		MimeType: "text/plain",
		Data:     []byte{},
	}

	if val.Cookies != nil {
		setCookies(val.Cookies, ctx)
	}

	if val.DataFromFile != "" {
		uploadFile(val.DataFromFile, ctx)
	}

	if val.ResponseHeaders != nil {
		setHeaders(val.ResponseHeaders, ctx)
	}

	if val.RequestHeaders != nil {
		if !checkHeaders(val.RequestHeaders, ctx) {
			noErr = false
			errMsg.Error = "Request Headers doesn't equal"
		}
	}

	if val.Queries != nil {
		if checkQueries(val.Queries, ctx) {
			if _, ok := val.Queries["query_data"]; ok {
				args.Data = bytes.NewBufferString(val.Queries["query_data"]).Bytes()
			} else {
				if val.Data != "" && len(args.Data) == 0  {
					args.Data = bytes.NewBufferString(val.Data).Bytes()
				} else {
					noErr = false
					errMsg.Error = "Don't know what return..."
				}
			}
			args.Code = val.HttpStatus
			args.MimeType = val.MimeType
		} else {
			if val.Data == "" && len(args.Data) == 0  {
				noErr = false
				errMsg.Error = "Don't know what return..."
			} else {
				args.Code = val.HttpStatus
				args.MimeType = val.MimeType
				args.Data = bytes.NewBufferString(val.Data).Bytes()
			}
		}
	}

	if val.WaitTime != "" {
		waitTime, _ := time.ParseDuration(val.WaitTime)
		time.Sleep(waitTime)
	}

	if !noErr {
		args.Code = http.StatusBadRequest
		args.MimeType = "application/json"
		data, _ := json.Marshal(errMsg)
		args.Data = data
		ctx.Data(args.Code, args.MimeType, args.Data)
	} else {
		ctx.Data(args.Code, args.MimeType, args.Data)
	}
}
