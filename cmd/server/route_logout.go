package main

import (
	"net/http" // [THEM_MOI]
	"time" // [THEM_MOI]
)

func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{ // [THEM_MOI]
		Name:     "session_cookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther) // [THEM_MOI]
}
