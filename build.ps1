#!/usr/bin/env pwsh

Write-Host "🚀 连续点击器 - Windows PowerShell 构建脚本" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green

# 项目根目录
$ROOT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
$DIST_DIR = Join-Path $ROOT_DIR "dist"
$FONT_FILE = Join-Path $ROOT_DIR "SourceHanSansSC-Bold.otf"

# 输出清理
Write-Host "🧹 清理旧的 dist 目录..." -ForegroundColor Yellow
if (Test-Path $DIST_DIR) {
    Remove-Item $DIST_DIR -Recurse -Force
}
New-Item -ItemType Directory -Path $DIST_DIR -Force | Out-Null

# 检查字体文件
if (Test-Path $FONT_FILE) {
    Write-Host "✅ 字体文件存在: $FONT_FILE" -ForegroundColor Green
} else {
    Write-Host "⚠️ 字体文件未找到: $FONT_FILE" -ForegroundColor Yellow
}

# 检测当前架构
Write-Host "🔍 检测系统架构..." -ForegroundColor Cyan
if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
    Write-Host "🖥️ 检测到 AMD64 架构" -ForegroundColor Green
    $GOARCH = "amd64"
} elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    Write-Host "🖥️ 检测到 ARM64 架构" -ForegroundColor Green
    $GOARCH = "arm64"
} else {
    Write-Host "🖥️ 检测到 $($env:PROCESSOR_ARCHITECTURE) 架构" -ForegroundColor Yellow
    $GOARCH = "amd64"
}

# 设置环境变量
Write-Host "🔧 设置编译环境..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:CGO_ENABLED = "1"

# 编译项目
Write-Host "🚀 开始编译 Windows 版本..." -ForegroundColor Green
Write-Host "📱 目标系统: Windows ($GOARCH)" -ForegroundColor Cyan

$OUTPUT_FILE = Join-Path $DIST_DIR "clicker_win_$GOARCH.exe"

try {
    & go build -o $OUTPUT_FILE .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Windows 版本编译完成" -ForegroundColor Green
    } else {
        Write-Host "❌ 编译失败，错误代码: $LASTEXITCODE" -ForegroundColor Red
        Read-Host "按回车键退出"
        exit $LASTEXITCODE
    }
} catch {
    Write-Host "❌ 编译过程中发生错误: $_" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host ""
Write-Host "✅ 构建完成！" -ForegroundColor Green
Write-Host "📁 输出目录: $DIST_DIR" -ForegroundColor Cyan
Write-Host ""

Write-Host "📝 生成的文件：" -ForegroundColor Yellow
Get-ChildItem $DIST_DIR | Format-Table Name, Length, LastWriteTime

Write-Host ""
Write-Host "💡 使用说明：" -ForegroundColor Cyan
Write-Host "   - 运行程序: .\dist\clicker_win_$GOARCH.exe" -ForegroundColor White
Write-Host "   - 快捷键: Ctrl+F 开始，Ctrl+G 停止" -ForegroundColor White
Write-Host ""

Read-Host "按回车键退出"
