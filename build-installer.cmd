@echo off
chcp 65001 >nul
echo ========================================
echo   火宝一键漫剧 - 安装包打包脚本
echo ========================================
echo.
echo 此脚本需要安装 Inno Setup 或 NSIS
echo.
echo 推荐使用便携版打包脚本: build-portable.cmd
echo.
echo 如需创建安装包，请：
echo 1. 安装 Inno Setup: https://jrsoftware.org/isdl.php
echo 2. 创建 installer.iss 配置文件
echo 3. 运行 Inno Setup 编译
echo.
echo 或使用 NSIS:
echo 1. 安装 NSIS: https://nsis.sourceforge.io/
echo 2. 创建 installer.nsi 配置文件
echo 3. 运行 makensis 编译
echo.
pause
