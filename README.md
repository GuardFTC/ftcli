# ftcli

个人开发 CLI 工具，用于自动化日常开发工作流：环境启动、项目构建、Maven 打包、CSV 导入 Doris、AI 助手、常用软件一键打开。

```
ftcli/
├── main.go              # 入口
├── cmd/
│   ├── root.go          # 根命令 ftcli (v1.0.0)
│   ├── ai/              # AI 助手（对接 ftcli 后端）
│   ├── build/           # 构建流水线（kill → package → start）
│   ├── env/             # 开发环境启动
│   ├── monitor/         # 系统资源监控（内存/CPU）
│   ├── open/            # 常用软件启动
│   ├── package/         # Java Maven 打包
│   ├── sql/             # CSV → Doris Stream Load 导入
│   └── wmai/            # 完美 AI API Key 用量统计
└── util/                # 公共工具（命令执行、进程管理、Docker、浏览器打开等）
```

## 环境要求

- Go 1.24+
- Maven 3.x（Java 打包/构建需要）
- Docker（部分环境服务依赖容器）
- Windows + macOS 双平台支持，各命令配置按 `runtime.GOOS` 自动区分

## 安装与构建

```bash
# 本地运行
go run .

# 编译为可执行文件
go build -o ftcli.exe .

# 通过 build 命令自举编译（输出到 ../../bin/ftcli.exe）
ftcli build -p ftcli -t go
```

## 命令总览

| 命令 | 说明 |
|------|------|
| `ftcli env` | 启动项目开发环境 |
| `ftcli package` | Java Maven 打包 |
| `ftcli build` | 完整构建流水线：kill → 打包 → 后台启动 |
| `ftcli sql` | CSV 数据通过 Stream Load 导入 Doris |
| `ftcli ai` | AI 助手（流式聊天 / 文档上传 / 管理页面） |
| `ftcli open` | 一键打开常用开发软件 |
| `ftcli monitor` | 系统资源监控（内存/CPU） |
| `ftcli wmai` | 完美 AI API Key 用量统计 |

---

## ftcli env

启动项目所需的中间件和后端服务。支持 Docker 容器启动与后台进程启动（含端口检测、进程幂等 kill）。

```bash
ftcli env                  # 启动默认项目（prospect-platform）
ftcli env -p ftcli         # 启动 ftcli 项目环境
ftcli env -l               # 列出内置项目及配置
ftcli env -b               # 查看所有后台服务运行状态
ftcli env --bl ftcli       # 滚动查看指定服务日志（支持服务名或容器名）
ftcli env --blc chroma     # 强制按容器名查看 docker 日志（绕过服务名匹配）
ftcli env --bk ftcli       # 停止指定后台服务
```

> `--bl` 按服务名匹配日志：background 类型查看日志文件，docker 类型查看容器日志。当服务名与容器名冲突时（如 kafka 服务包含 kafka 容器），用 `--blc` 直接按容器名查看。

**内置项目（Windows）**

| 项目 | 启动服务 |
|------|----------|
| `prospect-platform` | Nacos (8848)、Sentinel (8849)、Redis (Docker) |
| `ftcli` | Redis (Docker)、Chroma (Docker)、ES + ElasticVue (Docker Compose)、ftcli 后端 (6680) |
| `logging-mon` | Kafka + Kafka UI + Zookeeper (Docker Compose) |
| `ftc-loader` | ftcli-doris-stream-loader (6677) |

**内置项目（macOS）**

| 项目 | 启动服务 |
|------|----------|
| `ftcli` | ftcli 后端 (6680) |
| `ftc-loader` | ftcli-doris-stream-loader (6677) |

---

## ftcli package

执行 Maven `clean → install → package`，打包前自动 kill 相关 Java 进程，完成后打开输出目录。

```bash
ftcli package              # 打包默认项目（prospect-platform）
ftcli package -p ftcli     # 打包 ftcli 后端
ftcli package -p logging-mon
ftcli package -l           # 列出内置项目

# 手动指定路径（项目未在内置列表时）
ftcli package -P <pom路径> -m <settings路径> -o <输出目录>
```

**内置项目（Windows + macOS 双平台）**：`prospect-platform`、`logging-mon`、`ftcli`、`ftc-loader`

打包流程：kill 相关 Java 进程 → 依次执行 `mvn clean`、`mvn install`、`mvn package`（均带 `-DskipTests=true`）→ 打开输出目录（Windows 用 `explorer`，macOS 用 `open`）

---

## ftcli build

