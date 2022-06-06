package main

import (
	"example.com/employee/internal/handlers/postcards"
	"example.com/employee/pkg/handlerHelpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"net/http"
	"os"
	"time"
)

func AddMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(1 * time.Second))
}

func AddRoutes(router *chi.Mux) {
	templates := template.Must(
		template.ParseFiles("./static/templates/index.html"),
	)

	flag := os.Getenv("FLAG")
	if flag == "" {
		panic("please, set FLAG environment variable")
	}

	postcardService := postcards.NewService(templates, "MegaPostcard 3000", flag)

	fs := http.FileServer(http.Dir("static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", handlerHelpers.TemplateResponse(postcardService.GetIndex))
	router.Post("/create-postcard", handlerHelpers.TemplateResponse(postcardService.CreatePostcard))
}

func main() {
	router := chi.NewRouter()
	AddMiddlewares(router)
	AddRoutes(router)

	err := http.ListenAndServe(":13337", router)
	if err != nil {
		panic(err)
	}
}
