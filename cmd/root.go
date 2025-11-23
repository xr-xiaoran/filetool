// cmd/root.go
package cmd

import (
	"filetool/handler"
	"filetool/util"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	dir         string
	action      string
	target      string
	replace     string
	ext         string
	concurrency int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "filetool",
	Short: "一个实用的本地文件批量处理工具",
	Long:  `filetool 是一个用 Go 语言实现的命令行工具，支持文件查找、内容替换和格式转换（TXT->JSON）等功能。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 校验参数
		if err := validateArgs(); err != nil {
			fmt.Printf("参数错误: %v\n", err)
			cmd.Help()
			os.Exit(1)
		}

		// 2. 初始化统计信息
		stats := util.NewStats()

		// 3. 打印开始信息
		fmt.Printf("开始处理目录: %s，操作类型: %s，目标后缀: %s，最大并发数: %d\n", dir, action, ext, concurrency)

		// 4. 并发遍历目录，获取文件通道
		fileCh := make(chan string)
		var wgWalk sync.WaitGroup
		wgWalk.Add(1)
		go func() {
			defer wgWalk.Done()
			if err := util.WalkDir(dir, ext, fileCh); err != nil {
				fmt.Printf("目录遍历失败: %v\n", err)
			}
		}()

		// 关闭文件通道的 Goroutine
		go func() {
			wgWalk.Wait()
			close(fileCh)
		}()

		// 5. 根据 action 选择处理器
		var fileHandler func(string, *util.Stats) error
		switch action {
		case "find":
			fileHandler = func(f string, s *util.Stats) error {
				return handler.FindHandler(f, target, s)
			}
		case "replace":
			fileHandler = func(f string, s *util.Stats) error {
				return handler.ReplaceHandler(f, target, replace, s)
			}
		case "convert":
			fileHandler = func(f string, s *util.Stats) error {
				return handler.ConvertHandler(f, s)
			}
		}

		// 6. 并发处理文件
		var wgProcess sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wgProcess.Add(1)
			go func() {
				defer wgProcess.Done()
				for filePath := range fileCh {
					if err := fileHandler(filePath, stats); err != nil {
						stats.AddError(filePath, err)
						fmt.Printf("处理失败: %s, 错误: %v\n", filePath, err)
					}
				}
			}()
		}

		// 7. 等待所有处理完成并打印统计结果
		wgProcess.Wait()
		printStats(stats)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 这里定义标志和配置。
	rootCmd.Flags().StringVarP(&dir, "dir", "d", "", "指定要处理的根目录 (必填)")
	rootCmd.MarkFlagRequired("dir")

	rootCmd.Flags().StringVarP(&action, "action", "a", "", "指定操作类型 (必填, 支持 find/replace/convert)")
	rootCmd.MarkFlagRequired("action")

	rootCmd.Flags().StringVarP(&target, "target", "t", "", "查找/替换的目标字符串 (action=find 或 replace 时必填)")
	rootCmd.Flags().StringVarP(&replace, "replace", "r", "", "替换后的字符串 (action=replace 时必填)")
	rootCmd.Flags().StringVarP(&ext, "ext", "e", "", "指定文件后缀 (可选, 如 \".txt\")")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 5, "最大并发数 (可选, 默认 5)")
}

// validateArgs 进行更细致的参数校验
func validateArgs() error {
	if dir == "" {
		return fmt.Errorf("-dir 不能为空")
	}
	if action == "" {
		return fmt.Errorf("-action 不能为空")
	}
	if action != "find" && action != "replace" && action != "convert" {
		return fmt.Errorf("不支持的 action: %s", action)
	}
	if (action == "find" || action == "replace") && target == "" {
		return fmt.Errorf("action 为 %s 时, -target 不能为空", action)
	}
	if action == "replace" && replace == "" {
		return fmt.Errorf("action 为 replace 时, -replace 不能为空")
	}
	if concurrency <= 0 {
		return fmt.Errorf("concurrency 必须大于 0")
	}
	return nil
}

// printStats 打印最终的统计结果
func printStats(stats *util.Stats) {
	fmt.Println("\n处理统计:")
	fmt.Printf("总文件数: %d, 成功数: %d, 失败数: %d\n", stats.Total, stats.Success, stats.Fail)
	switch action {
	case "find":
		fmt.Printf("匹配到的总行数: %d\n", stats.Matches)
	case "replace":
		fmt.Printf("替换的总次数: %d\n", stats.Replaces)
	case "convert":
		fmt.Printf("转换成功的 JSON 文件数: %d\n", stats.Converts)
	}
}
