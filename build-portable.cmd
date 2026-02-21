@echo off
chcp 65001 >nul
echo ========================================
echo   火宝一键漫剧 - 便携版打包脚本
echo ========================================
echo.

REM 设置变量
set BUILD_DIR=huobao-drama-portable
set VERSION=v1.0.6-comfyui

REM 设置代理
set HTTP_PROXY=http://127.0.0.1:7890
set HTTPS_PROXY=http://127.0.0.1:7890

echo [1/6] 清理旧的构建目录...
if exist %BUILD_DIR% (
    rmdir /s /q %BUILD_DIR%
)
mkdir %BUILD_DIR%

echo.
echo [2/6] 构建前端...
cd web
call npm run build
if %errorlevel% neq 0 (
    echo 前端构建失败！
    cd ..
    pause
    exit /b 1
)
cd ..

echo.
echo [3/6] 编译后端...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-w -s" -o %BUILD_DIR%/huobao-drama.exe .
if %errorlevel% neq 0 (
    echo 后端编译失败！
    pause
    exit /b 1
)

echo.
echo [4/6] 复制必要文件...

REM 复制前端构建产物
xcopy /E /I /Y web\dist %BUILD_DIR%\web\dist

REM 复制配置文件
mkdir %BUILD_DIR%\configs
copy configs\config.example.yaml %BUILD_DIR%\configs\config.yaml

REM 复制数据库迁移文件
xcopy /E /I /Y migrations %BUILD_DIR%\migrations

REM 复制文档
copy README-CN.md %BUILD_DIR%\
copy LICENSE %BUILD_DIR%\

REM 创建数据目录
mkdir %BUILD_DIR%\data
mkdir %BUILD_DIR%\data\storage

REM 复制工作流目录
mkdir %BUILD_DIR%\workflows
copy workflows\.gitkeep %BUILD_DIR%\workflows\
copy workflows\README.md %BUILD_DIR%\workflows\

echo.
echo [5/6] 创建启动脚本...

REM 创建启动脚本
(
echo @echo off
echo chcp 65001 ^>nul
echo echo ========================================
echo echo   火宝一键漫剧 - 便携版
echo echo ========================================
echo echo.
echo echo 正在启动服务...
echo start "" huobao-drama.exe
echo echo.
echo echo 服务已启动！
echo echo 访问地址: http://localhost:5678
echo echo.
echo echo 按任意键退出此窗口（不影响服务运行）
echo pause ^>nul
) > %BUILD_DIR%\启动.cmd

REM 创建停止脚本
(
echo @echo off
echo chcp 65001 ^>nul
echo echo ========================================
echo echo   停止火宝一键漫剧服务
echo echo ========================================
echo echo.
echo taskkill /F /IM huobao-drama.exe ^>nul 2^>^&1
echo if %%errorlevel%% equ 0 (
echo     echo 服务已停止！
echo ^) else (
echo     echo 未找到运行中的服务
echo ^)
echo echo.
echo pause
) > %BUILD_DIR%\停止.cmd

REM 创建使用说明
(
echo # 火宝一键漫剧 - 便携版使用说明
echo.
echo ## 快速开始
echo.
echo 1. 双击 `启动.cmd` 启动服务
echo 2. 浏览器访问 http://localhost:5678
echo 3. 使用完毕后双击 `停止.cmd` 停止服务
echo.
echo ## 配置说明
echo.
echo - 配置文件: `configs/config.yaml`
echo - 数据目录: `data/`
echo - 工作流目录: `workflows/`
echo.
echo ## 注意事项
echo.
echo 1. 首次运行会自动创建数据库
echo 2. AI 服务需要在 Web 界面中配置 API Key
echo 3. 确保端口 5678 未被占用
echo.
echo ## 更多文档
echo.
echo 查看 README-CN.md 获取完整文档
echo.
) > %BUILD_DIR%\使用说明.md

echo.
echo [6/6] 打包压缩...
powershell -Command "Compress-Archive -Path '%BUILD_DIR%' -DestinationPath 'huobao-drama-%VERSION%-windows-x64.zip' -Force"

echo.
echo ========================================
echo   打包完成！
echo ========================================
echo.
echo 输出文件: huobao-drama-%VERSION%-windows-x64.zip
echo 大小: 
powershell -Command "(Get-Item 'huobao-drama-%VERSION%-windows-x64.zip').Length / 1MB | ForEach-Object { '{0:N2} MB' -f $_ }"
echo.
echo 便携版目录: %BUILD_DIR%\
echo.
echo 提示：
echo - 可以直接使用 %BUILD_DIR% 目录
echo - 或分发 .zip 压缩包
echo.
pause
