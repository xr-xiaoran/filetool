// main.go
package main

import (
	"filetool/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("执行命令失败: %v", err)
	}
}
