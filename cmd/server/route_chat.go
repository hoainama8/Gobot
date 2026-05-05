package main

import (
	"botai/botaigolang/internal/agent" // [THEM_MOI]
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (app *Application) chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	message := strings.TrimSpace(r.FormValue("message")) // [SUA]
	if message == "" {                                   // [SUA]
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	userID := strings.TrimSpace(r.Header.Get("X-User-ID")) // [THEM_MOI]
	if userID == "" {                                      // [THEM_MOI]
		userID = "anonymous" // [THEM_MOI]
	}
	// [SUA] Goi dung chu ky Chat hien tai: (ctx, userInput, userId).
	aiResp, err := app.agent.Chat(ctx, message, userID) // [SUA]
	if err != nil {
		log.Printf("[THEM_MOI] /chat error: %v", err)
		http.Error(w, fmt.Sprintf("[THEM_MOI] %v", err), http.StatusInternalServerError)
		return
	}

	payload := map[string]any{ // [SUA]
		"text":               aiResp.Text,
		"photourl":           aiResp.Photourl,
		"videourl":           aiResp.Videourl,
		"is_order_completed": aiResp.IsOrderCompleted, // [THEM_MOI]
	}
	if aiResp.IsOrderCompleted { // [THEM_MOI]
		order := aiResp.OrderInfo
		order.EnsureID() // [SUA]
		if order.Status == "" {
			order.Status = "new" // [THEM_MOI]
		}
		order.NormalizePay()                                             // [THEM_MOI]
		order.PaymentQRURL = agent.NewPaymentQRBuilder().BuildURL(order) // [SUA]
		aiResp.OrderInfo = order                                         // [THEM_MOI]
		aiResp.PaymentQRURL = order.PaymentQRURL                         // [THEM_MOI]
		if !aiResp.OrderSaved {                                          // [THEM_MOI]
			if err := app.order.Add_order(order); err != nil {
				http.Error(w, fmt.Sprintf("[THEM_MOI] save order failed: %v", err), http.StatusInternalServerError)
				return
			}
		}
		payload["payment_qr_url"] = aiResp.PaymentQRURL // [THEM_MOI]
		payload["order_id"] = aiResp.OrderInfo.ID_order // [THEM_MOI]
	}

	// [SUA] Luon chuyen response thanh JSON truoc khi tra ve client.
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // [SUA]
	w.WriteHeader(http.StatusOK)                                      // [SUA]
	if err := json.NewEncoder(w).Encode(payload); err != nil {        // [SUA]
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
