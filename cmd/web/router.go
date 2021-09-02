package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yaji1122/bookings-go/pkg/config"
	"github.com/yaji1122/bookings-go/pkg/handler"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//mux := pat.New()
	//mux.Get("/", http.HandlerFunc(handler.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handler.Repo.About))
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
