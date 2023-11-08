package main

import (
	"context"
	"net/http"
)

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
    return app.Session.LoadAndSave(next)
}

func NameMiddleWare(r http.Request) {
}
func NameMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract the user's name from the cookie or your session
        name := app.Session.GetString(r.Context(), "name")

        // Set the name to the request context
        ctx := context.WithValue(r.Context(), "user_name", name)

        // Call the next handler with the updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// checks if the user is logged in to prevent them from accessing parts that 
// requier user accounts if they are not 
func IsLoggedIn(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !IsUserLoggedIn(r) {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Returns a boolean if they are logged in
func IsUserLoggedIn(r *http.Request) bool {
    return app.Session.Exists(r.Context(), "isLoggedIn")
}
