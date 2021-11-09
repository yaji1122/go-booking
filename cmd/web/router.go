package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/handler"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//mux := pat.New()
	//mux.Get("/", http.HandlerFunc(handler.Repo.Home))
	//mux.Get("/about", http.HandlerFunc(handler.Repo.About))
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(WriteToConsole)
	mux.Get("/booking", handler.Repo.Booking)
	mux.Get("/contact", handler.Repo.Contact)
	mux.Get("/room", handler.Repo.Room)
	mux.Get("/reservation", handler.Repo.Reservation)
	mux.Get("/", handler.Repo.Index)
	mux.Get("/summary", handler.Repo.Summary)

	mux.Post("/search-availability", handler.Repo.Availability)
	mux.Post("/reservation", handler.Repo.PostReservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
