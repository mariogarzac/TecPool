package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
    mux := chi.NewRouter()

    // middleware
    mux.Use(middleware.Recoverer)
    mux.Use(SessionLoad)

    // routes
    mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
    mux.Post("/register", http.HandlerFunc(handlers.Repo.PostRegister))

    mux.Get("/dashboard", http.HandlerFunc(handlers.Repo.Dashboard))
    mux.Get("/login", http.HandlerFunc(handlers.Repo.Login))
    mux.Post("/login", http.HandlerFunc(handlers.Repo.PostLogin))

    mux.With(IsLoggedIn).Get("/create-trip", http.HandlerFunc(handlers.Repo.CreateTrip))
    mux.With(IsLoggedIn).Post("/create-trip", http.HandlerFunc(handlers.Repo.PostCreateTrip))

    // fileServer := http.FileServer(http.Dir("./static/"))
    // mux.Handle("/static/*", http.StripPrefix("/static",fileServer))

    return mux
}
