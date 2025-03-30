@echo off
setlocal

REM Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Please run this script as Administrator
    pause
    exit /b 1
)

REM Set paths
set "BUILD_DIR=.\build"
set "EXE_NAME=gowin.exe"

REM Build the Go application
echo Building application...
go build -o "%BUILD_DIR%\%EXE_NAME%" .\cmd\main.go
if %errorLevel% neq 0 (
    echo Failed to build application
    pause
    exit /b 1
)

@REM copy assets folder to installation directory
xcopy /E /I assets "%BUILD_DIR%/assets"


echo Check build directory for executable file 
pause