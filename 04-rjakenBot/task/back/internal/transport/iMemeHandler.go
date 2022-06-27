package transport

import "github.com/valyala/fasthttp"

type IMemeHandler interface {
	GenerateMeme(httpCtx *fasthttp.RequestCtx)
}
