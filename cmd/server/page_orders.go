package main

import (
	"botai/botaigolang/internal/agent" // [THEM_MOI]
	"encoding/json" // [THEM_MOI]
	"errors" // [THEM_MOI]
	"net/http" // [THEM_MOI]
	"os" // [THEM_MOI]
)

type orderHistoryView struct { // [THEM_MOI]
	IDOrder string
	Status  string
	Count   int
}

func loadOrderHistory() ([]orderHistoryView, error) { // [THEM_MOI]
	const orderFile = "data/order/order.json"
	data, err := os.ReadFile(orderFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []orderHistoryView{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []orderHistoryView{}, nil
	}
	var orders []agent.OrderInfo
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}
	views := make([]orderHistoryView, 0, len(orders))
	for _, o := range orders {
		views = append(views, orderHistoryView{
			IDOrder: o.ID_order,
			Status:  o.Status,
			Count:   len(o.Items),
		})
	}
	return views, nil
}

func (app *Application) OrdersPage(w http.ResponseWriter, r *http.Request) { // [THEM_MOI]
	if r.URL.Path != "/orders" {
		app.notFound(w)
		return
	}

	userID := ""
	if c, err := r.Cookie("session_cookie"); err == nil {
		if a, ok := app.currentAccount(c.Value); ok {
			userID = a.Webchat
		}
	}
	if userID == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	orders, err := loadOrderHistory()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{ // [THEM_MOI]
		UserID: userID,
		Orders: orders,
	}
	app.render(w, http.StatusOK, "orders.html", data)
}
