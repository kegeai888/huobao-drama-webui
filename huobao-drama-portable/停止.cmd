@echo off
chcp 65001 >nul
echo ========================================
echo   停止火宝一键漫剧服务
echo ========================================
echo.
taskkill /F /IM huobao-drama.exe >nul 2>&1
if %errorlevel% equ 0 (
    echo 服务已停止！
) else (
    echo 未找到运行中的服务
)
echo.
pause
