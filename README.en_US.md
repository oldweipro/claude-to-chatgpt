# Claude to ChatGPT API Adaptation

[简体中文](README.md) | English

## Introduction

This project adapts the chat functionality interface of [Claude](https://claude.ai) to the standard OpenAI API
interfaces.

After starting this project, you can call the interface `http://127.0.0.1:8080/v1/chat/completions` of this project
according to the interface documentation of [v1/chat/completions](https://platform.openai.com/docs/api-reference/chat)
to get the same data structure returned by [OpenAI API](https://platform.openai.com/docs/api-reference/chat). This
facilitates users who have developed based on the results
of [OpenAI API](https://platform.openai.com/docs/api-reference/chat) to quickly switch over.

## Runtime Environment

Requires [Go](https://go.dev/dl/) version 1.20 or above.

## Get Source Code

```
git clone https://github.com/oldweipro/claude-to-chatgpt.git
```

## Run

Enter project directory

```
cd claude-to-chatgpt
```

Get dependencies

```shell
# Tidy go.mod
go mod tidy

# Download dependencies in go.mod
go mod download
```

Run

```shell
go run main.go
```

Use `-c` to specify the configuration file `config-dev.yaml`

```shell
go run main.go -c config-dev.yaml
```

## Configuration

If the configuration file does not exist, the program will create it automatically.

If the configuration information filled in after startup is incorrect, just modify the configuration file directly and
save it. The program will automatically reload.

| Configuration       | Description                                                 | Example Value        |
|---------------------|-------------------------------------------------------------|----------------------|
| claude              | Claude related configuration                                |                      |  
| - base-url          | Claude service address, optional                            | https://claude.ai    |
| - session-key       | Unique identifier of current conversation session, required | sk-ant-xxxx-8MgHUgAA |
| - organization-uuid | Organization unique ID, optional                            |                      |
| proxy               | Proxy configuration, optional                               |                      |
| - protocol          | Protocol                                                    | http                 | 
| - host              | Proxy server address                                        | 127.0.0.1            |
| - port              | Proxy server port                                           | 7890                 |
| - username          | Authentication username                                     | admin                |
| - password          | Authentication password                                     | admin                |

## Deployment

### Compile

You can compile executable files for different platforms.

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

### Run

Copy the compiled executable file to the corresponding server directory, grant execution permissions and run it.

## Go Compile Command Parameters

| Parameter   | Description             | Optional Values                           |
|-------------|-------------------------|-------------------------------------------|
| CGO_ENABLED | Whether to enable Cgo   | 0: Disable Cgo<br>1: Enable Cgo (default) |
| GOOS        | Target operating system | linux, windows, darwin etc.               |
| GOARCH      | Target architecture     | amd64, 386, arm etc.                      | 
| go build    | Execute Go compile      |                                           |