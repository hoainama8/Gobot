package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"net/http"
)

func (app *Application) ChatPage(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	if r.URL.Path != "/chat" {
		app.notFound(w)
		return
	}
	userID := "" // [THEM_MOI]
	if c, err := r.Cookie("session_cookie"); err == nil { // [THEM_MOI]
		if accounts, loadErr := service.LoadAccounts(); loadErr == nil { // [THEM_MOI]
			if account, ok := service.FindBySessionCookie(c.Value, accounts); ok { // [THEM_MOI]
				userID = account.Webchat // [THEM_MOI]
			}
		}
	}
	data := &templateData{ // [THEM_MOI]
		UserID: userID, // [THEM_MOI]
	}
	app.render(w, http.StatusOK, "chat.html", data) // [SUA]
}
