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

REM Install the service
build\gowin-service.exe install
if %errorlevel% neq 0 (
	echo Failed to install the service.
	exit /b 1
)

REM Start the service
build\gowin-service.exe start
if %errorlevel% neq 0 (
	echo Failed to start the service.
	exit /b 1
)

sc failure "Gowin" reset= 30000 actions= restart/5000/restart/10000/restart/20000

endlocal

