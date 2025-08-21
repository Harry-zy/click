# 连续点击器 Makefile

# 变量定义
BINARY_NAME=click
BINARY_NAME_WINDOWS=click.exe
MAIN_FILES=main.go system_clicker.go config.go logger.go
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)

# 默认目标
.PHONY: all
all: clean build

# 构建程序
.PHONY: build
build:
	@echo "正在构建连续点击器..."
	go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o ${BINARY_NAME} ${MAIN_FILES}
	@echo "构建完成: ${BINARY_NAME}"

# 构建Windows版本
.PHONY: build-windows
build-windows:
	@echo "正在构建Windows版本..."
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o ${BINARY_NAME_WINDOWS} ${MAIN_FILES}
	@echo "构建完成: ${BINARY_NAME_WINDOWS}"

# 构建macOS版本
.PHONY: build-macos
build-macos:
	@echo "正在构建macOS版本..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" -o ${BINARY_NAME}-macos ${MAIN_FILES}
	@echo "构建完成: ${BINARY_NAME}-macos"

# 构建Linux版本
.PHONY: build-linux
build-linux:
	@echo "正在构建Linux版本..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUERSION}" -o ${BINARY_NAME}-linux ${MAIN_FILES}
	@echo "构建完成: ${BINARY_NAME}-linux"

# 运行程序
.PHONY: run
run:
	@echo "正在运行程序..."
	go run ${MAIN_FILES}

# 测试程序
.PHONY: test
test:
	@echo "正在运行测试..."
	go test ./...

# 清理构建文件
.PHONY: clean
clean:
	@echo "正在清理构建文件..."
	rm -f ${BINARY_NAME} ${BINARY_NAME_WINDOWS} ${BINARY_NAME}-macos ${BINARY_NAME}-linux
	@echo "清理完成"

# 安装依赖
.PHONY: deps
deps:
	@echo "正在安装依赖..."
	go mod tidy
	go mod download
	@echo "依赖安装完成"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "正在格式化代码..."
	go fmt ./...
	@echo "代码格式化完成"

# 代码检查
.PHONY: lint
lint:
	@echo "正在检查代码..."
	golangci-lint run
	@echo "代码检查完成"

# 帮助信息
.PHONY: help
help:
	@echo "连续点击器 Makefile 使用说明:"
	@echo ""
	@echo "可用目标:"
	@echo "  all          - 清理并构建程序"
	@echo "  build        - 构建程序"
	@echo "  build-windows- 构建Windows版本"
	@echo "  build-macos  - 构建macOS版本"
	@echo "  build-linux  - 构建Linux版本"
	@echo "  run          - 运行程序"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  deps         - 安装依赖"
	@echo "  fmt          - 格式化代码"
	@echo "  lint         - 代码检查"
	@echo "  help         - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build        # 构建程序"
	@echo "  make run          # 运行程序"
	@echo "  make build-windows # 构建Windows版本"
