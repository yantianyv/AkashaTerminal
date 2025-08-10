package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	
	"voidClint/pkg/types"
)

type SiliconFlowProvider struct {
	config types.APIConfig
}

func NewProvider(config types.APIConfig) (*SiliconFlowProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("SiliconFlow API key is required")
	}
	
	return &SiliconFlowProvider{
		config: config,
	}, nil
}

func (p *SiliconFlowProvider) SendRequest(prompt string, state interface{}) (string, error) {
	url := "https://api.siliconflow.com/v1/completions"
	if p.config.APIBase != "" {
		url = p.config.APIBase
	}
	
	requestBody := map[string]interface{}{
		"model":       p.config.Model,
		"prompt":      prompt,
		"max_tokens":  p.config.MaxTokens,
		"temperature": 0.7,
	}
	
	if strings.HasPrefix(p.config.Model, "yi-") {
		requestBody["quant_mode"] = "int8"
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
		return "", fmt.Errorf("SiliconFlow API error: %s - %s", resp.Status, string(body))
	}
	
	var response struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no completions received from SiliconFlow")
	}
	
	return response.Choices[0].Text, nil
}

func (p *SiliconFlowProvider) GetName() string {
	return "硅基流动"
}

func (p *SiliconFlowProvider) GetModel() string {
	return p.config.Model
}

func (p *SiliconFlowProvider) SupportsFeature(feature string) bool {
	switch feature {
	case "quantization":
		return strings.HasPrefix(p.config.Model, "yi-")
	default:
		return false
	}
}
