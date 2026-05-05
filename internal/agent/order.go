package agent

import (
	"encoding/json"
	"fmt"
	"net/url" // [THEM_MOI]
	"os"
	"path/filepath"
	"time"
)

// OrderInfo đại diện cho thông tin đơn hàng
type OrderInfo struct {
	ID_order     string   `json:"id_order"`
	TableNumber  *int     `json:"table_number"`
	DiningOption *string  `json:"dining_option"`
	Items        []item   `json:"items"`
	From         customer `json:"customer"`
	Status       string   `json:"status"`
	Pay          pay      `json:"pay"`
	Note         *string  `json:"note"`
	PaymentQRURL string   `json:"payment_qr_url"` // [THEM_MOI]
}
type pay struct {
	Status string  `json:"status"`
	Amount *int    `json:"amount"`
	PayID  *string `json:"pay_id"`
}

type item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    *int   `json:"price"`
}

type customer struct {
	Webchat string  `json:"webchat"`
	Name    *string `json:"name"`
	Phone   *string `json:"phone"`
}

type PaymentQRBuilder struct{} // [THEM_MOI]

func NewOrder() *OrderInfo {
	return &OrderInfo{}
}

func (a *Agent) CreateOrderTool() any { // [THEM_MOI]
	return map[string]any{
		"type":        "function",
		"name":        "create_order",
		"description": "Tao don hang moi sau khi khach hang da xac nhan va da cung cap du thong tin bat buoc.",
		"parameters": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"table_number": map[string]any{
					"type":        []any{"integer", "null"},
					"description": "So ban neu khach dung tai cho.",
				},
				"dining_option": map[string]any{
					"type":        "string",
					"enum":        []any{"tai cho", "mang ve"},
					"description": "Hinh thuc dung mon.",
				},
				"items": map[string]any{
					"type":        "array",
					"description": "Danh sach mon khach dat.",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"name":     map[string]any{"type": "string"},
							"quantity": map[string]any{"type": "integer"},
							"price":    map[string]any{"type": []any{"integer", "null"}},
						},
						"required": []string{"name", "quantity"},
					},
				},
				"from": map[string]any{
					"type":        "object",
					"description": "Thong tin khach hang.",
					"properties": map[string]any{
						"webchat": map[string]any{"type": []any{"string", "null"}},
						"name":    map[string]any{"type": []any{"string", "null"}},
						"phone":   map[string]any{"type": []any{"string", "null"}},
					},
				},
				"note": map[string]any{
					"type":        []any{"string", "null"},
					"description": "Ghi chu cua khach neu co.",
				},
				"pay": map[string]any{ // [THEM_MOI]
					"type":        "object",
					"description": "Thong tin thanh toan cua don hang.",
					"properties": map[string]any{
						"status": map[string]any{"type": "string"},
						"amount": map[string]any{"type": []any{"integer", "null"}},
						"pay_id": map[string]any{"type": []any{"string", "null"}},
					},
				},
			},
			"required": []string{"dining_option", "items"},
		},
	}
}

func (a *Agent) ExecuteCreateOrder(arguments json.RawMessage, userID string) (map[string]any, OrderInfo, error) { // [THEM_MOI]
	order, err := ParseCreateOrderArguments(arguments)
	if err != nil {
		return nil, OrderInfo{}, err
	}
	if order.From.Webchat == "" {
		order.From.Webchat = userID
	}
	order.EnsureID()
	if order.Status == "" {
		order.Status = "new"
	}
	order.NormalizePay()                                  // [THEM_MOI]
	order.PaymentQRURL = NewPaymentQRBuilder().BuildURL(order) // [SUA]
	if err := NewOrder().Add_order(order); err != nil {
		return nil, OrderInfo{}, err
	}
	result := map[string]any{
		"success":        true,
		"message":        "Don hang da duoc tao thanh cong.",
		"order_info":     order,
		"payment_qr_url": order.PaymentQRURL,
	}
	return result, order, nil
}

