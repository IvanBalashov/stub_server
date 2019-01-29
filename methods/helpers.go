package methods

import (
	"github.com/gin-gonic/gin"
)

func setHeaders(headers map[string]string, ctx *gin.Context) {
	if headers == nil {
		return
	}

	for key, val := range headers {
		ctx.Header(key, val)
	}
}

func setCookies(cookies map[string][]Cookie, ctx *gin.Context) {
	if cookies == nil {
		return
	}

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

func uploadFile(path string, ctx *gin.Context) {
	ctx.File(path)
}

func checkQueries(queries map[string]string, ctx *gin.Context) bool {
	for key, val := range queries {
		if key == "query_data" {
			continue
		}
		if ctx.Query(key) != val {
			return false
		}
	}
	return true
}

func checkHeaders(headers map[string]string, ctx *gin.Context) bool {
	for key, val := range headers {
		if ctx.GetHeader(key) != val {
			return false
		}
	}
	return true
}

func checkPostForm(form map[string]string, ctx *gin.Context) bool {
	for key, val := range form {
		if ctx.PostForm(key) != val {
			return false
		}
	}
	return true
}
