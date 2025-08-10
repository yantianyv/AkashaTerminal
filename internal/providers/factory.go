package providers

import (
	"fmt"
	
	"AkashaTerminal/pkg/types"
	"AkashaTerminal/internal/providers/openai"
	"AkashaTerminal/internal/providers/azure"
	"AkashaTerminal/internal/providers/deepseek"
	"AkashaTerminal/internal/providers/bailian"
	"AkashaTerminal/internal/providers/siliconflow"
	"AkashaTerminal/internal/providers/custom"
)

// CreateProvider 根据配置创建供应商实例
func CreateProvider(config types.APIConfig) (types.AIProvider, error) {
	switch config.Provider {
	case "openai":
		return openai.NewProvider(config)
	case "azure":
		return azure.NewProvider(config)
	case "deepseek":
		return deepseek.NewProvider(config)
	case "bailian":
		return bailian.NewProvider(config)
	case "siliconflow":
		return siliconflow.NewProvider(config)
	case "custom":
		return custom.NewProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}