一键完成 kill → Maven 打包 → 后台启动。支持 Java、Go 及混合构建。

```bash
ftcli build                # 默认：ftcli 项目，java + go 全部构建
ftcli build -p ftcli -t java
ftcli build -p ftcli -t go
ftcli build -p ftcli -t all
ftcli build -l             # 列出内置项目及支持类型
```

**构建流程**

- **Java**：kill 进程 → Maven 打包（clean → install → package）→ 后台启动（端口存活检测，超时 15 秒）
- **Go**：
  - Windows：编译到临时文件 `ftcli_new.exe`，生成 `ftcli_replace.bat` 脚本，当前进程退出后延迟自动替换目标 exe
  - macOS/Unix：直接 `go build` 覆盖目标文件

**内置项目**：`ftcli`（支持 java / go / all，Windows + macOS 双平台）、`ftc-loader`（仅 java）

---

## ftcli sql

调用 Doris Stream Load 服务接口，将本地 CSV 文件高效导入到 Doris 指定库表中。支持千万级数据量，导入过程流式处理，无 OOM 风险。

```bash
ftcli sql -c data.csv -d dw_tile -t ads_bi_af_ltvroas_d_i
ftcli sql -c data.csv -d dw_tile -t ads_bi_af_ltvroas_d_i -p /data/csv/
```

**参数说明**

| 参数 | 说明 |
|------|------|
| `-c` | CSV 文件名（必填） |
| `-d` | 目标数据库（必填） |
| `-t` | 目标表（必填） |
| `-p` | CSV 文件所在目录（可选，默认为系统 Downloads 目录） |

**默认 CSV 目录**

| 系统 | 路径 |
|------|------|
| Windows | `C:\Users\Administrator\Downloads\` |
| macOS | `/Users/m/Downloads/` |

**依赖服务**：需要 `ftcli-doris-stream-loader` 服务运行在 `localhost:6677`（HTTP 超时 10 分钟）。

**输出示例**

```
开始导入 | 文件: /Users/m/Downloads/data.csv | 目标: dw_tile.ads_bi_af_ltvroas_d_i
正在导入中，请耐心等待...
================================================================================
导入成功!
--------------------------------------------------------------------------------
  总行数:     4999960
  加载行数:   4950000
  过滤行数:   49960
  批次数:     100
  耗时:       32846ms (32.8s)
================================================================================
```

---

## ftcli ai

对接 ftcli Java 后端服务（默认 `http://localhost:6680`），提供流式聊天、文档上传、管理页面入口。

```bash
ftcli ai -l                # 本地文库流式聊天
ftcli ai -w                # 网络检索流式聊天
ftcli ai -u <路径>         # 上传文档到知识库
ftcli ai -f                # 浏览器打开文档管理页面
ftcli ai -t                # 浏览器打开工具管理页面
ftcli ai -s                # 浏览器打开技能管理页面
ftcli ai -S <地址>         # 指定后端地址（默认 localhost:6680）
```

聊天模式：输入 `exit` 或 `Ctrl+C` 退出，`clear` 清屏。输出自动做 Markdown → ANSI 颜色转换（标题黄色、代码青色、引用灰色）。

---

## ftcli open

一键启动日常开发软件。无参数时启动所有 `Always=true` 的软件，也可按名称指定。

```bash
ftcli open                 # 启动所有默认软件
ftcli open -l              # 列出支持的软件及默认启动状态
ftcli open goland          # 启动指定软件（空格分隔多个）
ftcli open goland webstorm
```

**内置软件（Windows）**

| 默认启动 | 软件 |
|----------|------|
| ✓ | edge、v2ray、docker、wechat、idea、goland、datagrip、kiro、yuque、typora、sublime |
| ✗ | chrome、we、webstorm、cursor、apipost、virtual、RDM、draw.io |

**内置软件（macOS）**

| 默认启动 | 软件 |
|----------|------|
| ✓ | v2ray、ishot、kh、edge、wechat、we、idea、datagrip、kiro、yuque、typora、sublime |
| ✗ | goland、webstorm、arm、apipost、docker、draw.io、tabby |

---

## ftcli monitor

系统资源监控，输出内存和 CPU 使用情况，带进度条和多核利用率分列展示。

```bash
ftcli monitor              # 输出内存 + CPU 全部信息
ftcli monitor -m           # 仅输出内存信息
ftcli monitor -c           # 仅输出 CPU 信息
```

**输出示例**

