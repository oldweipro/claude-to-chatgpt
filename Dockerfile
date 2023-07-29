FROM golang:alpine as builder

WORKDIR /go/src/github.com/oldweipro/claude-to-chatgpt
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o claude-to-chatgpt .

FROM alpine:latest

LABEL MAINTAINER="oldwei@oldweipro@gmail.com"

WORKDIR /go/src/github.com/oldweipro/claude-to-chatgpt

COPY --from=0 /go/src/github.com/oldweipro/claude-to-chatgpt ./
COPY --from=0 /go/src/github.com/oldweipro/claude-to-chatgpt/config.yaml ./

EXPOSE 8787
ENTRYPOINT ./claude-to-chatgpt -c config.yaml
