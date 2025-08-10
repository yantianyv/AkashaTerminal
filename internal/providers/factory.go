package providers

import (
	"fmt"
	
	"voidClint/pkg/types"
	"voidClint/internal/providers/openai"
	"voidClint/internal/providers/azure"
	"voidClint/internal/providers/deepseek"
	"voidClint/internal/providers/bailian"
	"voidClint/internal/providers/siliconflow"
	"voidClint/internal/providers/custom"
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
