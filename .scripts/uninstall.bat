@REM run script from root file as administrator

@echo off
setlocal

REM Check if go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
	echo Go is not installed. Please install Go and try again.
	exit /b 1
)

REM Build the service
go build -o build/gowin-service.exe ./cmd/service
if %errorlevel% neq 0 (
	echo Failed to build the service.
	exit /b 1
)

REM Stop the service
build\gowin-service.exe stop
if %errorlevel% neq 0 (
	echo Failed to stop the service.
)

REM Uninstall the service
build\gowin-service.exe uninstall
if %errorlevel% neq 0 (
	echo Failed to uninstall the service.
)

endlocal

