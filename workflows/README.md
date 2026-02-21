# ComfyUI 工作流目录

此目录用于存放 ComfyUI 工作流配置文件。

## 📁 目录说明

- 用户可以将自定义的 ComfyUI 工作流 JSON 文件放在此目录
- 系统会自动扫描此目录下的所有 `.json` 文件
- 工作流文件可以通过 Web 界面选择和使用

## 🔒 Git 忽略规则

为了保护用户隐私和避免提交大量自定义配置：
- 所有 `.json` 工作流文件已被 `.gitignore` 排除
- 仅保留此 README 和目录结构

## 📝 工作流文件格式

工作流文件应该是标准的 ComfyUI API 格式 JSON 文件，例如：

```json
{
  "prompt": {
    "1": {
      "class_type": "LoadImage",
      "inputs": {
        "image": "example.png"
      }
    },
    ...
  }
}
```

## 🎯 使用方法

1. 从 ComfyUI 导出工作流（API 格式）
2. 将 JSON 文件放入此目录
3. 在 Web 界面的工作流选择器中选择使用

## 📚 相关文档

- [ComfyUI 集成指南](../docs/COMFYUI_INTEGRATION.md)
- [ComfyUI 快速开始](../docs/COMFYUI_QUICKSTART.md)
- [工作流配置指南](../docs/comfyui/WORKFLOW_GUIDE.md)
