package middleware

import (
	"example.com/letoctf/rjakenbot/pkg/sizes"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"runtime"
)

func getStackTrace() string {
	buf := make([]byte, 4*sizes.KB)
	written := runtime.Stack(buf, false)

	return string(buf[:written])
}

func Recoverer(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(httpCtx *fasthttp.RequestCtx) {
		defer func() {
			rvr := recover()
			if rvr == nil {
				return
			}

			method := string(httpCtx.Method())
			path := string(httpCtx.Path())

			log.WithFields(log.Fields{
				"method":     method,
				"path":       path,
				"stacktrace": getStackTrace(),
			}).Errorf("panic occurred while handling %s request at path %s: %s", method, path, rvr)

			httpCtx.SetStatusCode(fasthttp.StatusInternalServerError)
		}()

		handler(httpCtx)
	}
}
