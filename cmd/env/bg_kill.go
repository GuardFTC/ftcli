// Package env @Author:冯铁城 [17615007230@163.com] 2026-06-08 17:30:00
package env

import (
	"fmt"
	"slices"

	"ftcli/util"
)

// runBgKill 停止后台服务
func runBgKill(serviceName string) {

	//1.获取当前系统对应的项目集合
	systemProjects := envCmdProjectPropertiesMap[system]
	if systemProjects == nil {
		fmt.Println("当前系统无服务配置")
		return
	}

	//2.查找服务配置
	serviceType, propertyValues := findServiceConfig(serviceName)

	//3.根据类型执行停止操作
	switch serviceType {
	case "background":
		killBackgroundService(serviceName, propertyValues)
	case "docker":
		killDockerService(serviceName, propertyValues)
	case "docker-container":
		killDockerContainer(serviceName)
	default:
		fmt.Printf("未找到服务[%s]的配置，用 -b 查看可选服务\n", serviceName)
	}
}

// findServiceConfig 根据服务名或容器名查找服务配置，返回类型和配置值
func findServiceConfig(serviceName string) (string, []string) {

	//1.获取当前系统对应的项目集合
	systemProjects := envCmdProjectPropertiesMap[system]
	if systemProjects == nil {
		return "", nil
	}

	//2.遍历所有项目查找匹配的服务
	for _, projectProperties := range systemProjects {
		for name, propertyValues := range projectProperties {

			//3.服务名精确匹配
			if name == serviceName {
				return propertyValues[0], propertyValues
			}

			//4.如果是docker类型，检查容器名是否匹配
			if propertyValues[0] == "docker" {
				containers := propertyValues[2:]
				if slices.Contains(containers, serviceName) {
					return "docker-container", propertyValues
				}
			}
		}
	}

	//3.未找到返回空
	return "", nil
}

// killBackgroundService 停止后台进程服务
func killBackgroundService(serviceName string, propertyValues []string) {

	//1.解析kill配置
	killName := propertyValues[3]
	killKeyword := propertyValues[4]

	//2.执行kill
	fmt.Printf(">>> 停止服务: %s\n", serviceName)
	if err := util.KillProcess(killName, killKeyword); err != nil {
		fmt.Printf("停止失败: %v\n", err)
		return
	}

	//3.打印完成提示
	fmt.Printf(">>> 服务 %s 已停止\n", serviceName)
}

// killDockerService 停止docker服务组（按服务名停止整组）
func killDockerService(serviceName string, propertyValues []string) {

	//1.解析配置
	composeFile := propertyValues[1]
	containers := propertyValues[2:]

	//2.如果有compose文件，使用docker compose down停止整组
	if composeFile != "" {
		fmt.Printf(">>> 停止docker compose组: %s\n", serviceName)
		if err := util.StopDockerCompose(composeFile); err != nil {
			fmt.Printf("停止失败: %v\n", err)
			return
		}
		fmt.Printf(">>> docker compose组 %s 已停止\n", serviceName)
		return
	}

	//3.无compose文件，逐个停止容器
	fmt.Printf(">>> 停止docker服务: %s\n", serviceName)
	for _, containerName := range containers {
		fmt.Printf(">>> 停止容器: %s\n", containerName)
		if err := util.StopDockerContainer(containerName); err != nil {
			fmt.Printf("停止容器[%s]失败: %v\n", containerName, err)
			continue
		}
		fmt.Printf(">>> 容器 %s 已停止\n", containerName)
	}
}

// killDockerContainer 停止单个docker容器（按容器名停止）
func killDockerContainer(containerName string) {

	//1.执行停止
	fmt.Printf(">>> 停止容器: %s\n", containerName)
	if err := util.StopDockerContainer(containerName); err != nil {
		fmt.Printf("停止失败: %v\n", err)
		return
	}

	//2.打印完成提示
	fmt.Printf(">>> 容器 %s 已停止\n", containerName)
}
