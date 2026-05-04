package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Agent struct {
	config  *ConfigFile
	History []Message
	APIKey  string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type InteractionRequest struct {
	Model                   string     `json:"model"`
	Input                   any        `json:"input"`           // [SUA]
	Tools                   []any      `json:"tools,omitempty"` // [THEM_MOI]
	Response_format         any        `json:"response_format"` // [SUA]
	System_instruction      string     `json:"system_instruction"`
	Generation_config       ConfigFile `json:"generation_config"`
	Previous_interaction_id *string    `json:"previous_interaction_id,omitempty"`
}

type interactionResponse struct { // [THEM_MOI]
	ID      string              `json:"id"` // [THEM_MOI]
	Outputs []interactionOutput `json:"outputs"`
}

type interactionOutput struct { // [THEM_MOI]
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Text      string          `json:"text"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func (r *interactionResponse) FirstFunctionCall(name string) (interactionOutput, bool) { // [THEM_MOI]
	if r == nil {
		return interactionOutput{}, false
	}
	for _, output := range r.Outputs {
		if output.Type == "function_call" && output.Name == name {
			return output, true
		}
	}
	return interactionOutput{}, false
}

func NewAgent(apikey string) *Agent {
	return &Agent{
		config:  &ConfigFile{},
		History: []Message{},
		APIKey:  apikey,
	}
}

// LoadConfig doc file cau hinh agent (vi du: config.json)
func (a *Agent) LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("khong the doc file config: %v", err)
	}
	return json.Unmarshal(data, &a.config)
}

func (a *Agent) Chat(ctx context.Context, userInput string, userId string) (*AIStructuredResponse, error) { // [SUA]
	// [THEM_MOI] Defensive check: do not call remote API without key.
	if a.APIKey == "" {
		return nil, fmt.Errorf("GEMINIKEY is required") // [SUA]
	}
	// 1. Them cau hoi cua user vao lich su noi bo
	a.History = append(a.History, Message{Role: "user", Content: userInput})

	// 3. Chuan bi Request (su dung struct InteractionRequest da co o tren)
	reqBody := InteractionRequest{
		Model:                   "gemini-3-flash-preview",
		Input:                   userInput,
		Tools:                   []any{a.CreateOrderTool()}, // [THEM_MOI]
		Response_format:         a.Response_format(),
		System_instruction:      a.SystemPrompt(),
		Generation_config:       a.Generation_config(),
		Previous_interaction_id: a.Prev_ID(),
	}
	body, wrapped, err := a.sendInteraction(ctx, reqBody) // [SUA]
	if err != nil {
		return nil, err // [SUA]
	}

	// 4. Luu phan hoi cua AI vao lich su
	a.History = append(a.History, Message{Role: "model", Content: string(body)})

	var toolOrder OrderInfo                                        // [THEM_MOI]
	toolCalled := false                                            // [THEM_MOI]
	if call, ok := wrapped.FirstFunctionCall("create_order"); ok { // [THEM_MOI]
		toolResult, createdOrder, err := a.ExecuteCreateOrder(call.Arguments, userId) // [THEM_MOI]
		if err != nil {
			return nil, err // [THEM_MOI]
		}
		toolCalled = true                // [THEM_MOI]
		toolOrder = createdOrder         // [THEM_MOI]
		resultInput := []map[string]any{ // [THEM_MOI]
			{
				"type":    "function_result",
				"name":    call.Name,
				"call_id": call.ID,
				"result":  toolResult,
			},
		}
		nextReq := InteractionRequest{ // [THEM_MOI]
			Model:                   "gemini-3-flash-preview",
			Input:                   resultInput,
			Response_format:         a.Response_format(),
			System_instruction:      a.SystemPrompt(),
			Generation_config:       a.Generation_config(),
			Previous_interaction_id: &wrapped.ID,
		}
		body, _, err = a.sendInteraction(ctx, nextReq) // [THEM_MOI]
		if err != nil {
			return nil, err // [THEM_MOI]
		}
		a.History = append(a.History, Message{Role: "model", Content: string(body)}) // [THEM_MOI]
	}

	result, err := a.ParseAIResponse(string(body))
	if err != nil {
		return nil, fmt.Errorf("AI tra ve sai dinh dang JSON: %v", err) // [SUA]
	}
	if toolCalled { // [THEM_MOI]
		result.OrderSaved = true       // [THEM_MOI]
		result.IsOrderCompleted = true // [THEM_MOI]
		result.OrderInfo = toolOrder   // [SUA]
		if result.PaymentQRURL == "" {
			result.PaymentQRURL = toolOrder.PaymentQRURL // [THEM_MOI]
		}
	}
	// 5. Tu dong luu lich su xuong file sau moi lan chat
	if err := a.SaveHistory(userId); err != nil { // [SUA]
		return nil, fmt.Errorf("save history failed: %v", err) // [SUA]
	}
	return result, nil
}

func (a *Agent) sendInteraction(ctx context.Context, reqBody InteractionRequest) ([]byte, *interactionResponse, error) { // [THEM_MOI]
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("loi ma hoa request: %v", err)
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/interactions?key=%s", a.APIKey)
	reqCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 65 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API Error (%d): %s", resp.StatusCode, string(body))
	}
	var wrapped interactionResponse
	if err := json.Unmarshal(body, &wrapped); err != nil {
		return nil, nil, fmt.Errorf("khong parse duoc response wrapper: %v", err)
	}
	return body, &wrapped, nil
}

// ParseAIResponse chuyen doi chuoi JSON tu AI thanh Struct de xu ly logic
func (a *Agent) ParseAIResponse(rawJson string) (*AIStructuredResponse, error) {
	var wrapped interactionResponse // [SUA]
	if err := json.Unmarshal([]byte(rawJson), &wrapped); err != nil {
		return nil, fmt.Errorf("khong parse duoc response wrapper: %v", err) // [SUA]
	}

	var inner string                      // [THEM_MOI]
	for _, out := range wrapped.Outputs { // [SUA]
		if strings.TrimSpace(out.Text) != "" {
			inner = out.Text
			break
		}
	}
	if inner == "" {
		return nil, fmt.Errorf("khong tim thay outputs[].text")
	}

	inner = strings.TrimSpace(inner) // [THEM_MOI]
	inner = strings.TrimPrefix(inner, "```json")
	inner = strings.TrimPrefix(inner, "```")
	inner = strings.TrimSuffix(inner, "```")
	inner = strings.TrimSpace(inner)

	var payload map[string]any // [SUA]
	if err := json.Unmarshal([]byte(inner), &payload); err != nil {
		return nil, fmt.Errorf("khong parse duoc outputs[].text thanh object: %v", err)
	}

	result := &AIStructuredResponse{} // [THEM_MOI]
	if v, ok := payload["text"].(string); ok {
		result.Text = v
	}
	if v, ok := payload["photourl"].([]any); ok {
		result.Photourl = make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result.Photourl = append(result.Photourl, s)
			}
		}
	}
	if v, ok := payload["videourl"].([]any); ok {
		result.Videourl = make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result.Videourl = append(result.Videourl, s)
			}
		}
	}

	return result, nil
}
