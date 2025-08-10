package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	
	"AkashaTerminal/internal/config"
	"AkashaTerminal/internal/operations"
	"AkashaTerminal/internal/providers"
	"AkashaTerminal/internal/state"
	"AkashaTerminal/internal/utils"
	"AkashaTerminal/pkg/types"
)

var (
	profileFlag = flag.String("profile", "", "使用指定的API配置")
	maxDepth    = flag.Int("depth", 3, "目录扫描最大深度")
	maxTokens   = flag.Int("tokens", 8192, "最大上下文Token数")
)

func main() {
	flag.Parse()
	
	// 初始化配置管理器
	cfgMgr := config.NewConfigManager()
	if err := cfgMgr.Load(); err != nil {
		utils.ShowError("加载配置失败", err)
		os.Exit(1)
	}
	
	// 选择API配置
	profileName := *profileFlag
	if profileName == "" && cfgMgr.Default != "" {
		profileName = cfgMgr.Default
	}
	
	if profileName == "" && len(cfgMgr.Profiles) == 0 {
		utils.ShowWarning("未找到API配置，请先添加配置")
		os.Exit(1)
	}
	
	apiConfig, err := cfgMgr.GetProfile(profileName)
	if err != nil {
		utils.ShowError("获取配置失败", err)
		os.Exit(1)
	}
	
	// 创建供应商实例
	provider, err := providers.CreateProvider(apiConfig)
	if err != nil {
		utils.ShowError("创建AI提供程序失败", err)
		os.Exit(1)
	}
	
	// 初始化状态管理
	stateMgr := state.NewProjectState(*maxDepth, *maxTokens)
	if err := stateMgr.ScanInitialDirectory("."); err != nil {
		utils.ShowError("扫描目录失败", err)
		os.Exit(1)
	}
	
	// 初始化文件管理器
	fileMgr := operations.FileManager{}
	
	// 初始化Token管理
	tokenMgr := state.NewTokenManager(*maxTokens)
	
	fmt.Printf("\n✨ AkashaTerminal v1.0 - 智能代码助手\n")
	fmt.Printf("供应商: %s (%s)\n", provider.GetName(), provider.GetModel())
	fmt.Println("输入 '/exit' 退出, '/help' 查看帮助")
	fmt.Println(strings.Repeat("=", 50))
	
	// 主交互循环
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("\n> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		
		switch userInput {
		case "":
			continue
		case "/exit":
			fmt.Println("再见！")
			return
		case "/help":
			printHelp()
			continue
		case "/reload":
			stateMgr.ScanInitialDirectory(".")
			utils.ShowSuccess("目录状态已刷新")
			continue
		}
		
		// 构建完整提示
		prompt := buildFullPrompt(userInput, stateMgr)
		
		// 发送请求
		response, err := provider.SendRequest(ctx, prompt)
		if err != nil {
			utils.ShowError("AI请求失败", err)
			continue
		}
		
		// 解析操作指令
		operations, err := parseOperations(response)
		if err != nil {
			utils.ShowError("解析操作指令失败", err)
			continue
		}
		
		// 处理操作指令
		for _, op := range operations {
			if err := processOperation(ctx, op, fileMgr, stateMgr, tokenMgr); err != nil {
				utils.ShowError("操作执行失败", err)
			}
		}
		
		// 更新Token状态
		tokenMgr.AddRecord(userInput, operations)
		utils.DisplayTokenUsage(tokenMgr.GetTokenUsage())
	}
}

func buildFullPrompt(userInput string, stateMgr *state.ProjectState) string {
	stateJSON, _ := json.Marshal(stateMgr.GetCurrentState())
	
	// 添加系统消息
	prompt := fmt.Sprintf(`
系统信息:
- 当前目录: %s
- 扫描深度: %d层
- 文件总数: %d

用户请求:
%s

项目结构:
%s

操作指令:`,
		stateMgr.GetCWD(),
		stateMgr.GetMaxDepth(),
		len(stateMgr.GetFileStates()),
		userInput,
		stateMgr.GetDirectoryTree(),
	)
	
	// 添加当前文件状态
	for path, fileState := range stateMgr.GetFileStates() {
		if fileState.Content != "" {
			prompt += fmt.Sprintf("\n\n文件[%s]:\n```\n%s\n```", path, fileState.Content)
		} else {
			prompt += fmt.Sprintf("\n\n文件[%s]: (未加载内容)", path)
		}
	}
	
	return prompt
}

