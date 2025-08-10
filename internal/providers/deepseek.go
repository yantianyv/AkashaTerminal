package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	
	"AkashaTerminal/pkg/types"
)

type DeepSeekProvider struct {
	config types.APIConfig
}

func NewProvider(config types.APIConfig) (*DeepSeekProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("DeepSeek API key is required")
	}
	
	return &DeepSeekProvider{
		config: config,
	}, nil
}

func (p *DeepSeekProvider) SendRequest(prompt string, state interface{}) (string, error) {
	url := "https://api.deepseek.com/v1/chat/completions"
	if p.config.APIBase != "" {
		url = p.config.APIBase
	}
	
	requestBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []interface{}{
			{"role": "system", "content": "You are a helpful coding assistant"},
			{"role": "user", "content": prompt},
		},
		"max_tokens": p.config.MaxTokens,
	}
	
	if p.config.Version != "" {
		requestBody["version"] = p.config.Version
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API error: %s - %s", resp.Status, string(body))
	}
	
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no completions received from DeepSeek")
	}
	
	return response.Choices[0].Message.Content, nil
}

func (p *DeepSeekProvider) GetName() string {
	return "DeepSeek"
}

func (p *DeepSeekProvider) GetModel() string {
	return p.config.Model
}

func (p *DeepSeekProvider) SupportsFeature(feature string) bool {
	switch feature {
	case "long_context":
		return true
	case "multimodal":
		return p.config.Model == "deepseek-vision"
	default:
		return false
	}
}
