package agent

type AIStructuredResponse struct {
	Text             string    `json:"text"`
	Photourl         []string  `json:"photourl"`
	Videourl         []string  `json:"videourl"`
	IsOrderCompleted bool      `json:"is_order_completed"`
	OrderInfo        OrderInfo `json:"order_info"`
	PaymentQRURL     string    `json:"payment_qr_url"` // [THEM_MOI]
	OrderSaved       bool      `json:"-"`              // [THEM_MOI]
}

func (a *Agent) Response_format() any { // [SUA]
	return map[string]any{ // [THEM_MOI]
		"type": "object",
		"properties": map[string]any{
			"text": map[string]any{
				"type": "string",
			},
			"photourl": map[string]any{
				"type":  "array",
				"items": map[string]any{"type": "string"},
			},
			"videourl": map[string]any{
				"type":  "array",
				"items": map[string]any{"type": "string"},
			},
		},
		"required": []string{"text", "photourl", "videourl"}, // [THEM_MOI]
	}
}
