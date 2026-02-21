@echo off
chcp 65001 >nul
echo ========================================
echo   火宝一键漫剧 - 生产环境启动脚本
echo ========================================
echo.

REM 设置 Go 代理
set GOPROXY=https://goproxy.cn,direct

REM 检查是否已有进程在运行
echo [1/4] 检查现有进程...
for /f "tokens=2" %%i in ('netstat -ano ^| findstr :5678 ^| findstr LISTENING') do (
    echo 发现端口 5678 被占用，正在停止进程 %%i...
    taskkill /F /PID %%i >nul 2>&1
)

echo.
echo [2/4] 构建前端...
cd web
call npm run build
if %errorlevel% neq 0 (
    echo 前端构建失败！
    cd ..
    pause
    exit /b 1
)
cd ..
echo 前端构建完成！

echo.
echo [3/4] 启动后端服务 (端口 5678)...
echo 后端将同时提供 API 和前端静态文件服务
start "火宝生产服务" cmd /k "echo 服务启动中... && go run main.go"

REM 等待服务启动
timeout /t 3 /nobreak >nul

echo.
echo [4/4] 启动完成！
echo.
echo ========================================
echo   服务访问地址：
echo   - 应用首页:  http://localhost:5678
echo   - API 服务:  http://localhost:5678/api/v1
echo   - 健康检查:  http://localhost:5678/health
echo ========================================
echo.
echo 提示：
echo   - 服务窗口会自动打开
echo   - 关闭窗口即可停止服务
echo   - 运行 stop-dev.cmd 可停止所有服务
echo   - 按任意键退出此窗口（不影响服务运行）
echo.
pause
