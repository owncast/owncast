@echo off
setlocal enabledelayedexpansion

REM This script will make your local Owncast server available at a public URL.
REM It's particularly useful for testing on mobile devices or want to test
REM activitypub integration.
REM Pass a custom domain as an argument if you have previously set it up at
REM localhost.run. Otherwise, a random hostname will be generated.
REM SET DOMAIN=me.example.com && test\test-local.bat
REM Pass a port number as an argument if you are running Owncast on a different port.
REM By default, it will use port 8080.
REM SET PORT=8080 && test\test-local.bat

REM Set default port if not provided
if "%PORT%"=="" set PORT=8080
set HOST=localhost

echo Checking if web server is running on port %PORT%...

REM Using PowerShell to check if the port is open (equivalent to nc -zv)
powershell -Command "try { $tcpConnection = New-Object System.Net.Sockets.TcpClient; $tcpConnection.Connect('%HOST%', %PORT%); $tcpConnection.Close(); exit 0 } catch { exit 1 }"

if %ERRORLEVEL% equ 0 (
    echo Your web server is running on port %PORT%. Good.
) else (
    echo Please make sure your Owncast server is running on port %PORT%.
    exit /b 1
)

if not "%DOMAIN%"=="" (
    echo Attempting to use custom domain: %DOMAIN%
    ssh -R "%DOMAIN%":80:localhost:%PORT% localhost.run
) else (
    echo Using auto-generated hostname for tunnel.
    ssh -R 80:localhost:%PORT% localhost.run
)