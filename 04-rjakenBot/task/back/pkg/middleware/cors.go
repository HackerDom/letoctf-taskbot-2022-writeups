package middleware

import (
	"github.com/valyala/fasthttp"
)

func CorsAllowAll(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(httpCtx *fasthttp.RequestCtx) {
		handler(httpCtx)
		httpCtx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		httpCtx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	}
}
