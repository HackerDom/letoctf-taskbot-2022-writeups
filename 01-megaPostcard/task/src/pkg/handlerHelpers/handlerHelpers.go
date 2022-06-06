package handlerHelpers

import (
	"html/template"
	"net/http"
)

type ErrorResponseInfo struct {
	Message    string
	StatusCode int
}

type TemplateHandler func(r *http.Request) (tmpl *template.Template, data any, errResp *ErrorResponseInfo)

func TemplateResponse(templateHandler TemplateHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, errResp := templateHandler(r)
		if errResp != nil {
			http.Error(w, errResp.Message, errResp.StatusCode)
			return
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}
