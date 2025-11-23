// handler/replace.go
package handler

import (
	"filetool/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceHandler 处理文件内容替换功能
func ReplaceHandler(filePath, target, replaceWith string, stats *util.Stats) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	originalStr := string(content)
	newStr := strings.ReplaceAll(originalStr, target, replaceWith)
	replaceCount := strings.Count(originalStr, target)

	if replaceCount == 0 {
		fmt.Printf("文件 %s 中未找到目标字符串 %s，跳过。\n", filePath, target)
		stats.AddSuccess() // 虽然没替换，但也算处理成功
		return nil
	}

	// 构建新文件名
	dir := filepath.Dir(filePath)
	filename := filepath.Base(filePath)
	newFilePath := filepath.Join(dir, fmt.Sprintf("%s_new%s", strings.TrimSuffix(filename, filepath.Ext(filename)), filepath.Ext(filename)))

	err = os.WriteFile(newFilePath, []byte(newStr), 0644)
	if err != nil {
		return fmt.Errorf("写入新文件失败: %w", err)
	}

	fmt.Printf("处理完成: %s -> %s, 替换次数: %d\n", filePath, newFilePath, replaceCount)

	stats.AddSuccess()
	stats.AddReplaces(replaceCount)
	return nil
}
