package main

import (
	"github.com/justinas/nosurf"
	"log"
	"net/http"
)

// WriteToConsole simplee example
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hit the page.")
		next.ServeHTTP(w, r)
	})
}

// NoSurf for csrf check
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	//it uses cookies to make sure that the token it generates is available on per page basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
