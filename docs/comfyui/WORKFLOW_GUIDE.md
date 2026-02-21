# ComfyUI Workflow 使用指南

## 概述

ComfyUI 集成支持两种方式：
1. **内置 Workflow**：使用预定义的 SVD (Stable Video Diffusion) workflow
2. **自定义 Workflow**：从 JSON 文件加载自定义 workflow

## 使用自定义 Workflow

### 1. 准备 Workflow JSON 文件

在 `workflows/` 目录下创建你的 workflow JSON 文件，例如 `video_wan2_2_14B_i2v.json`。

### 2. 使用占位符

在 workflow JSON 中，可以使用以下占位符，系统会自动替换：

- `%prompt%` - 用户输入的提示词
- `%image_url%` - 输入图片的 URL
- `%width%` - 视频宽度
- `%height%` - 视频高度
- `%fps%` - 帧率
- `%duration%` - 时长

### 3. Workflow 示例

```json
{
  "1": {
    "class_type": "LoadImage",
    "inputs": {
      "image": "%image_url%"
    }
  },
  "2": {
    "class_type": "CLIPTextEncode",
    "inputs": {
      "text": "%prompt%",
      "clip": ["4", 1]
    }
  },
  "3": {
    "class_type": "VideoModel",
    "inputs": {
      "model": ["4", 0],
      "conditioning": ["2", 0],
      "image": ["1", 0],
      "width": %width%,
      "height": %height%,
      "fps": %fps%,
      "frames": 40
    }
  },
  "4": {
    "class_type": "CheckpointLoaderSimple",
    "inputs": {
      "ckpt_name": "video_wan2_2_14B_i2v"
    }
  },
  "5": {
    "class_type": "SaveVideo",
    "inputs": {
      "video": ["3", 0],
      "filename_prefix": "comfyui_output"
    }
  }
}
```

### 4. 在 AI 配置中使用

1. 进入 **AI 服务配置** -> **视频生成**
2. 点击 **添加配置**
3. 选择厂商：**ComfyUI**
4. 模型选择：
   - 内置模型：`svd` 或 `svd_xt`
   - 自定义 workflow：输入 JSON 文件路径，如 `workflows/video_wan2_2_14B_i2v.json`
5. Base URL：`http://127.0.0.1:8188`（或你的 ComfyUI 服务器地址）
6. API Key：留空（ComfyUI 默认不需要）

### 5. 文件路径说明

- **相对路径**：相对于项目根目录，如 `workflows/video_wan2_2_14B_i2v.json`
- **绝对路径**：完整路径，如 `C:/ComfyUI/workflows/my_workflow.json`

## 工作流程

1. 用户在前端选择 ComfyUI 配置并生成视频
2. 系统检测模型名称是否以 `.json` 结尾
3. 如果是 JSON 文件：
   - 读取 workflow 文件
   - 替换所有占位符（`%prompt%`、`%image_url%` 等）
   - 发送到 ComfyUI API
4. 如果不是 JSON 文件：
   - 使用内置的 SVD workflow
   - 根据模型名称（svd/svd_xt）配置参数

## 注意事项

1. **文件位置**：建议将 workflow 文件放在 `workflows/` 目录下
2. **JSON 格式**：确保 JSON 格式正确，可以使用 JSON 验证工具检查
3. **占位符**：占位符区分大小写，必须完全匹配
4. **ComfyUI 节点**：workflow 中使用的节点必须在你的 ComfyUI 安装中可用
5. **模型文件**：确保 workflow 中引用的模型文件（如 `video_wan2_2_14B_i2v`）已下载到 ComfyUI 的 models 目录

## 调试

如果 workflow 执行失败，检查：

1. ComfyUI 服务器日志
2. 后端日志中的 `[ComfyUI]` 标记信息
3. workflow JSON 文件是否存在且格式正确
4. 占位符是否正确替换

## 示例配置

### 配置 1：使用内置 SVD
- 厂商：ComfyUI
- 模型：svd
- Base URL：http://127.0.0.1:8188

### 配置 2：使用自定义 Workflow
- 厂商：ComfyUI
- 模型：workflows/video_wan2_2_14B_i2v.json
- Base URL：http://127.0.0.1:8188

## 更多信息

参考文档：
- [ComfyUI 官方文档](https://github.com/comfyanonymous/ComfyUI)
- [COMFYUI_INTEGRATION.md](docs/COMFYUI_INTEGRATION.md)
- [COMFYUI_QUICKSTART.md](docs/COMFYUI_QUICKSTART.md)
