// Package ai @Author:冯铁城 [17615007230@163.com] 2026-06-03 16:00:00
package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

// ChatPayload 聊天请求体
type ChatPayload struct {
	IsLocal     bool   `json:"isLocal"`
	ChatId      string `json:"chatId"`
	UserMessage string `json:"userMessage"`
}

// runChat 进入聊天交互模式
func runChat(isLocal bool) {

	//1.确定聊天模式名称
	modeName := "Web AI 智能检索"
	if isLocal {
		modeName = "Local AI 本地文库"
	}

	//2.获取会话ID
	chatId, err := getChatId()
	if err != nil {
		fmt.Printf("获取会话ID失败: %v\n", err)
		return
	}

	//3.初始化readline实例(支持方向键翻历史、行编辑)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "[User] > ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("初始化readline失败: %v\n", err)
		return
	}
	defer rl.Close()

	//4.监听Ctrl+C信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	//5.打印进入提示
	fmt.Println("===================================================================")
	fmt.Printf(">> 进入 %s 聊天模式 (输入 'exit' 或 按 Ctrl+C 退出) <<\n", modeName)
	fmt.Println("===================================================================")
	fmt.Println()

	//6.聊天循环
	for {

		//7.监听输入或退出信号
		inputChan := make(chan string, 1)
		errChan := make(chan error, 1)
		go func() {
			line, err := rl.Readline()
			if err != nil {
				errChan <- err
			} else {
				inputChan <- line
			}
		}()

		//8.等待输入或信号
		var input string
		select {
		case <-sigChan:
			input = "exit"
		case err := <-errChan:
			if errors.Is(err, readline.ErrInterrupt) || err == io.EOF {
				input = "exit"
			} else {
				fmt.Printf("\n读取输入失败: %v\n", err)
				return
			}
		case input = <-inputChan:
		}

		//9.处理退出
		input = strings.TrimSpace(input)
		if input == "exit" || input == "quit" {
			fmt.Println()
			fmt.Println("===================================================================")
			fmt.Println(">> 已安全退出 AI 聊天模式。感谢使用！")
			fmt.Println("===================================================================")
			return
		}

		//10.空输入跳过
		if input == "" {
			continue
		}

		//11.清屏命令
		if input == "clear" {
			clearScreen()
			fmt.Println("===================================================================")
			fmt.Printf(">> 进入 %s 聊天模式 (输入 'exit' 或 按 Ctrl+C 退出) <<\n", modeName)
			fmt.Println("===================================================================")
			fmt.Println()
			continue
		}

		//12.构建请求
		payload := ChatPayload{
			IsLocal:     isLocal,
			ChatId:      chatId,
			UserMessage: input,
		}

		//13.发送流式请求，打字机效果输出
		fmt.Print("[AI]   > ")
		err := doPostStream(baseURL+"/api/rest/v1/ai/chat/stream", payload, func(chunk string) {
			fmt.Print(cleanMarkdown(chunk))
		})

		//14.处理错误
		if err != nil {
			fmt.Printf("\n请求失败: %v\n", err)
		} else {
			fmt.Println()
		}
		fmt.Println()
	}
}

// getChatId 获取会话ID
func getChatId() (string, error) {

	//1.发送请求
	result, err := doGet(baseURL + "/api/rest/v1/ai/chatId")
	if err != nil {
		return "", err
	}

	//2.解析chatId
	var chatId string
	if err := json.Unmarshal(result.Data, &chatId); err != nil {
		return "", fmt.Errorf("解析会话ID失败: %v", err)
	}

	//3.返回
	return chatId, nil
}

// clearScreen 清屏
func clearScreen() {

	//1.根据系统选择不同命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	//2.绑定标准输出
	cmd.Stdout = os.Stdout

	//3.执行
	cmd.Run()
}
