package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/handlers"
	"golang.org/x/net/websocket"
)

func routes(app *config.AppConfig, s *handlers.Hub) http.Handler {
    mux := chi.NewRouter()

    // middleware
    mux.Use(middleware.Recoverer)
    mux.Use(SessionLoad)

    // Serve static files
    fileServer := http.FileServer(http.Dir("templates/static"))
    mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

    // routes
    mux.Get("/register", http.HandlerFunc(handlers.Repo.Register))
    mux.Post("/register", http.HandlerFunc(handlers.Repo.PostRegister))

    mux.Get("/login", http.HandlerFunc(handlers.Repo.Login))
    mux.Post("/login", http.HandlerFunc(handlers.Repo.PostLogin))

    mux.With(IsLoggedIn).Get("/logout", http.HandlerFunc(handlers.Repo.Logout))

    mux.With(IsLoggedIn).Get("/", http.HandlerFunc(handlers.Repo.Dashboard))

    mux.With(IsLoggedIn).Get("/create-trip", http.HandlerFunc(handlers.Repo.CreateTrip))
    mux.With(IsLoggedIn).Post("/create-trip", http.HandlerFunc(handlers.Repo.PostCreateTrip))

    mux.With(IsLoggedIn).Get("/trips", http.HandlerFunc(handlers.Repo.ActiveTrips))

    // Add the new route for searching trips by departure_time
    mux.With(IsLoggedIn).Post("/searchTrips", http.HandlerFunc(handlers.Repo.SearchTripsHandler))

    mux.With(IsLoggedIn).Get("/settings", http.HandlerFunc(handlers.Repo.ShowSettings))

    mux.With(IsLoggedIn).Get("/join-trip/{tripId}", http.HandlerFunc(handlers.Repo.JoinTrip))

    mux.With(IsLoggedIn).Handle("/ws/chat/{tripId}/{userId}", websocket.Handler(s.HandleWs))
    mux.With(IsLoggedIn).Get("/chat/{tripId}/{userId}", handlers.Repo.RenderChat)

    mux.With(IsLoggedIn).Post("/update-group-name", http.HandlerFunc(handlers.Repo.UpdateChatName))

    return mux
}
