@echo off
echo Continuous Clicker - Simple Windows Build (with console)
echo =======================================================

:: Set basic environment
set GOOS=windows
set CGO_ENABLED=1

:: Clean and create dist directory
if exist "dist" rmdir /s /q "dist"
mkdir "dist"

:: Simple build command
echo Building Windows executable...
go build -ldflags="" -o "dist\clicker.exe" .

if %ERRORLEVEL% EQU 0 (
    echo Build successful!
    echo Output: dist\clicker.exe
) else (
    echo Build failed!
    pause
)

echo.
pause
