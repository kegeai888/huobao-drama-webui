# ComfyUI Workflow 文件管理指南

## 功能概述

系统支持上传和管理 ComfyUI workflow JSON 文件，让你可以使用自定义的图片和视频生成流程。

## 使用方法

### 1. 上传 Workflow 文件

#### 方式一：通过 AI 配置界面上传

1. 进入 **AI 服务配置**
2. 选择 **图片生成** 或 **视频生成** 标签
3. 点击 **添加配置**
4. 厂商选择：**ComfyUI**
5. 在模型字段，点击 **上传 Workflow** 按钮
6. 选择你的 `.json` 文件（最大 10MB）
7. 上传成功后，文件会自动出现在下拉列表中

#### 方式二：直接放入 workflows 目录

将 `.json` 文件直接复制到项目根目录的 `workflows/` 文件夹中，然后在配置界面点击"刷新列表"按钮。

### 2. 使用 Workflow

1. 在模型下拉列表中选择：
   - **预设模型**：`sdxl`、`sd15`、`flux`（图片）或 `svd`、`svd_xt`（视频）
   - **自定义 Workflow**：选择你上传的 JSON 文件

2. 配置完成后点击"创建"

3. 在生成图片/视频时，系统会自动：
   - 加载选择的 workflow 文件
   - 替换其中的占位符
   - 发送到 ComfyUI 执行

### 3. 管理 Workflow 文件

- **查看列表**：在模型下拉列表中查看所有已上传的 workflow
- **刷新列表**：点击"刷新列表"按钮重新加载
- **删除文件**：在下拉列表中点击文件右侧的"删除"按钮

## Workflow 文件格式

### 支持的占位符

在 workflow JSON 文件中，可以使用以下占位符，系统会自动替换：

#### 图片生成占位符
- `%prompt%` - 用户输入的提示词
- `%negative_prompt%` - 负面提示词
- `%width%` - 图片宽度（默认 512）
- `%height%` - 图片高度（默认 512）
- `%steps%` - 采样步数（默认 20）
- `%cfg_scale%` - CFG 比例（默认 7.0）
- `%seed%` - 随机种子

#### 视频生成占位符
- `%prompt%` - 用户输入的提示词
- `%image_url%` - 输入图片的 URL
- `%width%` - 视频宽度
- `%height%` - 视频高度
- `%fps%` - 帧率（默认 8）
- `%duration%` - 时长（默认 5 秒）

### 示例：图片生成 Workflow

```json
{
  "3": {
    "class_type": "KSampler",
    "inputs": {
      "seed": %seed%,
      "steps": %steps%,
      "cfg": %cfg_scale%,
      "sampler_name": "euler",
      "scheduler": "normal",
      "denoise": 1.0,
      "model": ["4", 0],
      "positive": ["6", 0],
      "negative": ["7", 0],
      "latent_image": ["5", 0]
    }
  },
  "5": {
    "class_type": "EmptyLatentImage",
    "inputs": {
      "width": %width%,
      "height": %height%,
      "batch_size": 1
    }
  },
  "6": {
    "class_type": "CLIPTextEncode",
    "inputs": {
      "text": "%prompt%",
      "clip": ["4", 1]
    }
  },
  "7": {
    "class_type": "CLIPTextEncode",
    "inputs": {
      "text": "%negative_prompt%",
      "clip": ["4", 1]
    }
  }
}
```

### 示例：视频生成 Workflow

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
  }
}
```

## 文件要求

- **格式**：必须是有效的 JSON 文件
- **扩展名**：`.json`
- **大小**：最大 10MB
- **编码**：UTF-8

## 目录结构

```
huobao-drama/
├── workflows/                    # Workflow 文件目录
│   ├── sdxl_text2img.json       # 预设：SDXL 文生图
│   ├── video_wan2_2_14B_i2v.json # 预设：视频生成
│   └── my_custom_workflow.json   # 你上传的文件
```

## API 接口

如果需要通过 API 管理 workflow：

### 获取列表
```http
GET /api/v1/workflows?type=image
GET /api/v1/workflows?type=video
```

### 上传文件
```http
POST /api/v1/workflows/upload
Content-Type: multipart/form-data

file: <your-workflow.json>
```

### 获取文件内容
```http
GET /api/v1/workflows/:filename
```

### 删除文件
```http
DELETE /api/v1/workflows/:filename
```

## 注意事项

1. **文件命名**：
   - 建议使用英文和数字
   - 避免特殊字符
   - 建议包含类型标识（如 `_image`、`_video`）

2. **占位符**：
   - 占位符区分大小写
   - 数值型占位符不需要引号（如 `%width%`）
   - 字符串型占位符需要引号（如 `"%prompt%"`）

3. **ComfyUI 节点**：
   - 确保 workflow 中使用的节点在你的 ComfyUI 中已安装
   - 确保引用的模型文件存在于 ComfyUI 的 models 目录

4. **安全性**：
   - 系统会验证 JSON 格式
   - 文件名会自动清理，防止路径遍历攻击
   - 建议只上传可信来源的 workflow 文件

## 故障排除

### 问题：上传失败
**解决方案**：
- 检查文件是否为有效的 JSON 格式
- 确认文件大小不超过 10MB
- 确认文件扩展名为 `.json`

### 问题：Workflow 执行失败
**解决方案**：
- 检查 ComfyUI 服务器是否运行
- 查看 ComfyUI 日志，确认节点是否可用
- 验证占位符是否正确替换

### 问题：找不到上传的文件
**解决方案**：
- 点击"刷新列表"按钮
- 检查 `workflows/` 目录是否存在该文件
- 确认文件扩展名为 `.json`

## 相关文档

- [ComfyUI Workflow 指南](COMFYUI_WORKFLOW_GUIDE.md)
- [ComfyUI 集成文档](docs/COMFYUI_INTEGRATION.md)
- [ComfyUI 快速开始](docs/COMFYUI_QUICKSTART.md)
