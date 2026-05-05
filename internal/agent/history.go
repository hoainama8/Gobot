package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath" // [THEM_MOI]
)

func (a *Agent) SaveHistory(UserID string) error {
	// [THEM_MOI] Luu lich su chat vao data/history thay vi root project.
	dirPath := filepath.Join("data", "history")                                // [THEM_MOI]
	fileName := filepath.Join(dirPath, fmt.Sprintf("history_%s.json", UserID)) // [SUA]
	if err := os.MkdirAll(dirPath, 0755); err != nil {                         // [THEM_MOI]
		return err
	}
	data, err := json.MarshalIndent(a.History, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (a *Agent) LoadHistory(UserID string) error {
	// [THEM_MOI] Doc lich su chat tu data/history de dong bo voi SaveHistory.
	fileName := filepath.Join("data", "history", fmt.Sprintf("history_%s.json", UserID)) // [SUA]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil // [SUA] Khong co file thi bat dau lich su moi
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.History)
}
