# ComfyUI Workflow 文件路径修复

## 问题描述
当用户上传自定义 ComfyUI workflow JSON 文件（如 `老王文生图，ZImage，强 API.json`）并在 AI 配置中选择该文件作为模型时，系统无法正确加载文件，因为代码直接使用文件名作为路径。

## 原因分析
1. 用户在 AI 配置中选择的模型值是文件名，例如：`"老王文生图，ZImage，强 API.json"`
2. 代码检测到 `.json` 后缀后，直接将文件名传递给 `loadWorkflowFromFile()`
3. `loadWorkflowFromFile()` 使用 `os.ReadFile()` 读取文件，但文件名不是完整路径
4. 实际文件存储在 `workflows/` 目录下，需要完整路径：`workflows/老王文生图，ZImage，强 API.json`

## 解决方案
在调用 `loadWorkflowFromFile()` 之前，检查文件路径：
- 如果不是绝对路径（不以 `/` 开头，也不包含 `:\`），则自动添加 `workflows/` 前缀
- 这样可以支持：
  - 相对文件名：`老王文生图，ZImage，强 API.json` → `workflows/老王文生图，ZImage，强 API.json`
  - 相对路径：`workflows/custom.json` → `workflows/workflows/custom.json`（虽然不推荐）
  - 绝对路径：`/path/to/workflow.json` 或 `C:\path\to\workflow.json`（保持不变）

## 代码修改

### 图片生成客户端 (`pkg/image/comfyui_image_client.go`)
```go
if c.Model != "" && (len(c.Model) > 5 && c.Model[len(c.Model)-5:] == ".json") {
    // Model 是 JSON 文件名，构建完整路径
    workflowPath := c.Model
    // 如果不是绝对路径，则添加 workflows/ 前缀
    if !strings.HasPrefix(workflowPath, "/") && !strings.Contains(workflowPath, ":\\") {
        workflowPath = "workflows/" + workflowPath
    }
    
    workflow, err = c.loadWorkflowFromFile(workflowPath, prompt, options)
    if err != nil {
        return nil, fmt.Errorf("load workflow from file: %w", err)
    }
}
```

### 视频生成客户端 (`pkg/video/comfyui_client.go`)
```go
if c.Model != "" && (len(c.Model) > 5 && c.Model[len(c.Model)-5:] == ".json") {
    // Model 是 JSON 文件名，构建完整路径
    workflowPath := c.Model
    // 如果不是绝对路径，则添加 workflows/ 前缀
    if !strings.HasPrefix(workflowPath, "/") && !strings.Contains(workflowPath, ":\\") {
        workflowPath = "workflows/" + workflowPath
    }
    
    workflow, err = c.loadWorkflowFromFile(workflowPath, imageURL, prompt, options)
    if err != nil {
        return nil, fmt.Errorf("load workflow from file: %w", err)
    }
}
```

## 占位符替换机制
`loadWorkflowFromFile()` 函数会自动替换 workflow JSON 中的占位符：

### 图片生成占位符
- `%prompt%` - 用户输入的提示词
- `%negative_prompt%` - 负面提示词
- `%width%` - 图片宽度
- `%height%` - 图片高度
- `%steps%` - 采样步数
- `%cfg_scale%` - CFG 比例
- `%seed%` - 随机种子（如果指定）

### 视频生成占位符
- `%prompt%` - 用户输入的提示词
- `%image_url%` - 输入图片 URL
- `%width%` - 视频宽度
- `%height%` - 视频高度
- `%fps%` - 帧率
- `%duration%` - 时长（秒）

## 使用示例

### 1. 上传 workflow 文件
在 AI 配置界面，选择 ComfyUI 作为厂商，点击"上传 Workflow"按钮，上传包含 `%prompt%` 占位符的 JSON 文件。

### 2. Workflow JSON 示例
```json
{
  "18": {
    "inputs": {
      "value": "%prompt%"
    },
    "class_type": "PrimitiveStringMultiline",
    "_meta": {
      "title": "字符串（多行）"
    }
  }
}
```

### 3. 选择模型
在模型下拉框中选择上传的文件名，例如：`老王文生图，ZImage，强 API.json`

### 4. 生成图片
当用户输入提示词并生成图片时：
1. 系统读取 `workflows/老王文生图，ZImage，强 API.json`
2. 将 JSON 转换为字符串
3. 替换所有 `%prompt%` 为用户输入的提示词
4. 解析回 JSON 对象
5. 发送到 ComfyUI API

## 测试验证
```bash
# 启动服务
./start-prod.cmd

# 访问应用
http://localhost:5678

# 测试步骤：
1. 进入 AI 配置
2. 添加 ComfyUI 配置（BaseURL: http://192.168.2.76:8080）
3. 上传包含 %prompt% 的 workflow JSON
4. 选择上传的文件作为模型
5. 生成角色图片
6. 检查后端日志，应该看到：
   [ComfyUI-Image] Loaded workflow from: workflows/xxx.json
   [ComfyUI-Image] Replaced %prompt% with: [用户提示词]
```

## 注意事项
1. Workflow JSON 文件必须是有效的 JSON 格式
2. 占位符区分大小写，必须使用小写（如 `%prompt%` 而不是 `%PROMPT%`）
3. 文件名支持中文和特殊字符
4. 建议使用相对文件名，系统会自动添加 `workflows/` 前缀
5. 上传的文件会保存在 `./workflows/` 目录下
