# Continuous Clicker - PowerShell Build Script
Write-Host "Continuous Clicker - PowerShell Build" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# Set environment variables
$env:GOOS = "windows"
$env:CGO_ENABLED = "1"

# Clean and create dist directory
if (Test-Path "dist") {
    Remove-Item "dist" -Recurse -Force
}
New-Item -ItemType Directory -Name "dist" | Out-Null

# Build command
Write-Host "Building Windows executable..." -ForegroundColor Yellow
try {
    go build -ldflags="-H windowsgui -s -w" -o "dist\clicker.exe" .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build successful!" -ForegroundColor Green
        Write-Host "Output: dist\clicker.exe" -ForegroundColor Green
        
        # Show file info
        $file = Get-Item "dist\clicker.exe"
        Write-Host "File size: $([math]::Round($file.Length / 1MB, 2)) MB" -ForegroundColor Cyan
        Write-Host "Created: $($file.CreationTime)" -ForegroundColor Cyan
    } else {
        Write-Host "Build failed with exit code: $LASTEXITCODE" -ForegroundColor Red
    }
} catch {
    Write-Host "Build error: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "Usage instructions:" -ForegroundColor Yellow
Write-Host "  - Run program: .\dist\clicker.exe" -ForegroundColor White
Write-Host "  - Hotkeys: Ctrl+F to start, Ctrl+G to stop" -ForegroundColor White
Write-Host "  - Program runs in pure GUI mode, no console window" -ForegroundColor White
Write-Host ""
Read-Host "Press Enter to continue"