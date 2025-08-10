package app

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
	"github.com/yantianyv/AkashaTerminal/internal/commands"
)

func Run() error {
	rootCmd := &cobra.Command{
		Use:   "akasha",
		Short: "github.com/yantianyv/AkashaTerminal - 智能代码助手",
		Long: `github.com/yantianyv/AkashaTerminal 是一个基于 AI 的代码助手工具，
支持多种 AI 供应商，提供智能代码生成、分析和重构功能。`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("欢迎使用 github.com/yantianyv/AkashaTerminal！输入 'akasha help' 查看可用命令")
		},
	}

	// 添加子命令
	rootCmd.AddCommand(commands.NewRunCommand())
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewVersionCommand())

	return rootCmd.Execute()
}

