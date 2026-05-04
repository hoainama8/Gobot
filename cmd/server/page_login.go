package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"net/http"
)

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	userID := ""                                          // [THEM_MOI]
	if c, err := r.Cookie("session_cookie"); err == nil { // [THEM_MOI]
		if accounts, loadErr := service.LoadAccounts(); loadErr == nil { // [THEM_MOI]
			if account, ok := service.FindBySessionCookie(c.Value, accounts); ok { // [THEM_MOI]
				userID = account.Webchat // [THEM_MOI]
			}
		}
	}
	data := &templateData{ // [SUA]
		Form:      userLoginForm{}, // [THEM_MOI]
		UserID:    userID,          // [THEM_MOI]
		IsAuthen:  false, // [SUA]
	}
	app.render(w, http.StatusOK, "login.html", data) // [SUA]
} // [THEM_MOI]
