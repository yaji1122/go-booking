package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yaji1122/bookings-go/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.Configuration
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing test passed
	default:
		t.Error(fmt.Sprintf("test failed, type is %T", v))
	}
}
