package main

import "botai/botaigolang/internal/service" // [THEM_MOI]

func (app *Application) currentAccount(sessionCookie string) (service.Account, bool) { // [THEM_MOI]
	accounts, err := service.LoadAccounts()
	if err != nil {
		return service.Account{}, false
	}
	return service.FindBySessionCookie(sessionCookie, accounts)
}
