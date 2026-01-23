//go:build !windows && !linux

// 非 Windows 平台不提供 GUI。
package main

import "fmt"

func main() {
	fmt.Println("当前平台不支持 GUI，请使用 CLI 启动。")
}
