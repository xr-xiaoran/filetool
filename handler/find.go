// handler/find.go
package handler

import (
	"bufio"
	"filetool/util"
	"fmt"
	"os"
	"strings"
)

// FindHandler 处理文件查找功能
func FindHandler(filePath, target string, stats *util.Stats) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fmt.Printf("正在查找文件: %s\n", filePath)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	matchCount := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(line, target) {
			matchCount++
			fmt.Printf("  匹配行 %d: %s\n", lineNum, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	stats.AddSuccess()
	stats.AddMatches(matchCount)
	return nil
}
