package commands

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "启动 github.com/yantianyv/AkashaTerminal 交互式会话",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("github.com/yantianyv/AkashaTerminal 会话启动...")
			// 这里将添加主交互逻辑
		},
	}
	
	// 添加运行参数
	cmd.Flags().StringP("profile", "p", "", "指定使用的 API 配置")
	cmd.Flags().IntP("tokens", "t", 8192, "最大上下文 Token 数")
	cmd.Flags().IntP("depth", "d", 3, "目录扫描最大深度")
	
	return cmd
}

