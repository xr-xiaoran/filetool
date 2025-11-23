// handler/convert.go
package handler

import (
	"bufio"
	"encoding/json"
	"filetool/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConvertHandler 处理 TXT 到 JSON 的转换
func ConvertHandler(filePath string, stats *util.Stats) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	fmt.Printf("正在转换文件: %s\n", filePath)

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("第 %d 行格式错误, 应为 'key:value': %s", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		data[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON 序列化失败: %w", err)
	}

	newFilePath := filepath.ChangeExtension(filePath, ".json")
	err = os.WriteFile(newFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("写入 JSON 文件失败: %w", err)
	}

	fmt.Printf("转换完成: %s -> %s\n", filePath, newFilePath)

	stats.AddSuccess()
	stats.AddConverts()
	return nil
}
