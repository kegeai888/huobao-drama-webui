@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion
echo ========================================
echo   火宝一键漫剧 - 停止开发服务
echo ========================================
echo.

REM 停止后端服务 (端口 5678)
echo [1/3] 检查后端服务端口 5678...
set BACKEND_STOPPED=0
for /f "tokens=5" %%i in ('netstat -ano ^| findstr :5678 ^| findstr LISTENING') do (
    echo 发现进程 %%i 占用端口 5678
    echo 正在停止进程 %%i...
    taskkill /F /PID %%i >nul 2>&1
    if !errorlevel! equ 0 (
        echo [成功] 后端服务进程 %%i 已停止
        set BACKEND_STOPPED=1
    ) else (
        echo [失败] 无法停止进程 %%i
    )
)

if !BACKEND_STOPPED! equ 0 (
    echo [信息] 端口 5678 未被占用，无需停止
)

REM 等待端口释放
timeout /t 1 /nobreak >nul

REM 再次检查端口 5678
for /f "tokens=5" %%i in ('netstat -ano ^| findstr :5678 ^| findstr LISTENING') do (
    echo [警告] 端口 5678 仍被进程 %%i 占用，尝试强制停止...
    taskkill /F /PID %%i >nul 2>&1
)

echo.
echo [2/3] 检查前端服务端口 3012...
set FRONTEND_STOPPED=0
for /f "tokens=5" %%i in ('netstat -ano ^| findstr :3012 ^| findstr LISTENING') do (
    echo 发现进程 %%i 占用端口 3012
    echo 正在停止进程 %%i...
    taskkill /F /PID %%i >nul 2>&1
    if !errorlevel! equ 0 (
        echo [成功] 前端服务进程 %%i 已停止
        set FRONTEND_STOPPED=1
    ) else (
        echo [失败] 无法停止进程 %%i
    )
)

if !FRONTEND_STOPPED! equ 0 (
    echo [信息] 端口 3012 未被占用，无需停止
)

REM 等待端口释放
timeout /t 1 /nobreak >nul

REM 再次检查端口 3012
for /f "tokens=5" %%i in ('netstat -ano ^| findstr :3012 ^| findstr LISTENING') do (
    echo [警告] 端口 3012 仍被进程 %%i 占用，尝试强制停止...
    taskkill /F /PID %%i >nul 2>&1
)

echo.
echo [3/3] 验证端口状态...
set PORT_5678_FREE=1
set PORT_3012_FREE=1

for /f "tokens=5" %%i in ('netstat -ano ^| findstr :5678 ^| findstr LISTENING') do (
    echo [错误] 端口 5678 仍被进程 %%i 占用！
    set PORT_5678_FREE=0
)

for /f "tokens=5" %%i in ('netstat -ano ^| findstr :3012 ^| findstr LISTENING') do (
    echo [错误] 端口 3012 仍被进程 %%i 占用！
    set PORT_3012_FREE=0
)

if !PORT_5678_FREE! equ 1 (
    echo [✓] 端口 5678 已释放
)

if !PORT_3012_FREE! equ 1 (
    echo [✓] 端口 3012 已释放
)

echo.
echo ========================================
if !PORT_5678_FREE! equ 1 if !PORT_3012_FREE! equ 1 (
    echo   所有服务已成功停止，端口已释放
) else (
    echo   部分端口仍被占用，请手动检查
    echo   使用命令: netstat -ano ^| findstr :5678
    echo   使用命令: netstat -ano ^| findstr :3012
)
echo ========================================
echo.
pause
