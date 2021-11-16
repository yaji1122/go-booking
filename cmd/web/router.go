package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/handler"
	"net/http"
)

func routes(app *config.Configuration) http.Handler {
	//mux := pat.New()
	//mux.Get("/", http.HandlerFunc(handler.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handler.Repo.About))
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(WriteToConsole)
	mux.Get("/contact", handler.Contact)
	mux.Get("/room/*", handler.Room)
	mux.Get("/reservation", handler.Reservation)
	mux.Get("/", handler.Index)
	mux.Get("/summary", handler.Summary)

	mux.Post("/availability", handler.Availability)
	mux.Post("/reservation", handler.PostReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
