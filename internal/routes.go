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
	v1Mux := chi.NewRouter()
	mux.Mount("/v1", v1Mux)
	v1Mux.Get("/", app.Home)
	v1Mux.Get("/chat", app.ChatHandler)
	v1Mux.Get("/game", app.GameHandler)
	v1Mux.Get("/auth", app.Auth)
	v1Mux.Get("/user", app.UserHandler)
	// v1Mux.Post("/user", app.CreateUserHandler)
	return mux
}