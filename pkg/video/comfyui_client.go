package video

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// ComfyUIClient ComfyUI 视频生成客户端
// 支持通过 ComfyUI API 进行视频生成
type ComfyUIClient struct {
	BaseURL       string
	APIKey        string
	Model         string
	Endpoint      string
	QueryEndpoint string
	HTTPClient    *http.Client
}

// ComfyUIPromptRequest ComfyUI 提示词请求结构
type ComfyUIPromptRequest struct {
	Prompt     map[string]interface{} `json:"prompt"`
	ClientID   string                 `json:"client_id,omitempty"`
	ExtraData  map[string]interface{} `json:"extra_data,omitempty"`
	FrontQueue bool                   `json:"front,omitempty"`
}

// ComfyUIPromptResponse ComfyUI 提示词响应
type ComfyUIPromptResponse struct {
	PromptID string                 `json:"prompt_id"`
	Number   int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors,omitempty"`
}

// ComfyUIHistoryResponse ComfyUI 历史记录响应
type ComfyUIHistoryResponse map[string]struct {
	Prompt []interface{} `json:"prompt"`
	Outputs map[string]struct {
		Images []struct {
			Filename  string `json:"filename"`
			Subfolder string `json:"subfolder"`
			Type      string `json:"type"`
		} `json:"images,omitempty"`
		Videos []struct {
			Filename  string `json:"filename"`
			Subfolder string `json:"subfolder"`
			Type      string `json:"type"`
		} `json:"videos,omitempty"`
		Gifs []struct {
			Filename  string `json:"filename"`
			Subfolder string `json:"subfolder"`
			Type      string `json:"type"`
		} `json:"gifs,omitempty"`
	} `json:"outputs"`
	Status struct {
		StatusStr  string `json:"status_str"`
		Completed  bool   `json:"completed"`
		Messages   [][]interface{} `json:"messages"`
	} `json:"status"`
}

// ComfyUIQueueResponse ComfyUI 队列状态响应
type ComfyUIQueueResponse struct {
	QueueRunning [][]interface{} `json:"queue_running"`
	QueuePending [][]interface{} `json:"queue_pending"`
}

// NewComfyUIClient 创建 ComfyUI 客户端
func NewComfyUIClient(baseURL, apiKey, model, endpoint, queryEndpoint string) *ComfyUIClient {
	// 规范化 BaseURL，移除尾部斜杠
	baseURL = strings.TrimRight(baseURL, "/")
	
	if endpoint == "" {
		endpoint = "/prompt"
	}
	// 确保 endpoint 以斜杠开头
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	
	if queryEndpoint == "" {
		queryEndpoint = "/history/{prompt_id}"
	}
	// 确保 queryEndpoint 以斜杠开头
	if !strings.HasPrefix(queryEndpoint, "/") {
		queryEndpoint = "/" + queryEndpoint
	}

	return &ComfyUIClient{
		BaseURL:       baseURL,
		APIKey:        apiKey,
		Model:         model,
		Endpoint:      endpoint,
		QueryEndpoint: queryEndpoint,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}
}

