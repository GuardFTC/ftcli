// Package wmai @Author:冯铁城 [17615007230@163.com] 2026-07-16 16:00:00
package wmai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// flag变量
var (
	stat   bool
	apiKey string
)

// UsageResponse 用量查询响应体
type UsageResponse struct {
	Key              string  `json:"key"`
	Timezone         string  `json:"timezone"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	Cost             float64 `json:"cost"`
	Currency         string  `json:"currency"`
	RequestCount     int     `json:"request_count"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
}

// ErrorResponse 错误响应体
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	} `json:"error"`
}

// NewWmaiCommand 创建wmai命令
func NewWmaiCommand() *cobra.Command {

	//1.设置Flags
	wmaiCmd.Flags().BoolVarP(&stat, "stat", "s", false, "查询今日API Key用量统计")
	wmaiCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key（默认读取环境变量WM_AI_KEY）")

	//2.返回
	return wmaiCmd
}

// wmaiCmd wmai命令
var wmaiCmd = &cobra.Command{
	Use:   "wmai",
	Short: "Wanmei AI API Key usage statistics",
	Run: func(cmd *cobra.Command, args []string) {

		//1.根据flag执行对应操作
		switch {
		case stat:
			runStat()
		default:
			cmd.Help()
		}
	},
}

// runStat 查询今日用量统计
func runStat() {

	//1.获取API Key
	key := getAPIKey()
	if key == "" {
		fmt.Println("错误: 未找到API Key，请通过 -k 参数传入或设置环境变量 WM_AI_KEY")
		return
	}

	//2.计算今天00:00和当前时间（北京时间）
	beijing := time.FixedZone("CST", 8*3600)
	now := time.Now().In(beijing)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, beijing)

	//3.格式化时间
	startStr := start.Format("2006-01-02T15:04:05")
	endStr := now.Format("2006-01-02T15:04:05")

	//4.发送请求
	result, err := doUsageRequest(key, startStr, endStr)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	//5.格式化输出
	consoleUsage(result.Key, startStr, endStr, result.Cost, result.RequestCount, result.PromptTokens, result.CompletionTokens, result.TotalTokens)
}

// getAPIKey 获取API Key（flag优先，环境变量次之）
func getAPIKey() string {

	//1.flag优先
	if apiKey != "" {
		return apiKey
	}

	//2.环境变量次之
	return os.Getenv(envKeyName)
}

// doUsageRequest 发送用量查询请求
func doUsageRequest(key string, start string, end string) (*UsageResponse, error) {

	//1.构建URL
	url := defaultBaseURL + "/v1/key/usage?start_time=" + start + "&end_time=" + end

	//2.创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	//3.设置请求头
	req.Header.Set("Authorization", "Bearer "+key)

	//4.发送请求
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	//5.读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	//6.处理错误响应
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return nil, fmt.Errorf("[%s] %s: %s", errResp.Error.Code, errResp.Error.Message, errResp.Error.Details)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	//7.解析响应
	var result UsageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	//8.返回
	return &result, nil
}

// consoleUsage 格式化输出用量信息
func consoleUsage(key string, start string, end string, cost float64, requests int, prompt int, completion int, total int) {

	//1.格式化时间显示（去掉T，用空格替代）
	displayStart := strings.ReplaceAll(start, "T", " ")
	displayEnd := strings.ReplaceAll(end, "T", " ")

	//2.计算使用百分比
	percent := cost / dailyLimit * 100
	if percent > 100 {
		percent = 100
	}

	//3.构建进度条（24格）
	barLen := 24
	filled := int(percent / 100 * float64(barLen))
	if filled > barLen {
		filled = barLen
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)

	//4.打印标题
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("完美 AI 使用额度")
	fmt.Println(strings.Repeat("=", 80))

	//5.打印基本信息
	fmt.Printf("密钥 (Key)   :  %s\n", key)
	fmt.Printf("时间范围     :  %s ~ %s (今日)\n", displayStart, displayEnd)

	//6.打印分隔线
	fmt.Println(strings.Repeat("-", 80))

	//7.打印花费与进度条
	fmt.Printf("消费金额     :  $%.6f / $%.2f  [%s]  %.2f%%\n", cost, dailyLimit, bar, percent)
	fmt.Printf("请求次数     :  %d\n", requests)

	//8.打印分隔线
	fmt.Println(strings.Repeat("-", 80))

	//9.打印Token信息
	fmt.Printf("输入 (Prompt):  %s tokens\n", formatNumber(prompt))
	fmt.Printf("输出 (Comp.) :  %s tokens\n", formatNumber(completion))
	fmt.Printf("总计 (Total) :  %s tokens\n", formatNumber(total))

	//10.打印底部分隔线
	fmt.Println(strings.Repeat("=", 80))
}

// formatNumber 格式化数字（添加千位分隔符）
func formatNumber(n int) string {

	//1.转为字符串
	s := fmt.Sprintf("%d", n)

	//2.如果小于4位，直接返回
	if len(s) <= 3 {
		return s
	}

	//3.添加千位分隔符
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}

	//4.返回
	return string(result)
}
