// Package ai @Author:冯铁城 [17615007230@163.com] 2026-07-28 15:00:00
package ai

import (
	"fmt"
	"ftcli/util"
)

// runMcpWeb 打开MCP管理页面
func runMcpWeb() {

	//1.定义URL
	url := baseURL + "/mcp.html"

	//2.打开浏览器
	if err := util.OpenBrowser(url); err != nil {
		fmt.Printf("打开浏览器失败: %v\n", err)
		return
	}

	//3.日志打印
	fmt.Printf("已打开页面: %s\n", url)
}
