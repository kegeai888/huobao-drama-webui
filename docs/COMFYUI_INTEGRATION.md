# ComfyUI 视频生成集成指南

## 概述

本项目已集成 ComfyUI 作为视频生成的一个选项，支持通过 ComfyUI 的 API 进行图生视频操作。

## 功能特性

- 支持 ComfyUI 标准 API 协议
- 自动构建图生视频工作流
- 支持自定义模型和参数
- 异步任务状态查询
- 支持队列管理

## 配置步骤

### 1. 启动 ComfyUI 服务

确保你的 ComfyUI 服务已启动并开启 API 访问：

```bash
# 启动 ComfyUI 并监听所有接口
python main.py --listen 0.0.0.0 --port 8188
```

### 2. 在系统中配置 ComfyUI

访问系统设置 → AI 服务配置 → 视频生成标签页，添加新配置：

**配置参数：**
- **名称**: ComfyUI-本地 (或自定义名称)
- **厂商**: ComfyUI
- **Base URL**: `http://localhost:8188` (或你的 ComfyUI 服务地址)
- **API Key**: 留空 (ComfyUI 默认不需要)
- **模型**: 
  - `svd` - Stable Video Diffusion
  - `svd_xt` - Stable Video Diffusion XT (更长时长)
  - `custom` - 自定义模型名称
- **优先级**: 0-100 (数值越大优先级越高)

### 3. 配置示例

```yaml
# 本地 ComfyUI
Base URL: http://localhost:8188
API Key: (留空)
Model: svd

# 远程 ComfyUI (如果有认证)
Base URL: https://your-comfyui-server.com
API Key: your-api-key
Model: svd_xt
```

## API 端点

ComfyUI 客户端使用以下端点：

- **提交任务**: `POST /prompt`
- **查询队列**: `GET /queue`
- **查询历史**: `GET /history/{prompt_id}`
- **获取文件**: `GET /view?filename={name}&subfolder={folder}&type={type}`

## 工作流说明

### 默认工作流

系统内置了一个基础的图生视频工作流，包含以下节点：

1. **LoadImage** - 加载输入图片
2. **CLIPTextEncode** - 编码文本提示词
3. **CheckpointLoaderSimple** - 加载模型检查点
4. **VideoLinearCFGGuidance** - 视频 CFG 引导
5. **SVD_img2vid_Conditioning** - SVD 图生视频条件
6. **KSampler** - K 采样器
7. **VAEDecode** - VAE 解码
8. **VHS_VideoCombine** - 视频合成

### 自定义工作流

如果需要自定义工作流，可以修改 `pkg/video/comfyui_client.go` 中的 `buildWorkflow` 方法。

## 参数配置

### 视频生成参数

通过 `VideoOptions` 可以配置以下参数：

```go
options := &VideoOptions{
    Duration:    5,          // 视频时长(秒)
    FPS:         8,          // 帧率
    Resolution:  "512x512",  // 分辨率
    AspectRatio: "16:9",     // 宽高比
}
```

### 工作流参数

在 `buildWorkflow` 中可以调整：

- `video_frames`: 视频帧数 (Duration * FPS)
- `motion_bucket_id`: 运动强度 (0-255)
- `steps`: 采样步数
- `cfg`: CFG 强度
- `sampler_name`: 采样器类型 (euler, euler_a, dpm++, etc.)
- `scheduler`: 调度器 (karras, normal, etc.)

## 使用示例

### 通过 API 调用

```bash
# 创建视频生成任务
curl -X POST http://localhost:5678/api/v1/video-generations \
  -H "Content-Type: application/json" \
  -d '{
    "storyboard_id": 123,
    "image_url": "http://example.com/image.jpg",
    "prompt": "A beautiful sunset scene",
    "model": "svd",
    "duration": 5
  }'
```

### 在前端使用

1. 进入剧本编辑页面
2. 选择分镜
3. 点击"生成视频"
4. 在模型选择中选择 ComfyUI 配置的模型
5. 系统会自动使用 ComfyUI 生成视频

## 故障排查

### 常见问题

**1. 连接失败**
```
Error: failed to send request: dial tcp: connection refused
```
解决方案：
- 检查 ComfyUI 是否正在运行
- 确认 Base URL 配置正确
- 检查防火墙设置

**2. 工作流错误**
```
Error: workflow error: node_errors
```
解决方案：
- 检查模型文件是否存在于 ComfyUI 的 models 目录
- 确认工作流节点配置正确
- 查看 ComfyUI 控制台日志

**3. 任务超时**
```
Error: task timeout
```
解决方案：
- 增加客户端超时时间
- 检查 ComfyUI 服务器性能
- 减少视频帧数或降低分辨率

### 调试模式

启用调试日志查看详细信息：

```go
// 在 config.yaml 中设置
app:
  debug: true
```

查看日志输出：
```
[ComfyUI] Sending request to: http://localhost:8188/prompt
[ComfyUI] Task created - PromptID: abc123
[ComfyUI] Task status - ID: abc123, Status: processing, Completed: false
```

## 性能优化

### 建议配置

**本地部署 (推荐配置):**
- GPU: NVIDIA RTX 3060 或更高
- VRAM: 8GB 或更高
- RAM: 16GB 或更高

**参数优化:**
- 降低分辨率: `512x512` → `256x256`
- 减少帧数: `40 frames` → `25 frames`
- 降低采样步数: `20 steps` → `15 steps`

### 批量处理

对于大量视频生成任务，建议：
1. 使用队列管理避免并发过多
2. 设置合理的轮询间隔 (5-10秒)
3. 监控 GPU 使用率

## 扩展开发

### 添加自定义节点

如果你的 ComfyUI 安装了自定义节点，可以修改工作流：

```go
func (c *ComfyUIClient) buildWorkflow(imageURL, prompt string, options *VideoOptions) map[string]interface{} {
    workflow := map[string]interface{}{
        // 添加你的自定义节点
        "custom_node": map[string]interface{}{
            "class_type": "YourCustomNode",
            "inputs": map[string]interface{}{
                "param1": "value1",
            },
        },
    }
    return workflow
}
```

### 支持更多模型

在前端配置中添加新模型：

```typescript
// web/src/views/settings/AIConfig.vue
{ id: "comfyui", name: "ComfyUI", models: [
    "svd", 
    "svd_xt", 
    "animatediff",  // 新增
    "custom"
]},
```

## 参考资源

- [ComfyUI 官方文档](https://github.com/comfyanonymous/ComfyUI)
- [ComfyUI API 文档](https://github.com/comfyanonymous/ComfyUI/wiki/API)
- [Stable Video Diffusion](https://stability.ai/stable-video)

## 更新日志

### v1.0.6 (2026-02-20)
- ✨ 新增 ComfyUI 视频生成支持
- 🔧 支持自定义工作流配置
- 📝 完善文档和使用指南