func ParseCreateOrderArguments(arguments json.RawMessage) (OrderInfo, error) { // [THEM_MOI]
	var order OrderInfo
	if len(arguments) == 0 {
		return order, fmt.Errorf("function create_order thieu arguments")
	}
	if err := json.Unmarshal(arguments, &order); err == nil {
		return order, nil
	}
	var argText string
	if err := json.Unmarshal(arguments, &argText); err != nil {
		return order, fmt.Errorf("khong parse duoc arguments cua create_order: %v", err)
	}
	if err := json.Unmarshal([]byte(argText), &order); err != nil {
		return order, fmt.Errorf("khong parse duoc JSON arguments cua create_order: %v", err)
	}
	return order, nil
}

func (a *OrderInfo) EnsureID() { // [THEM_MOI]
	if a.ID_order == "" {
		a.ID_order = fmt.Sprintf("OD-%d", time.Now().UnixNano())
	}
}

func (a OrderInfo) TotalAmount() int { // [THEM_MOI]
	total := 0
	for _, item := range a.Items {
		if item.Price != nil && item.Quantity > 0 {
			total += *item.Price * item.Quantity
		}
	}
	return total
}

func (a *OrderInfo) NormalizePay() { // [THEM_MOI]
	if a.Pay.Amount == nil {
		total := a.TotalAmount()
		a.Pay.Amount = &total
	}
	if a.Pay.Status == "" {
		a.Pay.Status = "unpaid"
	}
}

func (a OrderInfo) PaymentAmount() int { // [THEM_MOI]
	if a.Pay.Amount != nil && *a.Pay.Amount > 0 {
		return *a.Pay.Amount
	}
	return a.TotalAmount()
}

func NewPaymentQRBuilder() *PaymentQRBuilder { // [THEM_MOI]
	return &PaymentQRBuilder{}
}

func (b *PaymentQRBuilder) BuildURL(order OrderInfo) string { // [THEM_MOI]
	orderID := order.ID_order
	if orderID == "" {
		orderID = fmt.Sprintf("OD-%d", time.Now().UnixNano())
	}
	amount := order.PaymentAmount()
	addInfo := fmt.Sprintf("Thanh toan don hang %s", orderID)
	bankID := os.Getenv("PAYMENT_BANK_ID")
	accountNo := os.Getenv("PAYMENT_ACCOUNT_NO")
	accountName := os.Getenv("PAYMENT_ACCOUNT_NAME")
	if bankID != "" && accountNo != "" {
		q := url.Values{}
		if amount > 0 {
			q.Set("amount", fmt.Sprintf("%d", amount))
		}
		q.Set("addInfo", addInfo)
		if accountName != "" {
			q.Set("accountName", accountName)
		}
		return fmt.Sprintf("https://img.vietqr.io/image/%s-%s-compact2.png?%s", url.PathEscape(bankID), url.PathEscape(accountNo), q.Encode())
	}
	if amount > 0 {
		addInfo = fmt.Sprintf("%s - So tien %d VND", addInfo, amount)
	}
	return fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=240x240&data=%s", url.QueryEscape(addInfo))
}

func (a OrderInfo) PaymentQRCodeURL() string { // [SUA]
	return NewPaymentQRBuilder().BuildURL(a) // [THEM_MOI]
}

func (a *OrderInfo) Add_order(neworder OrderInfo) error {

	dateString := time.Now().Format("2006-01-02")

	filePath := filepath.Join("data", "orders", fmt.Sprintf("%s.json", dateString))
	// Tạo thư mục cha nếu chưa có
	os.MkdirAll(filepath.Dir(filePath), 0755)

	var orders []OrderInfo

	// 1. Nếu file đã tồn tại, đọc dữ liệu cũ ra
	if _, err := os.Stat(filePath); err == nil {
		data, _ := os.ReadFile(filePath)
		json.Unmarshal(data, &orders)
	}

	// 2. Thêm đơn hàng mới vào mảng
	orders = append(orders, neworder)

	// 3. Ghi đè lại toàn bộ mảng vào file
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
