# ComfyUI 快速开始指南

## 5 分钟快速配置

### 步骤 1: 安装 ComfyUI

```bash
# 克隆 ComfyUI
git clone https://github.com/comfyanonymous/ComfyUI.git
cd ComfyUI

# 安装依赖
pip install -r requirements.txt

# 下载 SVD 模型 (约 10GB)
# 将模型放到 models/checkpoints/ 目录
```

### 步骤 2: 启动 ComfyUI

```bash
# 启动 ComfyUI 服务
python main.py --listen 0.0.0.0 --port 8188
```

访问 http://localhost:8188 确认服务正常运行。

### 步骤 3: 配置火宝短剧

1. 启动火宝短剧系统
2. 访问 **设置 → AI 服务配置 → 视频生成**
3. 点击 **添加配置**
4. 填写以下信息：

```
名称: ComfyUI-本地
厂商: ComfyUI
Base URL: http://localhost:8188
API Key: (留空)
模型: svd
优先级: 50
```

5. 点击 **保存**

### 步骤 4: 测试生成

1. 进入任意剧本的分镜编辑页面
2. 选择一个已有图片的分镜
3. 点击 **生成视频**
4. 在模型选择中选择 `svd`
5. 点击 **开始生成**

等待 1-3 分钟，视频生成完成！

## Docker 快速部署

### 使用 Docker Compose

创建 `docker-compose.comfyui.yml`:

```yaml
version: '3.8'

services:
  comfyui:
    image: yanwk/comfyui-boot:latest
    ports:
      - "8188:8188"
    volumes:
      - ./comfyui/models:/opt/ComfyUI/models
      - ./comfyui/output:/opt/ComfyUI/output
    environment:
      - CLI_ARGS=--listen 0.0.0.0
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

启动：

```bash
docker-compose -f docker-compose.comfyui.yml up -d
```

## 常见问题

### Q: 模型文件在哪里下载？

A: 访问 [Hugging Face](https://huggingface.co/stabilityai/stable-video-diffusion-img2vid-xt) 下载 SVD 模型。

### Q: 生成速度很慢怎么办？

A: 
- 确保使用 GPU 运行
- 降低分辨率和帧数
- 减少采样步数

### Q: 提示 "connection refused" 错误？

A: 
- 检查 ComfyUI 是否正在运行
- 确认端口 8188 未被占用
- 检查防火墙设置

### Q: 如何使用自定义模型？

A: 
1. 将模型放到 ComfyUI 的 `models/checkpoints/` 目录
2. 在配置中添加模型名称
3. 重启 ComfyUI 服务

## 下一步

- 阅读 [完整集成指南](COMFYUI_INTEGRATION.md)
- 查看 [配置示例](comfyui-config-example.yaml)
- 了解 [工作流自定义](COMFYUI_INTEGRATION.md#自定义工作流)

## 获取帮助

- 查看 [故障排查](COMFYUI_INTEGRATION.md#故障排查)
- 提交 [Issue](https://github.com/chatfire-AI/huobao-drama/issues)
- 加入项目交流群
