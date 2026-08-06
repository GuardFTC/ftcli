// Package sql @Author:冯铁城 [17615007230@163.com] 2025-10-31 20:12:16
package sql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// importRequest 导入请求参数
type importRequest struct {
	FilePath        string `json:"filePath"`
	Database        string `json:"database"`
	Table           string `json:"table"`
	ColumnSeparator string `json:"columnSeparator"`
	SkipHeader      bool   `json:"skipHeader"`
}

// importResponse 导入响应结果
type importResponse struct {
	Success        bool   `json:"success"`
	TotalRows      int64  `json:"totalRows"`
	LoadedRows     int64  `json:"loadedRows"`
	FilteredRows   int64  `json:"filteredRows"`
	UnselectedRows int64  `json:"unselectedRows"`
	CostTimeMs     int64  `json:"costTimeMs"`
	BatchCount     int    `json:"batchCount"`
	ErrorMessage   string `json:"errorMessage"`
}

// flag变量
var (
	csvFile string
	db      string
	table   string
	path    string
)

// NewSqlCommand 创建sql命令
func NewSqlCommand() *cobra.Command {

	//1.设置Flags
	sqlCmd.Flags().StringVarP(&csvFile, "csvFile", "c", "", "CSV文件名")
	sqlCmd.Flags().StringVarP(&db, "db", "d", "", "数据库")
	sqlCmd.Flags().StringVarP(&table, "table", "t", "", "表")
	sqlCmd.Flags().StringVarP(&path, "path", "p", "", "默认文件路径")

	//2.返回
	return sqlCmd
}

// sql命令 将csv数据通过Stream Load导入Doris
var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "import csv data to doris via stream load",
	Run: func(cmd *cobra.Command, args []string) {
		runCommand()
	},
}

// 运行导入命令
func runCommand() {

	//1.校验CSV文件参数
	if csvFile == "" {
		fmt.Println("错误: 请通过 -c 指定CSV文件名")
		return
	}

	//2.拼装CSV文件绝对路径
	if path == "" {
		path = defaultPath[system]
	}
	filePath := path + csvFile

	//3.校验文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("错误: 文件不存在 %s\n", filePath)
		return
	}

	//4.校验数据库和表参数
	if db == "" {
		fmt.Println("错误: 请通过 -d 指定目标数据库")
		return
	}
	if table == "" {
		fmt.Println("错误: 请通过 -t 指定目标表")
		return
	}

	//5.构建请求体
	reqBody := importRequest{
		FilePath:        filePath,
		Database:        db,
		Table:           table,
		ColumnSeparator: ",",
		SkipHeader:      true,
	}

	//6.调用Stream Load接口
	fmt.Printf("开始导入 | 文件: %s | 目标: %s.%s\n", filePath, db, table)
	result, err := callStreamLoadAPI(reqBody)
	if err != nil {
		fmt.Printf("接口调用失败: %s\n", err)
		return
	}

	//7.打印导入结果
	printResult(result)
}

// 调用Stream Load导入接口
func callStreamLoadAPI(reqBody importRequest) (*importResponse, error) {

	//1.序列化请求体为JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	//2.创建HTTP客户端（超时时间10分钟，适配海量数据导入场景）
	client := &http.Client{Timeout: httpTimeout}

	//3.构建POST请求
	req, err := http.NewRequest(http.MethodPost, streamLoadURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//4.发送请求
	fmt.Println("正在导入中，请耐心等待...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %w", err)
	}
	defer resp.Body.Close()

	//5.校验HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("接口返回异常状态码: %d", resp.StatusCode)
	}

	//6.解析响应体
	var result importResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	//7.返回
	return &result, nil
}

// 打印导入结果
func printResult(result *importResponse) {

	fmt.Println("================================================================================")

	if result.Success {
		fmt.Println("导入成功!")
	} else {
		fmt.Println("导入失败!")
		fmt.Printf("  错误信息: %s\n", result.ErrorMessage)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("  总行数:     %d\n", result.TotalRows)
	fmt.Printf("  加载行数:   %d\n", result.LoadedRows)
	fmt.Printf("  过滤行数:   %d\n", result.FilteredRows)
	fmt.Printf("  批次数:     %d\n", result.BatchCount)
	fmt.Printf("  耗时:       %dms (%.1fs)\n", result.CostTimeMs, float64(result.CostTimeMs)/1000)
	fmt.Println("================================================================================")
}
