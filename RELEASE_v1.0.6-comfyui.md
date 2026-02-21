# Release v1.0.6-comfyui

## 🎉 ComfyUI 增强版发布

这是火宝一键漫剧的 ComfyUI 增强版本，集成了强大的 ComfyUI 工作流支持，让 AI 视频生成更加灵活和专业。

### 本 ComfyUI 贡献者

**你们喜爱的老王**

- 📺 **B站**: [https://space.bilibili.com/97727630](https://space.bilibili.com/97727630)
- 🎬 **专注**: ComfyUI 工作流、AI 视频生成

**你们喜爱的老王**

- 📺 **B站**: [https://space.bilibili.com/97727630](https://space.bilibili.com/97727630)
- 🎬 **专注**: ComfyUI 工作流、AI 视频生成

---

## ✨ 主要特性

### ComfyUI 集成
- ✅ 完整的 ComfyUI API 集成
- ✅ 支持自定义工作流导入
- ✅ 图片生成和视频生成工作流
- ✅ 实时进度监控和预览

### 工作流管理
- ✅ 可视化工作流选择器
- ✅ 支持多种预设工作流
- ✅ 工作流参数配置
- ✅ 批量生成支持

### AI 服务支持
- ✅ OpenAI / Azure OpenAI
- ✅ 通义千问 / 智谱 AI
- ✅ DeepSeek / Moonshot
- ✅ 本地 Ollama 模型

### 用户体验
- ✅ 现代化 UI 设计
- ✅ 深色/浅色主题切换
- ✅ 中英文双语支持
- ✅ 响应式布局

---

## 📦 下载

### Windows 便携版 (推荐)

**文件**: `huobao-drama-v1.0.6-comfyui-windows-x64.zip` (约 20 MB)

**特点**:
- 开箱即用，无需安装
- 包含完整的前后端
- 双击启动，简单方便

**使用方法**:
1. 解压 zip 文件
2. 双击 `启动.cmd`
3. 浏览器访问 http://localhost:5678
4. 使用完毕后双击 `停止.cmd`

### 系统要求

- Windows 10/11 (64位)
- 4GB+ 内存
- 需要配置 ComfyUI 服务地址
- 需要配置 AI 服务 API Key

---

## 🚀 快速开始

### 1. 下载并解压

下载 `huobao-drama-v1.0.6-comfyui-windows-x64.zip` 并解压到任意目录

### 2. 配置 ComfyUI

编辑 `configs/config.yaml`，配置你的 ComfyUI 服务地址：

```yaml
comfyui:
  base_url: "http://127.0.0.1:8188"  # 你的 ComfyUI 地址
  timeout: 300
```

### 3. 启动服务

双击 `启动.cmd`，等待服务启动完成

### 4. 配置 AI 服务

在 Web 界面点击右上角的 "AI配置" 按钮，配置你的 AI 服务：
- 选择 AI 服务提供商
- 输入 API Key
- 选择模型

### 5. 开始创作

- 创建新项目
- 输入剧本或让 AI 生成
- 选择工作流生成图片/视频
- 使用时间轴编辑器制作最终视频

---

## 📚 文档

- [快速入门指南](docs/COMFYUI_QUICKSTART.md)
- [ComfyUI 集成说明](docs/COMFYUI_INTEGRATION.md)
- [ComfyUI 配置示例](docs/comfyui-config-example.yaml)
- [完整文档索引](docs/README.md)

---

## 🔧 从源码构建

如果你想从源码构建，请参考 [README-CN.md](README-CN.md) 中的开发指南。

---

## 📝 更新日志

### v1.0.6-comfyui (2026-02-21)

**新增功能**:
- 在 UI 标题栏添加 ComfyUI 贡献者信息
- 添加 B 站链接快速访问
- 优化便携版打包流程
- 添加代理支持

**改进**:
- 更新 README 作者信息
- 完善文档结构
- 优化启动脚本

**修复**:
- 修复工作流文件排除规则
- 修复生成文件的 .gitignore 配置

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

如果你有好的 ComfyUI 工作流，也欢迎分享到 `workflows/` 目录。

---

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

---

## 💬 联系方式

- **GitHub**: [https://github.com/ops120/huobao-drama-comfyui](https://github.com/ops120/huobao-drama-comfyui)
- **B站**: [https://space.bilibili.com/97727630](https://space.bilibili.com/97727630)

---

**感谢使用火宝一键漫剧 ComfyUI 版！** 🎬✨
