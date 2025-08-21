@echo off
chcp 65001 >nul
echo Continuous Clicker - Windows Build Script
echo ==========================================

:: Project root directory
set ROOT_DIR=%~dp0
set DIST_DIR=%ROOT_DIR%dist
set FONT_FILE=%ROOT_DIR%SourceHanSansSC-Bold.otf

:: Clean output directory
echo Cleaning old dist directory...
if exist "%DIST_DIR%" rmdir /s /q "%DIST_DIR%"
mkdir "%DIST_DIR%"

:: Check font file
if exist "%FONT_FILE%" (
    echo Font file found: %FONT_FILE%
) else (
    echo Font file not found: %FONT_FILE%
)

:: Detect current architecture
echo Detecting system architecture...
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    echo Detected AMD64 architecture
    set GOARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    echo Detected ARM64 architecture
    set GOARCH=arm64
) else (
    echo Detected %PROCESSOR_ARCHITECTURE% architecture
    set GOARCH=amd64
)

:: Set environment variables - ensure no console window
echo Setting build environment...
set GOOS=windows
set CGO_ENABLED=1

:: Build project - use console hiding parameters
echo Starting Windows build...
echo Target system: Windows (%GOARCH%)
echo Using console hiding parameters...

go build -ldflags="-H windowsgui -s -w" -o "%DIST_DIR%\clicker_win_%GOARCH%.exe" .

if %ERRORLEVEL% EQU 0 (
    echo Windows build completed successfully
) else (
    echo Build failed with error code: %ERRORLEVEL%
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo Build completed!
echo Output directory: %DIST_DIR%
echo.
echo Generated files:
dir "%DIST_DIR%"

echo.
echo Usage instructions:
echo    - Run program: .\dist\clicker_win_%GOARCH%.exe
echo    - Hotkeys: Ctrl+F to start, Ctrl+G to stop
echo    - Program will run in pure GUI mode, no console window
echo.
pause
