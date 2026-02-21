# ComfyUI 集成总结

## 修改文件清单

### 后端代码

1. **新增文件**
   - `pkg/video/comfyui_client.go` - ComfyUI 客户端实现
   - `pkg/video/comfyui_client_test.go` - 单元测试

2. **修改文件**
   - `application/services/video_generation_service.go` - 添加 ComfyUI 支持
   - `application/services/video_merge_service.go` - 添加 ComfyUI 支持

### 前端代码

1. **修改文件**
   - `web/src/views/settings/AIConfig.vue` - 添加 ComfyUI 配置选项

### 文档

1. **新增文档**
   - `docs/COMFYUI_INTEGRATION.md` - 完整集成指南
   - `docs/comfyui-config-example.yaml` - 配置示例
   - `docs/comfyui-usage-example.sh` - 使用示例脚本
   - `docs/COMFYUI_SUMMARY.md` - 本文档

2. **更新文档**
   - `README-CN.md` - 添加 ComfyUI 功能说明和更新日志

## 核心功能

### 1. ComfyUI 客户端 (`comfyui_client.go`)

**主要功能：**
- 提交视频生成任务到 ComfyUI
- 查询任务状态（队列和历史）
- 自动构建图生视频工作流
- 支持自定义参数配置

**API 端点：**
- `POST /prompt` - 提交任务
- `GET /queue` - 查询队列
- `GET /history/{prompt_id}` - 查询历史
- `GET /view` - 获取生成的视频文件

**工作流节点：**
1. LoadImage - 加载输入图片
2. CLIPTextEncode - 编码提示词
3. CheckpointLoaderSimple - 加载模型
4. VideoLinearCFGGuidance - CFG 引导
5. SVD_img2vid_Conditioning - 图生视频条件
6. KSampler - 采样器
7. VAEDecode - VAE 解码
8. VHS_VideoCombine - 视频合成

### 2. 服务集成

**视频生成服务 (`video_generation_service.go`)：**
```go
case "comfyui":
    endpoint = "/prompt"
    queryEndpoint = "/history/{prompt_id}"
    return video.NewComfyUIClient(baseURL, apiKey, model, endpoint, queryEndpoint), nil
```

**视频合并服务 (`video_merge_service.go`)：**
```go
case "comfyui":
    endpoint = "/prompt"
    queryEndpoint = "/history/{prompt_id}"
    return video.NewComfyUIClient(config.BaseURL, config.APIKey, model, endpoint, queryEndpoint), nil
```

### 3. 前端配置

**AI 配置页面新增选项：**
```typescript
video: [
  // ... 其他提供商
  { 
    id: "comfyui", 
    name: "ComfyUI", 
    models: ["svd", "svd_xt", "custom"] 
  },
]
```

**端点配置：**
```typescript
else if (provider === "comfyui") {
  endpoint = "/prompt";
}
```

## 使用流程

### 1. 配置 ComfyUI

在系统设置 → AI 服务配置 → 视频生成：

```yaml
名称: ComfyUI-本地
厂商: ComfyUI
Base URL: http://localhost:8188
API Key: (留空)
模型: svd
优先级: 50
```

### 2. 生成视频

通过 API 或前端界面：

```bash
POST /api/v1/video-generations
{
  "storyboard_id": 123,
  "image_url": "http://example.com/image.jpg",
  "prompt": "A beautiful scene",
  "model": "svd",
  "duration": 5
}
```

### 3. 查询状态

系统自动轮询任务状态：
- 检查队列状态
- 查询历史记录
- 提取视频 URL
- 更新数据库

## 技术特点

### 1. 异步处理
- 任务提交后立即返回
- 后台轮询状态
- 完成后自动更新

### 2. 错误处理
- 连接失败重试
- 工作流错误检测
- 超时保护

### 3. 灵活配置
- 支持自定义端点
- 可选 API Key 认证
- 多模型支持

### 4. 扩展性
- 易于添加新节点
- 支持自定义工作流
- 模块化设计

## 测试

### 单元测试

```bash
cd pkg/video
go test -v
```

**测试覆盖：**
- 客户端创建
- 工作流构建
- URL 构建
- 参数验证

### 集成测试

使用提供的脚本：

```bash
chmod +x docs/comfyui-usage-example.sh
./docs/comfyui-usage-example.sh
```

## 性能考虑

### 1. 资源需求
- GPU: RTX 3060+ (8GB VRAM)
- RAM: 16GB+
- 存储: 根据模型大小

### 2. 优化建议
- 降低分辨率 (512x512)
- 减少帧数 (25-40 frames)
- 调整采样步数 (15-20 steps)

### 3. 并发控制
- 队列管理
- 轮询间隔 (5秒)
- 超时设置 (10分钟)

## 故障排查

### 常见问题

1. **连接失败**
   - 检查 ComfyUI 是否运行
   - 验证 Base URL
   - 检查防火墙

2. **工作流错误**
   - 确认模型文件存在
   - 检查节点配置
   - 查看 ComfyUI 日志

3. **任务超时**
   - 增加超时时间
   - 优化参数
   - 检查服务器性能

### 调试模式

启用调试日志：

```yaml
# config.yaml
app:
  debug: true
```

查看日志：
```
[ComfyUI] Sending request to: http://localhost:8188/prompt
[ComfyUI] Task created - PromptID: abc123
[ComfyUI] Task status - ID: abc123, Status: processing
```

## 未来改进

### 短期
- [ ] 支持更多 ComfyUI 节点
- [ ] 添加工作流模板库
- [ ] 优化错误提示

### 中期
- [ ] 支持批量处理
- [ ] 添加进度回调
- [ ] 实现工作流可视化编辑

### 长期
- [ ] 支持 ComfyUI 插件
- [ ] 分布式任务调度
- [ ] 性能监控和优化

## 参考资源

- [ComfyUI GitHub](https://github.com/comfyanonymous/ComfyUI)
- [ComfyUI API 文档](https://github.com/comfyanonymous/ComfyUI/wiki/API)
- [Stable Video Diffusion](https://stability.ai/stable-video)
- [项目集成指南](COMFYUI_INTEGRATION.md)

## 贡献者

- 初始实现: 2026-02-20
- 版本: v1.0.6

## 许可证

遵循项目主许可证 (CC BY-NC-SA 4.0)
