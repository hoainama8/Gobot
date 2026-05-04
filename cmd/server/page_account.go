package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"net/http"                           // [THEM_MOI]
)

func (app *Application) AccountPage(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	if r.URL.Path != "/account" {
		app.notFound(w)
		return
	}

	var account service.Account
	userID := ""
	if c, err := r.Cookie("session_cookie"); err == nil {
		if a, ok := app.currentAccount(c.Value); ok {
			account = a
			userID = a.Webchat
		}
	}
	if userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := &templateData{ // [THEM_MOI]
		Account:  account,
		UserID:   userID,
		IsAuthen: true,
	}
	app.render(w, http.StatusOK, "account.html", data)
}
