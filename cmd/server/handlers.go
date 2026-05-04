package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter" // [THEM_MOI]

	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	// Thiết lập handler tùy chỉnh cho lỗi 404
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Route cho các tệp tĩnh
	fileServer := http.FileServer(http.Dir("./web/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodGet, "/signup", app.Signup) // [THEM_MOI]
	router.HandlerFunc(http.MethodPost, "/signup", app.userSignup) // [THEM_MOI]
	router.HandlerFunc(http.MethodGet, "/login", app.Login)
	router.HandlerFunc(http.MethodPost, "/login", app.userLogin)
	router.HandlerFunc(http.MethodPost, "/logout", app.userLogout) // [THEM_MOI]
	router.HandlerFunc(http.MethodGet, "/account", app.AccountPage) // [THEM_MOI]
	router.HandlerFunc(http.MethodGet, "/orders", app.OrdersPage) // [THEM_MOI]
	router.HandlerFunc(http.MethodGet, "/chat", app.ChatPage)
	router.HandlerFunc(http.MethodPost, "/chat", func(w http.ResponseWriter, r *http.Request) { app.chat(w, r) }) // [SUA]
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
