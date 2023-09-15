package main

import (
    "net/http"
)

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
    return app.Session.LoadAndSave(next)
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
