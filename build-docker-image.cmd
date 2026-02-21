@echo off
chcp 65001 >nul
echo ========================================
echo   火宝一键漫剧 - Docker 镜像打包
echo ========================================
echo.

set VERSION=v1.0.6-comfyui
set IMAGE_NAME=huobao-drama

echo [1/3] 构建 Docker 镜像...
docker build -t %IMAGE_NAME%:%VERSION% .
if %errorlevel% neq 0 (
    echo Docker 镜像构建失败！
    pause
    exit /b 1
)

echo.
echo [2/3] 标记为 latest...
docker tag %IMAGE_NAME%:%VERSION% %IMAGE_NAME%:latest

echo.
echo [3/3] 导出镜像文件...
docker save -o huobao-drama-%VERSION%-docker.tar %IMAGE_NAME%:%VERSION%

echo.
echo ========================================
echo   Docker 镜像打包完成！
echo ========================================
echo.
echo 镜像文件: huobao-drama-%VERSION%-docker.tar
echo 大小: 
powershell -Command "(Get-Item 'huobao-drama-%VERSION%-docker.tar').Length / 1MB | ForEach-Object { '{0:N2} MB' -f $_ }"
echo.
echo 使用方法：
echo 1. 导入镜像: docker load -i huobao-drama-%VERSION%-docker.tar
echo 2. 运行容器: docker run -d -p 5678:5678 %IMAGE_NAME%:%VERSION%
echo.
pause