```
=== MEMORY MONITOR =============================================================
Total      :   31.77 GB  [████████████████░░░░░░░░]  80.00%
Used       :   25.61 GB
Available  :    6.16 GB
=== CPU MONITOR ================================================================
Cores (P/L): 12 / 20                                       Total Usage:   2.01%
Core 00:   0.52% | Core 05:   0.00% | Core 10:   1.04% | Core 15:   0.00%
Core 01:   0.00% | Core 06:   0.00% | Core 11:   0.00% | Core 16:   5.18%
Core 02:   0.00% | Core 07:   0.00% | Core 12:   0.00% | Core 17:   0.52%
Core 03:   0.00% | Core 08:   4.66% | Core 13:   0.00% | Core 18:  15.03%
Core 04:   0.00% | Core 09:   0.00% | Core 14:   1.04% | Core 19:  12.44%
================================================================================
```

---

## ftcli wmai

查询完美 AI（Wanmei AI）API Key 的今日用量统计，包括消费金额、请求次数、Token 用量，带进度条可视化额度消耗。

```bash
ftcli wmai -s                   # 查询今日用量（读取环境变量 WM_AI_KEY）
ftcli wmai -s -k <API Key>      # 指定 API Key 查询
```

**配置**

| 配置项 | 值 |
|--------|------|
| API 地址 | `https://api.ai.wanmei.net` |
| 环境变量 | `WM_AI_KEY` |
| 每日额度上限 | $20.00 |
| 时区 | 北京时间（CST, UTC+8） |

**输出示例**

```
================================================================================
完美 AI 使用额度
================================================================================
密钥 (Key)   :  sk-xxxxxx
时间范围     :  2026-07-28 00:00:00 ~ 2026-07-28 16:30:00 (今日)
--------------------------------------------------------------------------------
消费金额     :  $3.142000 / $20.00  [████░░░░░░░░░░░░░░░░░░░░]  15.71%
请求次数     :  42
--------------------------------------------------------------------------------
输入 (Prompt):  1,234 tokens
输出 (Comp.) :  567 tokens
总计 (Total) :  1,801 tokens
================================================================================
```

---

## 配置说明

各子命令的配置集中在对应 `config.go` 中，按 `runtime.GOOS` 区分系统。修改路径、端口、进程关键字等直接编辑即可，无需改业务逻辑。

| 模块 | 配置文件 |
|------|----------|
| env | `cmd/env/config.go` |
| package | `cmd/package/config.go` |
| build | `cmd/build/config.go` |
| sql | `cmd/sql/config.go` |
| ai | `cmd/ai/config.go` |
| open | `cmd/open/config.go` |
| wmai | `cmd/wmai/config.go` |

---

## 工具层（util）

公共工具函数位于 `util/` 目录，供各子命令复用：

| 文件 | 函数 | 说明 |
|------|------|------|
| `cmd.go` | `RunCommand` / `RunCommandInDir` | 前台执行命令，输出实时打印到控制台 |
| | `RunCommandBackground` | 后台启动进程，日志输出到文件，端口存活检测（超时 15 秒） |
| | `RunCommandBackgroundNoCheck` | 后台启动进程，不做存活检测（用于启动 GUI 软件） |
| | `ensureUTF8Console` | Windows 下设置 UTF-8 代码页（chcp 65001），解决中文乱码 |
| `normal.go` | `KillProcess` / `IsProcessRunning` | 按进程名 + 命令行关键字 kill/检查进程 |
| | `OpenBrowser` | 跨平台用默认浏览器打开 URL |
| | `IsNumeric` / `GetProjectItems` | 数据类型判断、项目配置项读取 |
| `docker.go` | `IsDockerContainerRunning` / `StopDockerContainer` | Docker 容器状态检查/停止 |
| | `StopDockerCompose` | 停止整个 Docker Compose 组 |
| | `TailDockerLog` | 滚动查看容器日志（`docker logs -f`，Ctrl+C 退出） |
| `proc_windows.go` | `setDetachAttrs` | Windows：创建新进程组（`CREATE_NEW_PROCESS_GROUP`），子进程脱离父终端 |
| `proc_unix.go` | `setDetachAttrs` | Unix：创建新会话（`Setsid`），子进程脱离父终端 |

---

## 依赖

| 库 | 用途 |
|----|------|
| [cobra](https://github.com/spf13/cobra) | CLI 框架 |
| [gopsutil/v3](https://github.com/shirou/gopsutil) | 进程管理（kill、端口检测） |
| [gopsutil/v4](https://github.com/shirou/gopsutil) | 系统资源监控（内存/CPU） |

---

## 作者

冯铁城 — 17615007230@163.com
