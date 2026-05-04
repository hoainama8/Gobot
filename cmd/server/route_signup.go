package main

import (
	"botai/botaigolang/internal/service" // [THEM_MOI]
	"botai/botaigolang/internal/validator" // [THEM_MOI]
	"net/http" // [THEM_MOI]
	"time" // [THEM_MOI]

	"golang.org/x/crypto/bcrypt" // [THEM_MOI]
)

type userSignupForm struct { // [THEM_MOI]
	Name                string `form:"name"` // [THEM_MOI]
	Phone               string `form:"phone"` // [THEM_MOI]
	Passwd              string `form:"passwd"` // [THEM_MOI]
	ConfirmPassword     string `form:"confirm_password"` // [THEM_MOI]
	validator.Validator `form:"-"` // [THEM_MOI]
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	var form userSignupForm
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Phone), "phone", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Passwd), "passwd", "This field cannot be blank")
	form.CheckField(len(form.Passwd) >= 6, "passwd", "Password must be at least 6 characters")
	form.CheckField(validator.NotBlank(form.ConfirmPassword), "confirm_password", "This field cannot be blank")
	form.CheckField(form.Passwd == form.ConfirmPassword, "confirm_password", "Passwords do not match")

	accounts, err := service.LoadAccounts() // [THEM_MOI]
	if err != nil {
		app.serverError(w, err)
		return
	}
	if _, exists := service.FindByPhone(form.Phone, accounts); exists { // [THEM_MOI]
		form.CheckField(false, "phone", "Phone already exists")
	}

	if !form.Valid() {
		data := &templateData{
			Form: form, // [SUA]
			Account: service.Account{ // [THEM_MOI]
				Name:    form.Name,
				Phone:   form.Phone,
				Passwd:  form.Passwd,
			},
		}
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	webchatID := service.NewWebchatID(accounts) // [THEM_MOI]
	sessionCookieValue, err := service.NewSessionCookieValue() // [THEM_MOI]
	if err != nil {
		app.serverError(w, err)
		return
	}
	hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(form.Passwd), bcrypt.DefaultCost) // [THEM_MOI]
	if err != nil {
		app.serverError(w, err)
		return
	}
	accounts = append(accounts, service.Account{ // [THEM_MOI]
		Name:          form.Name, // [SUA]
		Phone:         form.Phone, // [SUA]
		Passwd:        string(hashedPasswd), // [SUA]
		Webchat:       webchatID, // [SUA]
		Role:          "customer", // [THEM_MOI]
		SessionCookie: sessionCookieValue, // [THEM_MOI]
	})
	if err := service.SaveAccounts(accounts); err != nil { // [THEM_MOI]
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_cookie",
		Value:    sessionCookieValue, // [SUA]
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
