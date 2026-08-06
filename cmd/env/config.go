// Package env @Author:冯铁城 [17615007230@163.com] 2025-10-31 19:31:00
package env

import "runtime"

// 定义系统常量
const windows = "windows"
const mac = "darwin"

// 系统名称
var system = runtime.GOOS

// 默认项目
var defaultProject = "prospect-platform"

// env命令 系统-项目名称-项目配置-Map
// 配置格式说明:
//
//	background类型: "background", 日志路径, 检测端口, kill进程名, kill关键字, 实际命令...
//	docker类型: "docker", compose文件路径(单容器为空), 容器名1, 容器名2...
var envCmdProjectPropertiesMap = map[string]map[string]map[string][]string{
	windows: {
		defaultProject: {
			"nacos": {
				"background",
				"C:\\Users\\Administrator\\project\\java\\logs\\nacos.log",
				"8848",
				"java", "nacos",
				"cmd",
				"/C", "C:\\Users\\Administrator\\base\\nacos\\bin\\startup.cmd",
				"-m", "standalone",
			},
			"sentinel": {
				"background",
				"C:\\Users\\Administrator\\project\\java\\logs\\sentinel.log",
				"8849",
				"java", "sentinel",
				"java",
				"-Dserver.port=8849",
				"-Dcsp.sentinel.dashboard.server=0.0.0.0:8849",
				"-Dproject.name=Platform",
				"-Dsentinel.dashboard.auth.username=platform",
				"-Dsentinel.dashboard.auth.password=VI7O8ezi18kaYiQupoT2tohAw4mOLi",
				"-jar", "C:\\Users\\Administrator\\base\\sentinel\\sentinel-dashboard-1.8.8.jar",
			},
			"redis": {
				"docker",
				"",
				"redis-stack",
			},
		},
		"ftcli": {
			"redis": {
				"docker",
				"",
				"redis-stack",
			},
			"chroma": {
				"docker",
				"",
				"chroma-server",
			},
			"es": {
				"docker",
				"C:\\Users\\Administrator\\base\\elastic search\\docker-compose.yaml",
				"es01", "elasticvue",
			},
			"ftcli": {
				"background",
				"C:\\Users\\Administrator\\project\\java\\logs\\ftcli-ai-server.log",
				"6680",
				"java", "ftcli",
				"java",
				"-Dfile.encoding=UTF-8",
				"-Dstdout.encoding=UTF-8",
				"-Dstderr.encoding=UTF-8",
				"-Dhttps.proxyHost=127.0.0.1",
				"-Dhttps.proxyPort=10808",
				"-jar", "C:\\Users\\Administrator\\project\\java\\ftcli-ai-server\\target\\ftcli-ai-server-0.0.1-SNAPSHOT.jar",
			},
		},
		"logging-mon": {
			"kafka": {
				"docker",
				"C:\\Users\\Administrator\\base\\kafka\\docker-compose.yml",
				"kafka-ui", "kafka", "zookeeper",
			},
		},
		"ftc-loader": {
			"ftc-loader": {
				"background",
				"C:\\Users\\Administrator\\project\\java\\logs\\ftcli-doris-stream-loader.log",
				"6677",
				"java", "ftcli-doris-stream-loader",
				"java",
				"-Dfile.encoding=UTF-8",
				"-Dstdout.encoding=UTF-8",
				"-Dstderr.encoding=UTF-8",
				"-jar", "C:\\Users\\Administrator\\project\\java\\ftcli-doris-stream-loader\\target\\ftcli-doris-stream-loader-0.0.1-SNAPSHOT.jar",
			},
		},
	},
	mac: {
		defaultProject: {
			//"nacos": {
			//	"background",
			//	"/Applications/project/java/logs/nacos.log",
			//	"8848",
			//	"java", "nacos",
			//	"bash",
			//	"/Applications/base/nacos/nacos/bin/startup.sh",
			//	"-m", "standalone",
			//},
			//"sentinel": {
			//	"background",
			//	"/Applications/project/java/logs/sentinel.log",
			//	"8849",
			//	"java", "sentinel",
			//	"java",
			//	"-Dserver.port=8849",
			//	"-Dcsp.sentinel.dashboard.server=0.0.0.0:8849",
			//	"-Dproject.name=Platform",
			//	"-Dsentinel.dashboard.auth.username=platform",
			//	"-Dsentinel.dashboard.auth.password=VI7O8ezi18kaYiQupoT2tohAw4mOLi",
			//	"-jar", "/Applications/base/sentinel/sentinel-dashboard-1.8.8.jar",
			//},
			//"redis": {
			//	"docker",
			//	"",
			//	"redis-dev",
			//},
		},
		"ftcli": {
			//"redis": {
			//	"docker",
			//	"",
			//	"redis-dev",
			//},
			//"chroma": {
			//	"docker",
			//	"",
			//	"chroma-server",
			//},
			"ftcli": {
				"background",
				"/Applications/project/java/logs/ftcli-ai-server.log",
				"6680",
				"java", "ftcli",
				"java",
				"-Dfile.encoding=UTF-8",
				"-Dstdout.encoding=UTF-8",
				"-Dstderr.encoding=UTF-8",
				"-jar", "/Applications/project/java/ftcli-ai-server/target/ftcli-ai-server-0.0.1-SNAPSHOT.jar",
			},
		},
		"dolp": {
			"zookeeper": {
				"background",
				"/Applications/project/java/logs/zookeeper.log",
				"2181",
				"java", "zookeeper",
				"bash",
				"/Applications/base/zookeeper/apache-zookeeper-3.8.6-bin/bin/zkServer.sh",
				"start",
			},
		},
		"ftc-loader": {
			"ftc-loader": {
				"background",
				"/Applications/project/java/logs/ftcli-doris-stream-loader.log",
				"6677",
				"java", "ftcli-doris-stream-loader",
				"java",
				"-Dfile.encoding=UTF-8",
				"-Dstdout.encoding=UTF-8",
				"-Dstderr.encoding=UTF-8",
				"-jar", "/Applications/project/java/ftcli-doris-stream-loader/target/ftcli-doris-stream-loader-0.0.1-SNAPSHOT.jar",
			},
		},
	},
}
