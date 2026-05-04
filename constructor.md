app/
├── cmd/
│ └── server/
│ └── main.go # Điểm khởi đầu (Entry point) của ứng dụng
├── internal/ # Mã nguồn riêng tư (Không thể import từ ngoài)
│ ├── auth/ # Xử lý OAuth2, Session, Login
│ ├── agent/ # Logic điều khiển Gemini (Prompt, Config)
│ ├── web/ # Các Handler (Route), Middleware
│ └── service/ # Logic nghiệp vụ (Business logic)
├── pkg/ # Thư viện dùng chung (Có thể tái sử dụng)
│ └── utils/ # Các hàm hỗ trợ tóm tắt text, lọc HTML
├── api/ # Định nghĩa API (OpenAPI/Swagger nếu cần)
├── configs/ # File cấu hình (.env, config.yaml)
├── web/ # Frontend (HTML, CSS, JS, React/Vue)
├── go.mod # Quản lý dependencies
├── go.sum
└── .env # Lưu GEMINI_API_KEY, CLIENT_ID (Không commit lên Git)
