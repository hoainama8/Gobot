package main

import (
	"botai/botaigolang/internal/agent" // [THEM_MOI]
	"encoding/json"
	"fmt"
	"html"
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
		"text":     aiResp.Text,
		"photourl": aiResp.Photourl,
		"videourl": aiResp.Videourl,
	}
	if aiResp.IsOrderCompleted { // [THEM_MOI]
		order := aiResp.OrderInfo
		order.EnsureID() // [SUA]
		if order.Status == "" {
			order.Status = "new" // [THEM_MOI]
		}
		order.NormalizePay()                                     // [THEM_MOI]
		order.PaymentQRURL = agent.NewPaymentQRBuilder().BuildURL(order) // [SUA]
		aiResp.OrderInfo = order                      // [THEM_MOI]
		aiResp.PaymentQRURL = order.PaymentQRURL      // [THEM_MOI]
		if !aiResp.OrderSaved {                       // [THEM_MOI]
			if err := app.order.Add_order(order); err != nil {
				http.Error(w, fmt.Sprintf("[THEM_MOI] save order failed: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}

	// [THEM_MOI] Neu goi tu HTMX thi tra fragment HTML de render truc tiep len page.
	if r.Header.Get("HX-Request") == "true" { // [THEM_MOI]
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		escapedUserMsg := html.EscapeString(message) // [THEM_MOI]
		if aiResp.IsOrderCompleted {
			fmt.Fprintf(w, `<article class="msg user"><h4>Ban</h4><p>%s</p></article><article class="msg ok"><h4>Don hang da duoc tao</h4><p>Ma don: <strong>%s</strong></p><p>Quet ma QR de thanh toan:</p><img src="%s" alt="QR thanh toan don hang %s" style="width:240px;height:240px;max-width:100%%;border:1px solid #eee;border-radius:8px;padding:8px;background:#fff;"></article>`, escapedUserMsg, html.EscapeString(aiResp.OrderInfo.ID_order), html.EscapeString(aiResp.PaymentQRURL), html.EscapeString(aiResp.OrderInfo.ID_order)) // [SUA]
			return
		}

		var b strings.Builder
		b.WriteString(`<article class="msg user"><h4>Ban</h4><p>` + escapedUserMsg + `</p></article>`) // [THEM_MOI]
		b.WriteString(`<article class="msg"><h4>Phan hoi</h4>`)
		if aiResp.Text != "" {
			b.WriteString(`<p><strong>text:</strong> ` + html.EscapeString(aiResp.Text) + `</p>`)
		}
		if len(aiResp.Photourl) > 0 {
			b.WriteString(`<div class="media-grid">`) // [SUA]
			for _, u := range aiResp.Photourl {
				escaped := html.EscapeString(u)
				b.WriteString(`<figure class="media-item"><img src="` + escaped + `" alt="Hinh anh tu agent" loading="lazy" style="max-width:100%;height:auto;border-radius:8px;border:1px solid #eee;"><figcaption><a href="` + escaped + `" target="_blank" rel="noopener noreferrer">Mo anh</a></figcaption></figure>`) // [SUA]
			}
			b.WriteString(`</div>`) // [SUA]
		}
		if len(aiResp.Videourl) > 0 {
			b.WriteString(`<div class="media-grid">`) // [SUA]
			for _, u := range aiResp.Videourl {
				escaped := html.EscapeString(u)
				b.WriteString(`<figure class="media-item"><video src="` + escaped + `" controls preload="metadata" style="max-width:100%;height:auto;border-radius:8px;border:1px solid #eee;"></video><figcaption><a href="` + escaped + `" target="_blank" rel="noopener noreferrer">Mo video</a></figcaption></figure>`) // [SUA]
			}
			b.WriteString(`</div>`) // [SUA]
		}
		b.WriteString(`</article>`)
		_, _ = w.Write([]byte(b.String()))
		return
	}

	// [SUA] Luon chuyen response thanh JSON truoc khi tra ve client.
	w.Header().Set("Content-Type", "application/json")         // [SUA]
	w.WriteHeader(http.StatusOK)                               // [SUA]
	if err := json.NewEncoder(w).Encode(payload); err != nil { // [SUA]
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
