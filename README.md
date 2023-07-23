# Claude 到 ChatGPT 的接口适配

简体中文 | [English](README.en_US.md)

## 介绍

本项目是把 [Claude](https://claude.ai) 聊天功能接口适配到 OpenAI API 标准接口。

当启动本项目后,就可以按照 [v1/chat/completions](https://platform.openai.com/docs/api-reference/chat)
的接口文档调用本项目接口 `http://127.0.0.1:8080/v1/chat/completions`
得到与 [OpenAI API](https://platform.openai.com/docs/api-reference/chat)
相同的数据结构,方便已经对 [OpenAI API](https://platform.openai.com/docs/api-reference/chat) 结果进行开发的用户快速切换。

## 运行环境

需要 [Go](https://go.dev/dl/) 1.20 及以上版本。

## 获取源码

```
git clone https://github.com/oldweipro/claude-to-chatgpt.git
```

## 运行

进入项目目录

```
cd claude-to-chatgpt
```

获取依赖

```shell
# 整理go.mod
go mod tidy

# 下载go.mod依赖
go mod download
```

运行

```shell
go run main.go
```

使用 `-c` 指定配置文件 `config-dev.yaml`

```shell
go run main.go -c config-dev.yaml
```

## 配置

配置文件如果不存在,程序会自动创建。

如果启动后填写的配置信息有误,直接修改配置文件并保存即可,程序会自动重新加载。

| 配置项                 | 说明                 | 示例值                  |
|---------------------|--------------------|----------------------|  
| claude              | Claude 相关配置        |                      |
| - base-url          | Claude服务地址,可选      | https://claude.ai    |
| - session-key       | 当前对话session唯一标识,必填 | sk-ant-xxxx-8MgHUgAA | 
| - organization-uuid | 组织唯一ID,可选          |                      |
| proxy               | 代理配置,可选            |                      |
| - protocol          | 协议                 | http                 |
| - host              | 代理服务器地址            | 127.0.0.1            |
| - port              | 代理服务器端口            | 7890                 |
| - username          | 认证用户名              |                      |
| - password          | 认证密码               |                      |

## 部署

### 编译

可针对不同平台编译生成可执行文件。

Windows:

```
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
```

Linux:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 
```

macOS:

```
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
```

### 运行

将编译好的可执行文件复制到服务器相应目录,赋予执行权限并运行。

## Go编译命令参数说明

| 参数          | 说明      | 可选值                      |
|-------------|---------|--------------------------|
| CGO_ENABLED | 是否启用Cgo | 0: 禁用Cgo<br>1: 启用Cgo(默认) |  
| GOOS        | 目标操作系统  | linux, windows, darwin等  |
| GOARCH      | 目标架构    | amd64, 386, arm等         |
| go build    | 执行Go编译  |                          |