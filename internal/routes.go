package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)	

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8081", "http://192.168.1.100:8081", "http://121.45.87.60:8081", "exp://192.168.1.100:8081"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300, // Maximum value not ignored by any of major browsers
		}))
	mux.Get("/", app.Home)
	mux.Get("/chat", app.ChatHandler)
	mux.Get("/game", app.GameHandler)
	mux.Get("/auth", app.Auth)
	return mux
}