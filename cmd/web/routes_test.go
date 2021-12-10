package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing test passed
	default:
		t.Error(fmt.Sprintf("test failed, type is %T", v))
	}
}
