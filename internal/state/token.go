package state

import (
	"encoding/json"
	"fmt"
	"strings"
	
	"voidClint/pkg/types"
)

// TokenEstimator 估算文本的 token 数量
type TokenEstimator struct{}

func (te *TokenEstimator) Estimate(text string) int {
	// 简易估算：4个英文字符 ≈ 1 token，1个汉字 ≈ 2 tokens
	total := 0
	for _, r := range text {
		if r <= 255 {
			total += 1
		} else {
			total += 2
		}
	}
	return total / 4
}

// TokenManager 管理 token 使用情况
type TokenManager struct {
	maxTokens      int
	currentToken   int
	history        []*ConversationRecord
}

type ConversationRecord struct {
	ID         int               `json:"id"`
	Role       string            `json:"role"`
	Content    string            `json:"content"`
	Operation  types.FileOperation `json:"operation"`
	TokenCount int               `json:"token_count"`
}

func NewTokenManager(maxTokens int) *TokenManager {
	return &TokenManager{
		maxTokens:    maxTokens,
		currentToken: 0,
	}
}

func (tm *TokenManager) AddRecord(record *ConversationRecord) error {
	estimator := TokenEstimator{}
	
	// 估算 token 使用
	tokens := estimator.Estimate(record.Content)
	if record.Operation.Content != "" {
		tokens += estimator.Estimate(record.Operation.Content)
	}
	
	record.TokenCount = tokens
	tm.currentToken += tokens
	
	// 应用清理策略
	return tm.applyCleanupStrategy()
}

func (tm *TokenManager) applyCleanupStrategy() error {
	threshold50 := tm.maxTokens / 2
	threshold75 := tm.maxTokens * 3 / 4
	threshold100 := tm.maxTokens
	
	switch {
	case tm.currentToken > threshold100:
		return fmt.Errorf("上下文Token超限 (%d/%d)", tm.currentToken, tm.maxTokens)
		
	case tm.currentToken > threshold75:
		tm.level2Cleanup()
		
	case tm.currentToken > threshold50:
		tm.level1Cleanup()
	}
	
	return tm.checkCleanupEffectiveness()
}

func (tm *TokenManager) level1Cleanup() {
	// 保留最近2条记录
	keep := min(2, len(tm.history))
	preserved := tm.history[len(tm.history)-keep:]
	
	// 清理3轮前的记录
	for i := 0; i < len(tm.history)-keep && tm.currentToken > tm.maxTokens/2; i++ {
		rec := tm.history[i]
		
		// 清理写入操作详情
		if rec.Operation.Action == "write" && rec.Operation.Content != "" {
			tokenSave := tm.tokenEstimator.Estimate(rec.Operation.Content)
			rec.Operation.Content = ""
			rec.Operation.OldText = ""
			rec.Operation = fmt.Sprintf("已清理的写入操作: %s", rec.Operation.Path)
			rec.TokenCount -= tokenSave
			tm.currentToken -= tokenSave
		}
		
		// 标记为可清理
		if i > 1 {
			tm.history[i] = nil
		}
	}
	
	// 压缩历史记录
	newHistory := make([]*ConversationRecord, 0, len(tm.history))
	for _, rec := range tm.history {
		if rec != nil {
			newHistory = append(newHistory, rec)
		}
	}
	tm.history = newHistory
}

func (tm *TokenManager) level2Cleanup() {
	// 先执行一级清理
	tm.level1Cleanup()
	
	// 清理非关键文件读取内容
	for i, rec := range tm.history {
		if i > 0 && rec.Operation.Action == "read" && tm.currentToken > tm.maxTokens/2 {
			tokenSave := tm.tokenEstimator.Estimate(rec.Operation.Content)
			rec.Operation.Content = ""
			rec.Operation = fmt.Sprintf("已清理的文件读取记录: %s", rec.Operation.Path)
			rec.TokenCount -= tokenSave
			tm.currentToken -= tokenSave
		}
	}
	
	// 简化历史消息
	for i, rec := range tm.history {
		if i > 0 && rec != nil && tm.currentToken > tm.maxTokens*3/5 {
			if !strings.HasPrefix(rec.Content, "[关键]") {
				tokenSave := rec.TokenCount - 10 // 保留基本信息
				rec.Content = fmt.Sprintf("[精简] 记录 #%d (%s)", rec.ID, rec.Operation.Action)
				rec.Operation = types.FileOperation{}
				rec.TokenCount = 10
				tm.currentToken -= tokenSave
			}
		}
	}
}

func (tm *TokenManager) checkCleanupEffectiveness() error {
	if tm.currentToken > tm.maxTokens {
		return fmt.Errorf("清理后仍超过100%% Token限制 (%d/%d)", tm.currentToken, tm.maxTokens)
	}
	
	if tm.currentToken > tm.maxTokens*3/4 {
		return fmt.Errorf("清理后仍超过75%% Token限制 (%d/%d)", tm.currentToken, tm.maxTokens)
	}
	
	if tm.currentToken > tm.maxTokens/2 {
		return fmt.Errorf("清理后仍超过50%% Token限制 (%d/%d)", tm.currentToken, tm.maxTokens)
	}
	
	return nil
}

func (tm *TokenManager) GetTokenUsage() (int, int) {
	return tm.currentToken, tm.maxTokens
}
