# ComfyUI 视频生成集成

## 🎉 新功能

火宝一键漫剧现已支持 ComfyUI 作为视频生成服务！

### ✨ 主要特性

- 🏠 **本地部署** - 完全控制你的数据和模型
- 🎨 **自定义工作流** - 灵活配置视频生成流程
- 🚀 **高性能** - 利用本地 GPU 加速
- 🔧 **易于集成** - 简单配置即可使用

## 📚 文档导航

### 快速开始
- [5分钟快速配置](docs/COMFYUI_QUICKSTART.md) - 最快上手指南
- [配置示例](docs/comfyui-config-example.yaml) - 配置参考

### 详细文档
- [完整集成指南](docs/COMFYUI_INTEGRATION.md) - 深入了解所有功能
- [技术总结](docs/COMFYUI_SUMMARY.md) - 技术实现细节
- [变更说明](COMFYUI_CHANGES.md) - 代码变更详情

### 示例代码
- [使用示例脚本](docs/comfyui-usage-example.sh) - API 调用示例

## 🚀 快速开始

### 1. 启动 ComfyUI

```bash
# 克隆并启动 ComfyUI
git clone https://github.com/comfyanonymous/ComfyUI.git
cd ComfyUI
pip install -r requirements.txt
python main.py --listen 0.0.0.0 --port 8188
```

### 2. 配置系统

访问 **设置 → AI 服务配置 → 视频生成**，添加：

```
名称: ComfyUI-本地
厂商: ComfyUI
Base URL: http://localhost:8188
模型: svd
```

### 3. 开始生成

在分镜编辑页面选择 ComfyUI 模型，点击生成视频！

## 📋 支持的模型

- **svd** - Stable Video Diffusion (推荐)
- **svd_xt** - Stable Video Diffusion XT (更长时长)
- **custom** - 自定义模型

## 🔧 系统要求

### 最低配置
- GPU: NVIDIA RTX 3060 (8GB VRAM)
- RAM: 16GB
- 存储: 20GB (含模型)

### 推荐配置
- GPU: NVIDIA RTX 3080+ (10GB+ VRAM)
- RAM: 32GB
- 存储: 50GB

## 📊 性能参考

| 配置 | 分辨率 | 帧数 | 生成时间 |
|------|--------|------|----------|
| RTX 3060 | 512x512 | 25 | ~60秒 |
| RTX 3080 | 512x512 | 40 | ~60秒 |
| RTX 4090 | 1024x1024 | 40 | ~45秒 |

## 🐛 故障排查

### 连接失败
```bash
# 检查 ComfyUI 是否运行
curl http://localhost:8188/queue

# 检查端口占用
netstat -an | grep 8188
```

### 生成失败
1. 确认模型文件已下载
2. 检查 GPU 内存是否充足
3. 查看 ComfyUI 控制台日志

### 性能问题
- 降低分辨率 (512x512)
- 减少帧数 (25 frames)
- 调整采样步数 (15 steps)

## 🤝 获取帮助

- 📖 查看 [完整文档](docs/COMFYUI_INTEGRATION.md)
- 🐛 提交 [Issue](https://github.com/chatfire-AI/huobao-drama/issues)
- 💬 加入项目交流群

## 📝 更新日志

### v1.0.6 (2026-02-20)
- ✨ 新增 ComfyUI 视频生成支持
- 📝 完善文档和使用指南
- 🧪 添加单元测试

## 📄 许可证

遵循项目主许可证 (CC BY-NC-SA 4.0)

---

**开始使用 ComfyUI，让视频生成更自由！** 🎬
