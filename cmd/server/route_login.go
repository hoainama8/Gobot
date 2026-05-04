package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"botai/botaigolang/internal/validator"
	"net/http"
	"time" // [THEM_MOI]

	"golang.org/x/crypto/bcrypt" // [THEM_MOI]
)

type userLoginForm struct {
	Phone               string `form:"phone"` // [SUA]
	Passwd              string `form:"passwd"` // [SUA]
	validator.Validator `form:"-"`
}

// Update the handler so it displays the login page.
func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form) // [SUA]
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Phone), "phone", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Passwd), "passwd", "This field cannot be blank")

	accounts, err := service.LoadAccounts() // [THEM_MOI]
	if err != nil {
		app.serverError(w, err)
		return
	}
	if form.Valid() {
		account, exists := service.FindByPhone(form.Phone, accounts) // [SUA]
		if !exists || bcrypt.CompareHashAndPassword([]byte(account.Passwd), []byte(form.Passwd)) != nil { // [SUA]
			form.AddNonFieldError("Phone or password is incorrect") // [SUA]
		} else {
			form.Phone = account.Phone // [THEM_MOI]
			if account.SessionCookie == "" { // [THEM_MOI]
				form.AddNonFieldError("Invalid account session, please signup again")
				data := &templateData{
					Form: form,
				}
				app.render(w, http.StatusUnprocessableEntity, "login.html", data)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "session_cookie",
				Value:    account.SessionCookie, // [SUA]
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(24 * time.Hour),
			})
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}
	}
	if !form.Valid() {
		data := &templateData{ // [THEM_MOI]
			Form: form, // [THEM_MOI]
		}
		app.render(w, http.StatusUnprocessableEntity, "login.html", data) // [SUA]
		return
	}
}
