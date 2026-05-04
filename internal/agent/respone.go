package agent

import (
	"encoding/json"
	"os"
)

// OrderItem đại diện cho từng món ăn
type OrderItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price,omitempty"`
}

// Cấu trúc để map với file JSON
type ConfigFile struct {
	Max_output_tokens int     `json:"max_output_tokens"` // [SUA]
	Temperature       float64 `json:"temperature"`
	Thinking_level    string  `json:"thinking_level"`
}

func (a *Agent) Generation_config() ConfigFile {
	// 1. Thiết lập giá trị mặc định ban đầu
	conf := ConfigFile{
		Max_output_tokens: 2048,
		Temperature:       1.0,
		Thinking_level:    "low",
	}

	// 2. Đường dẫn tới file config
	configPath := "config/config.json"

	// 3. Đọc file
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Nếu không đọc được file (file không tồn tại), trả về mặc định luôn
		return conf
	}

	// 4. Parse JSON
	var fileData ConfigFile
	if err := json.Unmarshal(data, &fileData); err != nil {
		// Nếu file lỗi định dạng, trả về mặc định
		return conf
	}

	// 5. Ghi đè giá trị nếu có trong file (kiểm tra khác giá trị zero)
	if fileData.Max_output_tokens > 0 {
		conf.Max_output_tokens = fileData.Max_output_tokens
	}

	// Lưu ý: Với float, nếu bạn muốn temp = 0 là hợp lệ thì cần dùng con trỏ,
	// nhưng theo yêu cầu mặc định là 1 nên ta check > 0.
	if fileData.Temperature > 0 {
		conf.Temperature = fileData.Temperature
	}

	if fileData.Thinking_level != "" {
		conf.Thinking_level = fileData.Thinking_level
	}

	return conf
}
