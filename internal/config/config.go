package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yantianyv/AkashaTerminal/pkg/types"
)

const (
	defaultConfigPath = ".config/akashaterminal/profiles.json"
)

// ConfigManager 管理所有 API 配置
type ConfigManager struct {
	Path         string
	Default      string
	Profiles     map[string]types.APIConfig
}

func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, defaultConfigPath)
	
	return &ConfigManager{
		Path: configPath,
	}
}

// Load 加载配置文件
func (cm *ConfigManager) Load() error {
	// 确保配置目录存在
	if err := os.MkdirAll(filepath.Dir(cm.Path), 0700); err != nil {
		return err
	}
	
	// 如果文件不存在，创建空配置
	if _, err := os.Stat(cm.Path); os.IsNotExist(err) {
		cm.Profiles = make(map[string]types.APIConfig)
		return cm.Save()
	}
	
	data, err := os.ReadFile(cm.Path)
	if err != nil {
		return err
	}
	
	var configData struct {
		DefaultProfile string                 `json:"default_profile"`
		Profiles       map[string]types.APIConfig `json:"profiles"`
	}
	
	if err := json.Unmarshal(data, &configData); err != nil {
		return err
	}
	
	cm.Default = configData.DefaultProfile
	cm.Profiles = configData.Profiles
	return nil
}

// Save 保存配置文件
func (cm *ConfigManager) Save() error {
	configData := struct {
		DefaultProfile string                 `json:"default_profile"`
		Profiles       map[string]types.APIConfig `json:"profiles"`
	}{
		DefaultProfile: cm.Default,
		Profiles:       cm.Profiles,
	}
	
	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(cm.Path, data, 0600)
}

// AddProfile 添加新配置
func (cm *ConfigManager) AddProfile(name string, config types.APIConfig) {
	cm.Profiles[name] = config
	cm.Save()
}

// DeleteProfile 删除配置
func (cm *ConfigManager) DeleteProfile(name string) {
	delete(cm.Profiles, name)
	cm.Save()
}

// SetDefault 设置默认配置
func (cm *ConfigManager) SetDefault(name string) error {
	if _, exists := cm.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}
	cm.Default = name
	return cm.Save()
}

// GetProfile 获取指定配置
func (cm *ConfigManager) GetProfile(name string) (types.APIConfig, error) {
	if config, exists := cm.Profiles[name]; exists {
		return config, nil
	}
	return types.APIConfig{}, fmt.Errorf("profile %s not found", name)
}