func parseOperations(response string) ([]types.FileOperation, error) {
	// 在实际实现中，这里需要解析AI返回的JSON结构
	// 以下是简化实现
	
	var operations []types.FileOperation
	
	// 假设AI返回的操作是JSON格式
	err := json.Unmarshal([]byte(response), &operations)
	if err != nil {
		return nil, fmt.Errorf("解析操作失败: %v", err)
	}
	
	return operations, nil
}

func processOperation(ctx context.Context, op types.FileOperation, fm operations.FileManager, 
	stateMgr *state.ProjectState, tokenMgr *state.TokenManager) error {
	
	switch op.Action {
	case "read":
		return handleReadOperation(op, fm, stateMgr)
		
	case "write", "create":
		if !utils.UserApproval(op, &fm) {
			return nil // 用户取消操作
		}
		return handleWriteOperation(op, fm, stateMgr)
		
	case "scan":
		return handleScanOperation(op, stateMgr)
		
	default:
		return fmt.Errorf("不支持的操作类型: %s", op.Action)
	}
}

func handleReadOperation(op types.FileOperation, fm operations.FileManager, stateMgr *state.ProjectState) error {
	// 解析安全路径
	path, err := fm.ResolvePath(stateMgr.GetCWD(), op.Path)
	if err != nil {
		return err
	}
	
	// 读取文件
	content, checksum, err := fm.ReadFile(path)
	if err != nil {
		return err
	}
	
	// 更新状态
	stateMgr.UpdateFileState(op.Path, content, checksum)
	
	utils.ShowSuccess(fmt.Sprintf("已读取文件: %s (%d字节)", op.Path, len(content)))
	return nil
}

func handleWriteOperation(op types.FileOperation, fm operations.FileManager, stateMgr *state.ProjectState) error {
	// 解析安全路径
	path, err := fm.ResolvePath(stateMgr.GetCWD(), op.Path)
	if err != nil {
		return err
	}
	
	// 执行操作
	switch op.Action {
	case "write":
		if op.Mode == "" {
			op.Mode = "replace" // 默认模式
		}
		if err := fm.WriteFile(op); err != nil {
			return err
		}
		utils.ShowSuccess(fmt.Sprintf("已更新文件: %s", op.Path))
		
	case "create":
		if err := fm.CreateFile(path, op.Content); err != nil {
			return err
		}
		utils.ShowSuccess(fmt.Sprintf("已创建文件: %s", op.Path))
	}
	
	// 更新状态
	content, checksum, _ := fm.ReadFile(path)
	stateMgr.UpdateFileState(op.Path, content, checksum)
	
	return nil
}

func handleScanOperation(op types.FileOperation, stateMgr *state.ProjectState) error {
	utils.ShowWarning(fmt.Sprintf("AI请求扫描目录: %s", op.Path))
	if !utils.GetUserConfirmation("确认扫描? (y/n) > ") {
		return nil
	}
	
	// 执行扫描
	if err := stateMgr.ScanAdditionalDirectory(op.Path); err != nil {
		return err
	}
	
	utils.ShowSuccess(fmt.Sprintf("已扫描目录: %s (%d个文件)", op.Path, len(stateMgr.GetScannedFiles(op.Path))))
	return nil
}

func printHelp() {
	fmt.Println("\n可用命令:")
	fmt.Println("  /exit       - 退出程序")
	fmt.Println("  /help       - 显示此帮助信息")
	fmt.Println("  /reload     - 重新扫描当前目录")
	fmt.Println()
	fmt.Println("操作支持:")
	fmt.Println("  AI可执行读取(read)、写入(write)、创建(create)文件和扫描(scan)目录操作")
	fmt.Println("  写入操作支持以下模式: replace(替换), insert(插入), append(追加)")
}
