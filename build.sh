#!/bin/bash
set -e

echo "🚀 连续点击器 - 智能构建脚本"
echo "================================"

# 项目根目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
DIST_DIR="$ROOT_DIR/dist"
FONT_FILE="$ROOT_DIR/SourceHanSansSC-Bold.otf"

# 输出清理
echo "🧹 清理旧的 dist 目录..."
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# 检查字体文件
if [ ! -f "$FONT_FILE" ]; then
    echo "⚠️ 字体文件未找到: $FONT_FILE"
else
    echo "✅ 字体文件存在: $FONT_FILE"
fi

# 检测当前系统
CURRENT_OS=$(uname -s)
CURRENT_ARCH=$(uname -m)

echo "🔍 检测到当前系统: $CURRENT_OS ($CURRENT_ARCH)"

# 根据当前系统编译对应版本
case "$CURRENT_OS" in
    "Darwin")
        echo "📱 检测到 macOS 系统，编译 macOS 版本..."
        
        # 检测架构
        if [ "$CURRENT_ARCH" = "arm64" ]; then
            echo "🔄 编译 Apple Silicon (ARM64) 版本..."
            CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/clicker_mac_arm64" .
            echo "✅ Apple Silicon 版本编译完成"
        else
            echo "🔄 编译 Intel (AMD64) 版本..."
            CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/clicker_mac_amd64" .
            echo "✅ Intel 版本编译完成"
        fi
        
        # 尝试编译另一个架构（如果可能）
        if [ "$CURRENT_ARCH" = "arm64" ]; then
            echo "🔄 尝试编译 Intel (AMD64) 版本..."
            if CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/clicker_mac_amd64" . 2>/dev/null; then
                echo "✅ Intel 版本编译完成"
            else
                echo "⚠️ Intel 版本编译失败（可能需要 Rosetta 2）"
            fi
        else
            echo "🔄 尝试编译 Apple Silicon (ARM64) 版本..."
            if CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/clicker_mac_arm64" . 2>/dev/null; then
                echo "✅ Apple Silicon 版本编译完成"
            else
                echo "⚠️ Apple Silicon 版本编译失败"
            fi
        fi
        ;;
        
    "Linux")
        echo "🐧 检测到 Linux 系统，编译 Linux 版本..."
        
        if [ "$CURRENT_ARCH" = "aarch64" ]; then
            echo "🔄 编译 ARM64 版本..."
            CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o "$DIST_DIR/clicker_linux_arm64" .
        else
            echo "🔄 编译 AMD64 版本..."
            CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/clicker_linux_amd64" .
        fi
        
        echo "✅ Linux 版本编译完成"
        ;;
        
    "MINGW"*|"MSYS"*|"CYGWIN"*)
        echo "🪟 检测到 Windows 系统，编译 Windows 版本..."
        
        if [ "$CURRENT_ARCH" = "x86_64" ]; then
            echo "🔄 编译 AMD64 版本..."
            CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/clicker_win_amd64.exe" .
        else
            echo "🔄 编译 ARM64 版本..."
            CGO_ENABLED=1 GOOS=windows GOARCH=arm64 go build -o "$DIST_DIR/clicker_win_arm64.exe" .
        fi
        
        echo "✅ Windows 版本编译完成"
        ;;
        
    *)
        echo "❓ 未知操作系统: $CURRENT_OS"
        echo "   尝试通用编译..."
        go build -o "$DIST_DIR/clicker_unknown" .
        echo "✅ 通用版本编译完成"
        ;;
esac

echo ""
echo "✅ 构建完成！"
echo "📁 输出目录: $DIST_DIR"
echo ""
echo "📝 生成的文件："
ls -la "$DIST_DIR"/

echo ""
echo "💡 使用说明："
echo "   - 在 macOS 上运行: ./dist/clicker_mac_*"
echo "   - 在 Linux 上运行: ./dist/clicker_linux_*"
echo "   - 在 Windows 上运行: ./dist/clicker_win_*.exe"