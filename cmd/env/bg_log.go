// Package env @Author:冯铁城 [17615007230@163.com] 2026-06-08 16:00:00
package env

import (
	"bufio"
	"fmt"
	"ftcli/util"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// runBgLog 滚动查看后台服务日志
func runBgLog(serviceName string) {

	//1.如果未指定服务名，列出可选服务
	if serviceName == "" || serviceName == "list" {
		printAvailableServices()
		return
	}

	//2.查找服务配置
	serviceType, propertyValues := findServiceConfig(serviceName)

	//3.根据类型查看日志
	switch serviceType {
	case "background":
		tailBackgroundLog(propertyValues)
	case "docker":
		tailDockerServiceLog(serviceName, propertyValues)
	case "docker-container":
		util.TailDockerLog(serviceName)
	default:
		fmt.Printf("未找到服务[%s]的日志配置，可选服务:\n", serviceName)
		printAvailableServices()
	}
}

// tailBackgroundLog 滚动查看后台进程的日志文件
func tailBackgroundLog(propertyValues []string) {

	//1.获取日志文件路径
	logFile := propertyValues[1]

	//2.滚动输出日志
	tailLog(logFile)
}

// tailDockerServiceLog 查看docker服务日志（如果有多个容器，查看第一个）
func tailDockerServiceLog(serviceName string, propertyValues []string) {

	//1.解析容器列表
	containers := propertyValues[2:]

	//2.如果只有一个容器，直接查看
	if len(containers) == 1 {
		util.TailDockerLog(containers[0])
		return
	}

	//3.多个容器时，提示用户选择
	fmt.Printf("服务[%s]包含多个容器，请指定容器名:\n", serviceName)
	for _, c := range containers {
		if c == serviceName {
			fmt.Printf("  * %s (与服务名同名，请用: ftcli env --blc %s)\n", c, c)
		} else {
			fmt.Printf("  * %s\n", c)
		}
	}
	fmt.Println()
	fmt.Println("用法: ftcli env --bl <容器名>")
}

// runBgLogContainer 强制按容器名查看docker容器日志
// 用于服务名与容器名冲突的场景（如kafka服务包含kafka容器）：
// --bl会优先匹配服务名导致走到列容器分支，此时用--blc绕过服务名匹配直接tail容器
func runBgLogContainer(containerName string) {
	util.TailDockerLog(containerName)
}

// printAvailableServices 打印可选的服务名
func printAvailableServices() {

	//1.获取当前系统对应的项目集合
	systemProjects := envCmdProjectPropertiesMap[system]
	if systemProjects == nil {
		fmt.Println("当前系统无服务配置")
		return
	}

	//2.打印提示
	fmt.Println("可选的服务:")

	//3.遍历所有项目
	for projectName, projectProperties := range systemProjects {

		//4.遍历项目配置
		for serviceName, propertyValues := range projectProperties {

			//5.根据类型打印
			switch propertyValues[0] {
			case "background":
				fmt.Printf("  * %s [background] (项目: %s)\n", serviceName, projectName)
			case "docker":
				containers := propertyValues[2:]
				fmt.Printf("  * %s [docker: %s] (项目: %s)\n", serviceName, strings.Join(containers, ","), projectName)
			}
		}
	}

	//6.打印用法提示
	fmt.Println()
	fmt.Println("用法: ftcli env --bl <服务名或容器名>")
}

// tailLog 滚动输出日志文件（类似 tail -f，先输出最后100行）
func tailLog(logFile string) {

	//1.打开日志文件
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("打开日志文件失败: %v\n", err)
		return
	}
	defer file.Close()

	//2.打印提示
	fmt.Printf(">>> 正在滚动查看日志: %s (按 Ctrl+C 退出)\n", logFile)
	fmt.Println("--------------------------------------------------------------------------------")

	//3.先输出最后100行
	printLastNLines(logFile, 100)

	//4.移动到文件末尾，开始追踪新内容
	file.Seek(0, io.SeekEnd)

	//5.监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//6.创建读取器
	reader := bufio.NewReader(file)

	//7.循环读取新内容
	for {

		//8.检查退出信号
		select {
		case <-sigChan:
			fmt.Println("\n>>> 已退出日志查看")
			return
		default:
		}

		//9.尝试读取一行
		line, err := reader.ReadString('\n')
		if err != nil {

			//10.如果是EOF，等待100ms后继续
			time.Sleep(100 * time.Millisecond)
			continue
		}

		//11.输出日志
		fmt.Print(line)
	}
}

// printLastNLines 打印文件最后N行
func printLastNLines(filePath string, n int) {

	//1.读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	//2.按行分割
	lines := strings.Split(string(data), "\n")

	//3.计算起始行
	start := len(lines) - n
	if start < 0 {
		start = 0
	}

	//4.输出
	for _, line := range lines[start:] {
		if line != "" {
			fmt.Println(line)
		}
	}
}