// GenerateVideo 生成视频
func (c *ComfyUIClient) GenerateVideo(imageURL, prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		FPS:         8,
		Resolution:  "512x512",
		AspectRatio: "16:9",
	}

	for _, opt := range opts {
		opt(options)
	}

	// 如果使用外部 workflow，处理图片
	var imageFilename string
	if c.Model != "" && (len(c.Model) > 5 && c.Model[len(c.Model)-5:] == ".json") {
		// 检查是否是 base64 数据
		if strings.HasPrefix(imageURL, "data:image/") {
			// 将 base64 数据转换为临时文件并上传
			fmt.Printf("[ComfyUI] Converting base64 to temp file and uploading (length: %d)\n", len(imageURL))
			filename, err := c.uploadBase64Image(imageURL)
			if err != nil {
				return nil, fmt.Errorf("upload base64 image: %w", err)
			}
			imageFilename = filename
			fmt.Printf("[ComfyUI] Base64 image uploaded successfully: %s\n", filename)
		} else {
			// 上传图片到 ComfyUI
			filename, err := c.uploadImage(imageURL)
			if err != nil {
				fmt.Printf("[ComfyUI] Warning: Failed to upload image, will use URL directly: %v\n", err)
				imageFilename = imageURL // 降级使用 URL
			} else {
				imageFilename = filename
				fmt.Printf("[ComfyUI] Image uploaded successfully: %s\n", filename)
			}
		}
	}

	// 构建 ComfyUI workflow
	// 如果 model 是 JSON 文件路径，则加载外部 workflow
	// 否则使用内置的 SVD workflow
	var workflow map[string]interface{}
	var err error
	
	if c.Model != "" && (len(c.Model) > 5 && c.Model[len(c.Model)-5:] == ".json") {
		// Model 是 JSON 文件名，构建完整路径
		workflowPath := c.Model
		
		// 移除可能已经存在的 workflows/ 或 workflows\ 前缀
		workflowPath = strings.TrimPrefix(workflowPath, "workflows/")
		workflowPath = strings.TrimPrefix(workflowPath, "workflows\\")
		
		// 如果不是绝对路径，则添加 workflows/ 前缀
		if !strings.HasPrefix(workflowPath, "/") && !strings.Contains(workflowPath, ":\\") {
			workflowPath = "workflows/" + workflowPath
		}
		
		workflow, err = c.loadWorkflowFromFile(workflowPath, imageURL, imageFilename, prompt, options)
		if err != nil {
			return nil, fmt.Errorf("load workflow from file: %w", err)
		}
	} else {
		// 使用内置 workflow
		workflow = c.buildWorkflow(imageURL, prompt, options)
	}

	reqBody := ComfyUIPromptRequest{
		Prompt:     workflow,
		ClientID:   fmt.Sprintf("huobao-drama-%d", time.Now().Unix()),
		FrontQueue: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 调试：打印请求体（截断以避免日志过长）
	reqBodyPreview := string(jsonData)
	if len(reqBodyPreview) > 500 {
		reqBodyPreview = reqBodyPreview[:500] + "..."
	}
	fmt.Printf("[ComfyUI] Request body preview: %s\n", reqBodyPreview)

	url := c.BaseURL + c.Endpoint
	fmt.Printf("[ComfyUI] Sending request to: %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var promptResp ComfyUIPromptResponse
	if err := json.Unmarshal(body, &promptResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(promptResp.NodeErrors) > 0 {
		return nil, fmt.Errorf("workflow error: %v", promptResp.NodeErrors)
	}

	fmt.Printf("[ComfyUI] Task created - PromptID: %s\n", promptResp.PromptID)

	result := &VideoResult{
		TaskID:    promptResp.PromptID,
		Status:    "processing",
		Completed: false,
	}

	return result, nil
}

// GetTaskStatus 查询任务状态
func (c *ComfyUIClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	// 先检查队列状态
	queueURL := c.BaseURL + "/queue"
	queueReq, err := http.NewRequest("GET", queueURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create queue request: %w", err)
	}

	if c.APIKey != "" {
		queueReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	queueResp, err := c.HTTPClient.Do(queueReq)
	if err != nil {
		return nil, fmt.Errorf("send queue request: %w", err)
	}
	defer queueResp.Body.Close()

	queueBody, err := io.ReadAll(queueResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read queue response: %w", err)
	}

	var queueStatus ComfyUIQueueResponse
	if err := json.Unmarshal(queueBody, &queueStatus); err != nil {
		return nil, fmt.Errorf("parse queue response: %w", err)
	}

	// 检查是否在队列中
	inQueue := false
	for _, item := range queueStatus.QueueRunning {
		if len(item) > 1 {
			if promptID, ok := item[1].(string); ok && promptID == taskID {
				inQueue = true
				break
			}
		}
	}
	for _, item := range queueStatus.QueuePending {
		if len(item) > 1 {
			if promptID, ok := item[1].(string); ok && promptID == taskID {
				inQueue = true
				break
			}
		}
	}

	if inQueue {
		return &VideoResult{
			TaskID:    taskID,
			Status:    "processing",
			Completed: false,
		}, nil
	}

	// 查询历史记录
	historyURL := c.BaseURL + "/history/" + taskID
	historyReq, err := http.NewRequest("GET", historyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create history request: %w", err)
	}

	if c.APIKey != "" {
		historyReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	historyResp, err := c.HTTPClient.Do(historyReq)
	if err != nil {
		return nil, fmt.Errorf("send history request: %w", err)
	}
	defer historyResp.Body.Close()

	historyBody, err := io.ReadAll(historyResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read history response: %w", err)
	}

	var history ComfyUIHistoryResponse
	if err := json.Unmarshal(historyBody, &history); err != nil {
		return nil, fmt.Errorf("parse history response: %w", err)
	}

	taskHistory, exists := history[taskID]
	if !exists {
		return &VideoResult{
			TaskID:    taskID,
			Status:    "processing",
			Completed: false,
		}, nil
	}

	result := &VideoResult{
		TaskID:    taskID,
		Status:    taskHistory.Status.StatusStr,
		Completed: taskHistory.Status.Completed,
	}

	// 提取视频 URL
	if taskHistory.Status.Completed {
		fmt.Printf("[ComfyUI] Task completed, extracting video URL. Outputs count: %d\n", len(taskHistory.Outputs))
		for nodeID, output := range taskHistory.Outputs {
			fmt.Printf("[ComfyUI] Node %s - Videos: %d, Gifs: %d, Images: %d\n", 
				nodeID, len(output.Videos), len(output.Gifs), len(output.Images))
			
			// 优先查找视频
			if len(output.Videos) > 0 {
				video := output.Videos[0]
				result.VideoURL = c.buildFileURL(video.Filename, video.Subfolder, video.Type)
				fmt.Printf("[ComfyUI] Found video: %s (subfolder: %s, type: %s)\n", 
					video.Filename, video.Subfolder, video.Type)
				fmt.Printf("[ComfyUI] Video URL: %s\n", result.VideoURL)
				break
			}
			// 其次查找 GIF
			if len(output.Gifs) > 0 {
				gif := output.Gifs[0]
				result.VideoURL = c.buildFileURL(gif.Filename, gif.Subfolder, gif.Type)
				fmt.Printf("[ComfyUI] Found GIF: %s (subfolder: %s, type: %s)\n", 
					gif.Filename, gif.Subfolder, gif.Type)
				fmt.Printf("[ComfyUI] GIF URL: %s\n", result.VideoURL)
				break
			}
			// 最后检查 images 数组中的视频文件（某些节点将视频放在 images 中）
			if len(output.Images) > 0 {
				for _, img := range output.Images {
					// 检查文件扩展名是否为视频格式
					if strings.HasSuffix(strings.ToLower(img.Filename), ".mp4") ||
						strings.HasSuffix(strings.ToLower(img.Filename), ".webm") ||
						strings.HasSuffix(strings.ToLower(img.Filename), ".mov") ||
						strings.HasSuffix(strings.ToLower(img.Filename), ".avi") {
						result.VideoURL = c.buildFileURL(img.Filename, img.Subfolder, img.Type)
						fmt.Printf("[ComfyUI] Found video in images array: %s (subfolder: %s, type: %s)\n", 
							img.Filename, img.Subfolder, img.Type)
						fmt.Printf("[ComfyUI] Video URL: %s\n", result.VideoURL)
						break
					}
				}
				if result.VideoURL != "" {
					break
				}
			}
		}

		if result.VideoURL == "" {
			result.Error = "no video output found"
			fmt.Printf("[ComfyUI] ERROR: No video output found in completed task\n")
		}
	}

	fmt.Printf("[ComfyUI] Task status - ID: %s, Status: %s, Completed: %v\n", 
		taskID, result.Status, result.Completed)

	return result, nil
}

// buildFileURL 构建文件访问 URL
func (c *ComfyUIClient) buildFileURL(filename, subfolder, fileType string) string {
	url := c.BaseURL + "/view"
	params := fmt.Sprintf("?filename=%s", filename)
	
	if subfolder != "" {
		params += fmt.Sprintf("&subfolder=%s", subfolder)
	}
	
	if fileType != "" {
		params += fmt.Sprintf("&type=%s", fileType)
	}
	
	return url + params
}

// buildWorkflow 构建 ComfyUI 工作流
// 这是一个基础的图生视频工作流模板，可以根据实际需求自定义
func (c *ComfyUIClient) buildWorkflow(imageURL, prompt string, options *VideoOptions) map[string]interface{} {
	// 基础工作流：加载图片 -> 视频生成 -> 保存
	workflow := map[string]interface{}{
		"1": map[string]interface{}{
			"class_type": "LoadImage",
			"inputs": map[string]interface{}{
				"image": imageURL,
			},
		},
		"2": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": prompt,
				"clip": []interface{}{"4", 1},
			},
		},
		"3": map[string]interface{}{
			"class_type": "VideoLinearCFGGuidance",
			"inputs": map[string]interface{}{
				"model":      []interface{}{"4", 0},
				"min_cfg":    1.0,
			},
		},
		"4": map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{
				"ckpt_name": c.Model,
			},
		},
		"5": map[string]interface{}{
			"class_type": "SVD_img2vid_Conditioning",
			"inputs": map[string]interface{}{
				"width":         512,
				"height":        512,
				"video_frames":  options.Duration * options.FPS,
				"motion_bucket_id": 127,
				"fps":           options.FPS,
				"augmentation_level": 0,
				"clip_vision":   []interface{}{"4", 2},
				"init_image":    []interface{}{"1", 0},
				"vae":           []interface{}{"4", 3},
			},
		},
		"6": map[string]interface{}{
			"class_type": "KSampler",
			"inputs": map[string]interface{}{
				"seed":       time.Now().Unix(),
				"steps":      20,
				"cfg":        2.5,
				"sampler_name": "euler",
				"scheduler":  "karras",
				"denoise":    1.0,
				"model":      []interface{}{"3", 0},
				"positive":   []interface{}{"5", 0},
				"negative":   []interface{}{"5", 1},
				"latent_image": []interface{}{"5", 2},
			},
		},
		"7": map[string]interface{}{
			"class_type": "VAEDecode",
			"inputs": map[string]interface{}{
				"samples": []interface{}{"6", 0},
				"vae":     []interface{}{"4", 3},
			},
		},
		"8": map[string]interface{}{
			"class_type": "VHS_VideoCombine",
			"inputs": map[string]interface{}{
				"frame_rate":    options.FPS,
				"loop_count":    0,
				"filename_prefix": "huobao_drama",
				"format":        "video/h264-mp4",
				"images":        []interface{}{"7", 0},
			},
		},
	}

	return workflow
}


// loadWorkflowFromFile 从 JSON 文件加载 workflow 并替换占位符
func (c *ComfyUIClient) loadWorkflowFromFile(filePath, imageURL, imageFilename, prompt string, options *VideoOptions) (map[string]interface{}, error) {
	// 读取 JSON 文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read workflow file: %w", err)
	}

	// 将 JSON 转换为字符串以便替换占位符
	workflowStr := string(data)

	// JSON 转义函数：转义特殊字符以便安全地插入 JSON 字符串
	escapeJSON := func(s string) string {
		// 使用 json.Marshal 来正确转义字符串
		b, _ := json.Marshal(s)
		// 去掉首尾的引号
		return string(b[1 : len(b)-1])
	}

	// 替换占位符（使用 JSON 转义）
	workflowStr = strings.ReplaceAll(workflowStr, "%prompt%", escapeJSON(prompt))
	workflowStr = strings.ReplaceAll(workflowStr, "%image_url%", escapeJSON(imageURL))
	workflowStr = strings.ReplaceAll(workflowStr, "%image_filename%", escapeJSON(imageFilename))
	
	// 替换其他可能的占位符
	if options != nil {
		// 从 Resolution 解析宽高 (例如 "512x512")
		width, height := 512, 512
		if options.Resolution != "" {
			fmt.Sscanf(options.Resolution, "%dx%d", &width, &height)
		}
		
		workflowStr = strings.ReplaceAll(workflowStr, "%width%", fmt.Sprintf("%d", width))
		workflowStr = strings.ReplaceAll(workflowStr, "%height%", fmt.Sprintf("%d", height))
		workflowStr = strings.ReplaceAll(workflowStr, "%fps%", fmt.Sprintf("%d", options.FPS))
		workflowStr = strings.ReplaceAll(workflowStr, "%duration%", fmt.Sprintf("%d", options.Duration))
	}

	// 解析为 map
	var workflow map[string]interface{}
	if err := json.Unmarshal([]byte(workflowStr), &workflow); err != nil {
		return nil, fmt.Errorf("parse workflow JSON: %w", err)
	}

	fmt.Printf("[ComfyUI] Loaded workflow from: %s\n", filePath)
	fmt.Printf("[ComfyUI] Replaced %%prompt%% with: %s\n", prompt)
	fmt.Printf("[ComfyUI] Replaced %%image_filename%% with: %s\n", imageFilename)

	return workflow, nil
}


