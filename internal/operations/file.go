package operations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	
	"voidClint/pkg/types"
)

// FileManager 处理所有文件操作
type FileManager struct{}

// ReadFile 读取文件内容
func (fm *FileManager) ReadFile(path string) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	
	content := string(data)
	checksum := fm.calculateChecksum(content)
	return content, checksum, nil
}

// WriteFile 写入文件，支持多种模式
func (fm *FileManager) WriteFile(op types.FileOperation) error {
	// 创建备份
	backupPath := op.Path + ".bak"
	if _, err := os.Stat(op.Path); err == nil {
		if err := os.Rename(op.Path, backupPath); err != nil {
			return fmt.Errorf("无法创建备份: %v", err)
		}
	}
	
	// 根据模式处理
	switch op.Mode {
	case "replace":
		return ioutil.WriteFile(op.Path, []byte(op.Content), 0644)
		
	case "insert":
		current, err := ioutil.ReadFile(backupPath)
		if err != nil {
			return err
		}
		newContent := string(current[:op.Offset]) + op.Content + string(current[op.Offset:])
		return ioutil.WriteFile(op.Path, []byte(newContent), 0644)
		
	case "append":
		f, err := os.OpenFile(op.Path, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(op.Content)
		return err
		
	default:
		return fmt.Errorf("不支持的写入模式: %s", op.Mode)
	}
}

// CreateFile 创建新文件
func (fm *FileManager) CreateFile(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}

func (fm *FileManager) calculateChecksum(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// PreviewContent 生成内容预览
func (fm *FileManager) PreviewContent(content string) string {
	const maxPreview = 300
	if len(content) <= maxPreview {
		return content
	}
	return content[:150] + "\n... [截断] ...\n" + content[len(content)-150:]
}

// ResolvePath 安全解析路径
func (fm *FileManager) ResolvePath(basePath, targetPath string) (string, error) {
	fullPath := filepath.Join(basePath, targetPath)
	
	// 检查路径是否超出项目范围
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", err
	}
	
	absFull, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	
	if !strings.HasPrefix(absFull, absBase) {
		return "", fmt.Errorf("尝试访问项目外部路径: %s", targetPath)
	}
	
	return absFull, nil
}
