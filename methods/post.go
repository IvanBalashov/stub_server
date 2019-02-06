package methods

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Post(answers map[string]Answers, ctx *gin.Context) {
	errMsg := ErrMsq{
		Error: "error",
	}

	noErr := true
	success := true
	args := ResponseArgs{
		Code:     200,
		MimeType: "text/plain",
		Data:     []byte{},
	}

	for _, val := range answers {
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
				success = false
			}
		}

		if val.Queries != nil {
			if checkQueries(val.Queries, ctx) {
				if _, ok := val.Queries["query_data"]; ok {
					args.Data = bytes.NewBufferString(val.Queries["query_data"]).Bytes()
				} else {
					if val.Data != "" && len(args.Data) == 0 {
						args.Data = bytes.NewBufferString(val.Data).Bytes()
					} else {
						noErr = false
						errMsg.Error = "Don't know what return..."
						success = false
					}
				}
				args.Code = val.HttpStatus
				args.MimeType = val.MimeType
			} else {
				if val.Data == "" && len(args.Data) == 0 {
					noErr = false
					errMsg.Error = "Don't know what return..."
					success = false
				} else {
					args.Code = val.HttpStatus
					args.MimeType = val.MimeType
					args.Data = bytes.NewBufferString(val.Data).Bytes()
				}
			}
		}

		if val.PostForm != nil {
			if checkPostForm(val.PostForm, ctx) {
				if _, ok := val.PostForm["post_data"]; ok {
					args.Data = bytes.NewBufferString(val.PostForm["post_data"]).Bytes()
					success = true
					noErr = true
				} else {
					if val.Data != "" && len(args.Data) == 0 {
						args.Data = bytes.NewBufferString(val.Data).Bytes()
					} else {
						noErr = false
						success = false
						errMsg.Error = "Post Form doesn't equal with form from config..."
					}
				}
				args.Code = val.HttpStatus
				args.MimeType = val.MimeType
			} else {
				if val.Data == "" {
					noErr = false
					success = false
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

		if success {
			break
		} else {
			success = true
			noErr = true
		}
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
