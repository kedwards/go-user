package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/usermgmt/health-status", app.GetHealthStatus)
	mux.Get("/usermgmt/users", app.GetAllUsers)
	mux.Post("/usermgmt/user", app.CreateUser)
	mux.Put("/usermgmt/user", app.UpdateUser)
	mux.Delete("/usermgmt/user/{user}", app.DeleteUser)

	return mux
}
