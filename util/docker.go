// Package common @Author:冯铁城 [17615007230@163.com] 2026-07-20 10:00:00
package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// IsDockerContainerRunning 检查docker容器是否正在运行，传入容器名，返回是否运行
func IsDockerContainerRunning(containerName string) bool {

	//1.执行docker inspect命令检查容器状态
	cmd := exec.Command("docker", "inspect", "--format", "{{.State.Running}}", containerName)

	//2.获取命令输出
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	//3.判断输出是否为true
	return strings.TrimSpace(string(output)) == "true"
}

// StopDockerContainer 停止单个docker容器，传入容器名，返回错误
func StopDockerContainer(containerName string) error {

	//1.执行docker stop命令
	cmd := exec.Command("docker", "stop", containerName)

	//2.重定向输出到控制台
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//3.运行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止容器[%s]失败: %v", containerName, err)
	}

	//4.默认返回
	return nil
}

// StopDockerCompose 停止整个docker compose组，传入compose文件路径，返回错误
func StopDockerCompose(composeFile string) error {

	//1.执行docker compose down命令
	cmd := exec.Command("docker", "compose", "-f", composeFile, "down")

	//2.重定向输出到控制台
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//3.运行命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止compose组[%s]失败: %v", composeFile, err)
	}

	//4.默认返回
	return nil
}

// TailDockerLog 滚动查看docker容器日志，传入容器名（阻塞，Ctrl+C退出）
func TailDockerLog(containerName string) {

	//1.打印提示
	fmt.Printf(">>> 正在滚动查看容器日志: %s (按 Ctrl+C 退出)\n", containerName)
	fmt.Println("--------------------------------------------------------------------------------")

	//2.执行docker logs -f命令
	cmd := exec.Command("docker", "logs", "-f", "--tail", "100", containerName)

	//3.重定向输出到控制台
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//4.监听退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//5.启动命令
	if err := cmd.Start(); err != nil {
		fmt.Printf("查看容器日志失败: %v\n", err)
		return
	}

	//6.等待退出信号或命令结束
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	//7.阻塞等待
	select {
	case <-sigChan:
		cmd.Process.Kill()
		fmt.Println("\n>>> 已退出日志查看")
	case err := <-done:
		if err != nil {
			fmt.Printf("日志查看异常退出: %v\n", err)
		}
	}
}
