package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	
	"voidClint/pkg/types"
)

type BailianProvider struct {
	config types.APIConfig
}

func NewProvider(config types.APIConfig) (*BailianProvider, error) {
	if config.AppID == "" || config.AccessKeyID == "" || config.AccessKeySecret == "" {
		return nil, fmt.Errorf("Bailian requires AppID, AccessKeyID and AccessKeySecret")
	}
	
	return &BailianProvider{
		config: config,
	}, nil
}

func (p *BailianProvider) generateSignature() (string, int64) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	strToSign := fmt.Sprintf("%d", timestamp)
	
	h := hmac.New(sha1.New, []byte(p.config.AccessKeySecret))
	h.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	
	return signature, timestamp
}

func (p *BailianProvider) SendRequest(prompt string, state interface{}) (string, error) {
	url := "https://bailian.aliyuncs.com/v2/app/completions"
	if p.config.APIBase != "" {
		url = p.config.APIBase
	}
	
	signature, timestamp := p.generateSignature()
	
	requestBody := map[string]interface{}{
		"appId":    p.config.AppID,
		"prompt":   prompt,
		"sessionId": time.Now().Unix(),
	}
	
	if p.config.AgentID != "" {
		requestBody["agentId"] = p.config.AgentID
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
	req.Header.Set("X-Bailian-AppId", p.config.AppID)
	req.Header.Set("X-Bailian-Token", signature)
	req.Header.Set("X-Bailian-Timestamp", fmt.Sprintf("%d", timestamp))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Bailian API error: %s - %s", resp.Status, string(body))
	}
	
	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Text string `json:"text"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	
	if !response.Success {
		return "", fmt.Errorf("Bailian API returned non-success response")
	}
	
	return response.Data.Text, nil
}

func (p *BailianProvider) GetName() string {
	return "阿里百炼"
}

func (p *BailianProvider) GetModel() string {
	if p.config.Model != "" {
		return p.config.Model
	}
	return "bailian-plus"
}

func (p *BailianProvider) SupportsFeature(feature string) bool {
	switch feature {
	case "enterprise":
		return true
	case "custom_model":
		return p.config.Model != "bailian-plus"
	default:
		return false
	}
}
