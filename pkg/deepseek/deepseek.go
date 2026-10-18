package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const chatTimeout = 60 * time.Second

type Message struct {
	Role  string
	Text  string
	Image string
}

func (m Message) MarshalJSON() ([]byte, error) {
	if m.Image == "" {
		return json.Marshal(struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{m.Role, m.Text})
	}
	parts := make([]map[string]any, 0, 2)
	if m.Text != "" {
		parts = append(parts, map[string]any{"type": "text", "text": m.Text})
	}
	parts = append(parts, map[string]any{"type": "image_url", "image_url": map[string]string{"url": m.Image}})
	return json.Marshal(struct {
		Role    string           `json:"role"`
		Content []map[string]any `json:"content"`
	}{m.Role, parts})
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var (
	apiKey  string
	baseURL = "https://api.deepseek.com"
	model   = "deepseek-v4-flash"
	httpCl  = &http.Client{Timeout: chatTimeout}
)

func Init(key, url, mdl string) {
	apiKey = key
	if url != "" {
		baseURL = url
	}
	if mdl != "" {
		model = mdl
	}
}

func Enabled() bool {
	return apiKey != ""
}

func Chat(messages []Message) (string, error) {
	if !Enabled() {
		return "", errors.New("DeepSeek API 未配置")
	}
	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: messages,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, baseURL+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := httpCl.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API 请求失败: %d %s", resp.StatusCode, string(data))
	}
	var result chatResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", errors.New("DeepSeek 返回空结果")
	}
	return extractContent(result.Choices[0].Message.Content)
}

func extractContent(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &parts); err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, p := range parts {
		if p.Type == "text" {
			sb.WriteString(p.Text)
		}
	}
	return sb.String(), nil
}
