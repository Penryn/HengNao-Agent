
IDL_PATH := $(shell pwd)/idl
# 自定义模块的执行命令
.PHONY: all
all: help

default: help

.PHONY: help
help: ## 显示帮助信息，列出所有可用的目标命令。
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Gen
.PHONY: gen
gen: ## 生成代码
	@cwgo server --type HTTP --service meeting_agent --module meeting_agent -I $(IDL_PATH) --idl $(IDL_PATH)/agent.proto

.PHONY: fmt
fmt: ## 格式化 Go 代码，使用 `gofmt`、`gofumpt` 和 `goimports`。
	@gofmt -l -w app
	@gofumpt -l -w app
	@goimports -l -w app

##@ Development Env
.PHONY: env-start
env-start: ## 启动所有中间件服务作为 Docker 容器。
	@docker compose up -d

.PHONY: env-stop
env-stop: ## 停止所有运行中的 Docker 容器。
	@docker compose down