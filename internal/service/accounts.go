package service

import (
	"crypto/rand"        // [THEM_MOI]
	"encoding/hex"       // [THEM_MOI]
	"encoding/json"      // [THEM_MOI]
	"errors"             // [THEM_MOI]
	"fmt"                // [THEM_MOI]
	mathrand "math/rand" // [SUA]
	"os"                 // [THEM_MOI]
	"path/filepath"      // [THEM_MOI]
	"time"               // [THEM_MOI]
)

type Account struct {
	Name          string `json:"name"`           // [SUA]
	Phone         string `json:"phone"`          // [SUA]
	Passwd        string `json:"passwd"`         // [SUA]
	Webchat       string `json:"webchat"`        // [SUA]
	Role          string `json:"role"`           // [THEM_MOI]
	SessionCookie string `json:"session_cookie"` // [THEM_MOI]
}

const AccountsFilePath = "account/accounts.json" // [THEM_MOI]

func LoadAccounts() ([]Account, error) { // [THEM_MOI]
	data, err := os.ReadFile(AccountsFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Account{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Account{}, nil
	}
	var accounts []Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

func SaveAccounts(accounts []Account) error { // [THEM_MOI]
	if err := os.MkdirAll(filepath.Dir(AccountsFilePath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(AccountsFilePath, data, 0o644)
}

func FindByWebchat(webchat string, accounts []Account) (Account, bool) { // [THEM_MOI]
	for _, a := range accounts {
		if a.Webchat == webchat {
			return a, true
		}
	}
	return Account{}, false
}

func FindByPhone(phone string, accounts []Account) (Account, bool) { // [THEM_MOI]
	for _, a := range accounts {
		if a.Phone == phone {
			return a, true
		}
	}
	return Account{}, false
}

func FindBySessionCookie(sessionCookie string, accounts []Account) (Account, bool) { // [THEM_MOI]
	for _, a := range accounts {
		if a.SessionCookie == sessionCookie {
			return a, true
		}
	}
	return Account{}, false
}

func NewWebchatID(accounts []Account) string { // [THEM_MOI]
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano())) // [SUA]
	for {
		id := fmt.Sprintf("wc%06d", r.Intn(1000000))
		if _, exists := FindByWebchat(id, accounts); !exists {
			return id
		}
	}
}

func NewSessionCookieValue() (string, error) { // [THEM_MOI]
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
