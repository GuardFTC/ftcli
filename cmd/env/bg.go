// Package env @Author:冯铁城 [17615007230@163.com] 2026-06-08 16:00:00
package env

import (
	"fmt"
	"ftcli/util"
	"os"
	"strings"
	"text/tabwriter"
)

// runListBgServices 查看所有后台运行进程状态
func runListBgServices() {

	//1.获取当前系统对应的项目集合
	systemProjects := envCmdProjectPropertiesMap[system]
	if systemProjects == nil {
		fmt.Println("当前系统无服务配置")
		return
	}

	//2.打印分割线
	fmt.Println("--------------------------------------------------------------------------------")

	//3.打印表头
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "| 项目\t| 服务名\t| 类型\t| 容器/端口\t| 状态\t| 日志路径\t")
	fmt.Fprintln(w, "--------------------------------------------------------------------------------")

	//4.遍历所有项目
	for projectName, projectProperties := range systemProjects {

		//5.遍历项目配置
		for serviceName, propertyValues := range projectProperties {

			//6.根据类型分别处理
			switch propertyValues[0] {
			case "background":
				printBackgroundService(w, projectName, serviceName, propertyValues)
			case "docker":
				printDockerService(w, projectName, serviceName, propertyValues)
			}
		}
	}

	//7.写入控制台
	fmt.Fprintln(w, "--------------------------------------------------------------------------------")
	w.Flush()
}

// printBackgroundService 打印后台进程服务状态
func printBackgroundService(w *tabwriter.Writer, projectName string, serviceName string, propertyValues []string) {

	//1.解析配置
	logFile := propertyValues[1]
	checkPort := propertyValues[2]
	killName := propertyValues[3]
	killKeyword := propertyValues[4]

	//2.检查进程是否存活
	status := "未运行"
	if util.IsProcessRunning(killName, killKeyword) {
		status = "运行中"
	}

	//3.打印
	fmt.Fprintf(w, "| %s\t| %s\t| background\t| %s\t| %s\t| %s\t\n",
		projectName, serviceName, checkPort, status, logFile)
}

// printDockerService 打印docker容器服务状态
func printDockerService(w *tabwriter.Writer, projectName string, serviceName string, propertyValues []string) {

	//1.解析配置: "docker", compose文件路径, 容器名1, 容器名2...
	containers := propertyValues[2:]

	//2.遍历容器检查状态
	for _, containerName := range containers {

		//3.检查容器是否运行
		status := "未运行"
		if util.IsDockerContainerRunning(containerName) {
			status = "运行中"
		}

		//4.打印
		fmt.Fprintf(w, "| %s\t| %s(%s)\t| docker\t| %s\t| %s\t| -\t\n",
			projectName, serviceName, containerName, strings.Join(containers, ","), status)
	}
}
