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
set "INSTALL_DIR=%ProgramFiles%\GoWin"
set "STARTUP_DIR=%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup"
set "EXE_NAME=gowin.exe"

REM Create installation directory
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Build the Go application
echo Building application...
go build -o "%INSTALL_DIR%\%EXE_NAME%" .\cmd\main.go
if %errorLevel% neq 0 (
    echo Failed to build application
    pause
    exit /b 1
)

@REM copy assets folder to installation directory
xcopy /E /I assets "%INSTALL_DIR%/assets"


REM Create startup script
echo @echo off > "%STARTUP_DIR%\start_gowin.bat"
echo start "" "%INSTALL_DIR%\%EXE_NAME%" >> "%STARTUP_DIR%\start_gowin.bat"

echo Installation completed successfully!
echo The application will start automatically when you log in.
pause