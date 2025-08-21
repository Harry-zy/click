@echo off
chcp 65001 >nul
echo 🚀 连续点击器 - Windows 快速构建脚本
echo ================================

:: 项目根目录
set ROOT_DIR=%~dp0
set DIST_DIR=%ROOT_DIR%dist
set FONT_FILE=%ROOT_DIR%SourceHanSansSC-Bold.otf

:: 输出清理
echo 🧹 清理旧的 dist 目录...
if exist "%DIST_DIR%" rmdir /s /q "%DIST_DIR%"
mkdir "%DIST_DIR%"

:: 检查字体文件
if exist "%FONT_FILE%" (
    echo ✅ 字体文件存在: %FONT_FILE%
) else (
    echo ⚠️ 字体文件未找到: %FONT_FILE%
)

:: 检测当前架构
echo 🔍 检测系统架构...
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    echo 🖥️ 检测到 AMD64 架构
    set GOARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    echo 🖥️ 检测到 ARM64 架构
    set GOARCH=arm64
) else (
    echo 🖥️ 检测到 %PROCESSOR_ARCHITECTURE% 架构
    set GOARCH=amd64
)

:: 设置环境变量
echo 🔧 设置编译环境...
set GOOS=windows
set CGO_ENABLED=1

:: 编译项目
echo 🚀 开始编译 Windows 版本...
echo 📱 目标系统: Windows (%GOARCH%)

go build -ldflags="-H windowsgui" -o "%DIST_DIR%\clicker_win_%GOARCH%.exe" .

if %ERRORLEVEL% EQU 0 (
    echo ✅ Windows 版本编译完成
) else (
    echo ❌ 编译失败，错误代码: %ERRORLEVEL%
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo ✅ 构建完成！
echo 📁 输出目录: %DIST_DIR%
echo.
echo 📝 生成的文件：
dir "%DIST_DIR%"

echo.
echo 💡 使用说明：
echo    - 运行程序: .\dist\clicker_win_%GOARCH%.exe
echo    - 快捷键: Ctrl+F 开始，Ctrl+G 停止
echo.
pause
