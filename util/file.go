// util/file.go
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// WalkDir 递归遍历目录，将符合后缀的文件路径发送到 fileCh
func WalkDir(rootDir, ext string, fileCh chan<- string) error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("访问路径 %s 失败: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if ext == "" || strings.HasSuffix(path, ext) {
			fileCh <- path
		}
		return nil
	})
}

// Stats 用于记录处理过程中的统计信息
type Stats struct {
	sync.Mutex
	Total    int
	Success  int
	Fail     int
	Matches  int // 用于 find
	Replaces int // 用于 replace
	Converts int // 用于 convert
	Errors   map[string]error
}

// NewStats 创建一个新的 Stats 实例
func NewStats() *Stats {
	return &Stats{
		Errors: make(map[string]error),
	}
}

// AddSuccess 增加成功计数
func (s *Stats) AddSuccess() {
	s.Lock()
	defer s.Unlock()
	s.Total++
	s.Success++
}

// AddError 增加失败计数并记录错误
func (s *Stats) AddError(filePath string, err error) {
	s.Lock()
	defer s.Unlock()
	s.Total++
	s.Fail++
	s.Errors[filePath] = err
}

// AddMatches 增加匹配计数
func (s *Stats) AddMatches(count int) {
	s.Lock()
	defer s.Unlock()
	s.Matches += count
}

// AddReplaces 增加替换计数
func (s *Stats) AddReplaces(count int) {
	s.Lock()
	defer s.Unlock()
	s.Replaces += count
}

// AddConverts 增加转换计数
func (s *Stats) AddConverts() {
	s.Lock()
	defer s.Unlock()
	s.Converts++
}
