@echo off
chcp 65001 >nul
echo ========================================
echo   火宝一键漫剧 - 开发环境启动脚本
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

for /f "tokens=2" %%i in ('netstat -ano ^| findstr :3012 ^| findstr LISTENING') do (
    echo 发现端口 3012 被占用，正在停止进程 %%i...
    taskkill /F /PID %%i >nul 2>&1
)

echo.
echo [2/4] 启动后端服务 (端口 5678)...
start "火宝后端服务" cmd /k "echo 后端服务启动中... && go run main.go"

REM 等待后端启动
timeout /t 3 /nobreak >nul

echo.
echo [3/4] 启动前端开发服务器 (端口 3012)...
cd web
start "火宝前端服务" cmd /k "echo 前端服务启动中... && npm run dev"
cd ..

echo.
echo [4/4] 启动完成！
echo.
echo ========================================
echo   服务访问地址：
echo   - 前端开发服务器: http://localhost:3012
echo   - 后端 API 服务:  http://localhost:5678
echo   - 健康检查:       http://localhost:5678/health
echo ========================================
echo.
echo 提示：
echo   - 两个服务窗口会自动打开
echo   - 关闭窗口即可停止对应服务
echo   - 按任意键退出此窗口（不影响服务运行）
echo.
pause
