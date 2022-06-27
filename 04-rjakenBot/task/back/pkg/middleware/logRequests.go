package middleware

import (
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func LogRequests(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(httpCtx *fasthttp.RequestCtx) {
		method := string(httpCtx.Method())
		path := string(httpCtx.Path())

		log.WithFields(log.Fields{
			"method": method,
			"path":   path,
		}).Infof("[%s] %s", method, path)

		handler(httpCtx)
	}
}
