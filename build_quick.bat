@echo off
echo 快速编译 Windows 版本...
set GOOS=windows
set CGO_ENABLED=1
go build -o clicker.exe .
if %ERRORLEVEL% EQU 0 (
    echo 编译成功！生成 clicker.exe
) else (
    echo 编译失败！
)
pause