// uploadImage 上传图片到 ComfyUI 服务器
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

	// 生成唯一文件名
	filename := fmt.Sprintf("huobao_drama_%d.jpg", time.Now().Unix())

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("write image data: %w", err)
	}

	// 添加其他字段（如果需要）
	writer.WriteField("overwrite", "true")

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close writer: %w", err)
	}

	// 上传到 ComfyUI
	uploadURL := c.BaseURL + "/upload/image"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	uploadResp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload image: %w", err)
	}
	defer uploadResp.Body.Close()

	uploadBody, err := io.ReadAll(uploadResp.Body)
	if err != nil {
		return "", fmt.Errorf("read upload response: %w", err)
	}

	if uploadResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed (status %d): %s", uploadResp.StatusCode, string(uploadBody))
	}

	// 解析响应获取文件名
	var uploadResult map[string]interface{}
	if err := json.Unmarshal(uploadBody, &uploadResult); err != nil {
		// 如果解析失败，直接返回我们生成的文件名
		return filename, nil
	}

	// 尝试从响应中获取文件名
	if name, ok := uploadResult["name"].(string); ok {
		return name, nil
	}

	return filename, nil
}

// uploadBase64Image 将 base64 图片数据上传到 ComfyUI 服务器
func (c *ComfyUIClient) uploadBase64Image(base64Data string) (string, error) {
	// 解析 base64 数据
	// 格式: data:image/png;base64,iVBORw0KGgo...
	parts := strings.SplitN(base64Data, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 data format")
	}

	// 解码 base64
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	// 从 MIME 类型推断文件扩展名
	ext := ".jpg"
	if strings.Contains(parts[0], "image/png") {
		ext = ".png"
	} else if strings.Contains(parts[0], "image/jpeg") || strings.Contains(parts[0], "image/jpg") {
		ext = ".jpg"
	} else if strings.Contains(parts[0], "image/webp") {
		ext = ".webp"
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("huobao_drama_%d%s", time.Now().Unix(), ext)

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("write image data: %w", err)
	}

	// 添加其他字段
	writer.WriteField("overwrite", "true")

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close writer: %w", err)
	}

	// 上传到 ComfyUI
	uploadURL := c.BaseURL + "/upload/image"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	uploadResp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload image: %w", err)
	}
	defer uploadResp.Body.Close()

	uploadBody, err := io.ReadAll(uploadResp.Body)
	if err != nil {
		return "", fmt.Errorf("read upload response: %w", err)
	}

	if uploadResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed (status %d): %s", uploadResp.StatusCode, string(uploadBody))
	}

	// 解析响应获取文件名
	var uploadResult map[string]interface{}
	if err := json.Unmarshal(uploadBody, &uploadResult); err != nil {
		// 如果解析失败，直接返回我们生成的文件名
		return filename, nil
	}

	// 尝试从响应中获取文件名
	if name, ok := uploadResult["name"].(string); ok {
		return name, nil
	}

	return filename, nil
}
