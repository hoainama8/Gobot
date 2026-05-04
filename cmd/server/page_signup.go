package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"net/http"                           // [THEM_MOI]
)

func (app *Application) Signup(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	userID := ""                                          // [THEM_MOI]
	if c, err := r.Cookie("session_cookie"); err == nil { // [THEM_MOI]
		if accounts, loadErr := service.LoadAccounts(); loadErr == nil { // [THEM_MOI]
			if account, ok := service.FindBySessionCookie(c.Value, accounts); ok { // [THEM_MOI]
				userID = account.Webchat // [THEM_MOI]
			}
		}
	}
	data := &templateData{ // [THEM_MOI]
		Form:    userSignupForm{},  // [THEM_MOI]
		Account: service.Account{}, // [THEM_MOI]
		UserID:  userID,            // [THEM_MOI]

	}
	app.render(w, http.StatusOK, "signup.html", data) // [THEM_MOI]
}
