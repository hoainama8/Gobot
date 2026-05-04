package agent

import (
	"encoding/json"
)

// Prev_ID trả về con trỏ chuỗi chứa ID hoặc nil (null) nếu không tìm thấy
func (a *Agent) Prev_ID() *string {
	// 1. Kiểm tra nếu lịch sử trống
	if len(a.History) == 0 {
		return nil
	}

	// 2. Duyệt ngược từ cuối danh sách tìm tin nhắn của "model"
	var lastModelContent string
	found := false
	for i := len(a.History) - 1; i >= 0; i-- {
		if a.History[i].Role == "model" {
			lastModelContent = a.History[i].Content
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	// 3. Giải mã JSON
	var tempResp struct {
		ID string `json:"id"`
	}

	// Nếu content không phải JSON hoặc lỗi parse, trả về nil
	if err := json.Unmarshal([]byte(lastModelContent), &tempResp); err != nil {
		return nil
	}

	if tempResp.ID == "" {
		return nil
	}

	// Trả về địa chỉ của biến ID
	return &tempResp.ID
}
