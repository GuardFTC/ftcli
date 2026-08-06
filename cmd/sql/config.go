// Package sql @Author:冯铁城 [17615007230@163.com] 2025-10-31 19:47:32
package sql

import (
	"runtime"
	"time"
)

// 定义系统常量
const windows = "windows"
const mac = "darwin"

// 系统名称
var system = runtime.GOOS

// 默认CSV文件路径
var defaultPath = map[string]string{
	windows: "C:\\Users\\Administrator\\Downloads\\",
	mac:     "/Users/m/Downloads/",
}

// Stream Load 服务地址
const streamLoadURL = "http://localhost:6677/api/csv/import"

// HTTP 请求超时时间（海量数据导入耗时较长，设为10分钟）
const httpTimeout = 10 * time.Minute
