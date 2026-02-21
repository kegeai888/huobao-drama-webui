# ComfyUI 图片上传修复

## 问题描述

ComfyUI 视频生成一直报错 "error" 状态，但 ComfyUI 后台没有收到任何任务。

## 根本原因

`pkg/video/comfyui_client.go` 中的 `uploadImage()` 函数只能处理 HTTP URL，无法处理本地文件路径。

当应用传递本地文件路径（如 `data/storage/images/xxx.png`）时：
- `http.Client.Get()` 尝试将其作为 HTTP URL 下载
- 下载失败，导致图片上传失败
- 代码降级使用原始路径作为 LoadImage 节点的输入
- LoadImage 节点需要 ComfyUI 服务器上的实际文件，导致任务失败

## 解决方案

修改 `uploadImage()` 函数，支持两种图片来源：

1. **HTTP URL**: 使用 `http.Client.Get()` 下载
2. **本地文件**: 使用 `os.ReadFile()` 直接读取

```go
func (c *ComfyUIClient) uploadImage(imageURL string) (string, error) {
	var imageData []byte
	var err error

	// 判断是本地文件还是 HTTP URL
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		// HTTP URL - 下载图片
		resp, err := c.HTTPClient.Get(imageURL)
		if err != nil {
			return "", fmt.Errorf("download image: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("download image failed: status %d", resp.StatusCode)
		}

		imageData, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("read image data: %w", err)
		}
	} else {
		// 本地文件路径 - 直接读取
		imageData, err = os.ReadFile(imageURL)
		if err != nil {
			return "", fmt.Errorf("read local image file: %w", err)
		}
		fmt.Printf("[ComfyUI] Read local image file: %s (%d bytes)\n", imageURL, len(imageData))
	}
	
	// ... 后续上传逻辑保持不变
}
```

## 测试结果

使用测试脚本 `test_video_id7.go` 验证：

```bash
go run test_video_id7.go
```

测试结果：
- ✓ 图片成功读取（1.5MB）
- ✓ 图片成功上传到 ComfyUI
- ✓ Workflow 成功提交
- ✓ 视频生成成功完成（约 5 分钟）

## 影响范围

- `pkg/video/comfyui_client.go` - 视频生成客户端
- 所有使用 ComfyUI 进行视频生成的功能

## 相关文件

- `pkg/video/comfyui_client.go` - 修复的主文件
- `test_video_id7.go` - 测试脚本
- `workflows/video_wan2_2_14B_i2v_api.json` - 测试使用的 workflow

## 注意事项

1. 图片上传需要时间，大图片（>1MB）可能需要几秒钟
2. Wan2.2 14B 模型生成视频需要较长时间（约 5-10 分钟）
3. 建议设置足够的超时时间（600 秒以上）
