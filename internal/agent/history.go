package agent

import (
	"encoding/json"
	"fmt"
	"os"
)

func (a *Agent) SaveHistory(UserID string) error {
	fileName := fmt.Sprintf("history_%s.json", UserID)
	data, err := json.MarshalIndent(a.History, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (a *Agent) LoadHistory(UserID string) error {
	fileName := fmt.Sprintf("history_%s.json", UserID)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil // Không có file thì thôi, bắt đầu lịch sử mới
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.History)
}
