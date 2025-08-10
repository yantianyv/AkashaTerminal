package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/fatih/color"
	"github.com/yantianyv/AkashaTerminal/pkg/types"
	"github.com/yantianyv/AkashaTerminal/internal/operations"
)

// UserPrompt 显示用户提示并获取输入
func UserPrompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// UserApproval 获取用户对操作的批准
func UserApproval(op types.FileOperation, fm *operations.FileManager) bool {
	color.Yellow("\n⚠️ 确认操作: %s", strings.ToUpper(op.Action))
	color.Cyan("路径: %s", op.Path)
	
	if op.Content != "" {
		fmt.Println("\n内容预览:")
		fmt.Println(fm.PreviewContent(op.Content))
	}
	
	if op.Mode != "" {
		color.Magenta("模式: %s", op.Mode)
	}
	
	fmt.Print("\n确认执行? (y/n) > ")
	return GetUserConfirmation()
}

// GetUserConfirmation 获取简单的Y/N确认
func GetUserConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// DisplayTokenUsage 显示Token使用情况
func DisplayTokenUsage(current, max int) {
	percentage := float64(current) / float64(max) * 100
	status := "正常"
	colorStatus := color.GreenString
	
	switch {
	case percentage > 90:
		status = "严重"
		colorStatus = color.RedString
	case percentage > 75:
		status = "警告"
		colorStatus = color.YellowString
	case percentage > 50:
		status = "注意"
		colorStatus = color.CyanString
	}
	
	fmt.Printf("\nToken使用情况: %s (%.1f%%) %s/%s\n", 
		colorStatus(status), percentage, formatTokenCount(current), formatTokenCount(max))
}

func formatTokenCount(count int) string {
	if count > 1000 {
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

// ShowError 显示错误信息
func ShowError(message string, err error) {
	color.Red("\n❌ 错误: %s", message)
	if err != nil {
		color.Red("    -> %v", err)
	}
}

// ShowSuccess 显示成功信息
func ShowSuccess(message string) {
	color.Green("\n✅ %s", message)
}

// ShowWarning 显示警告信息
func ShowWarning(message string) {
	color.Yellow("\n⚠️ 注意: %s", message)
}
