package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ComfyUIImageClient ComfyUI 图片生成客户端
type ComfyUIImageClient struct {
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
	PromptID   string                 `json:"prompt_id"`
	Number     int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors,omitempty"`
}

// ComfyUIHistoryResponse ComfyUI 历史记录响应
type ComfyUIHistoryResponse map[string]struct {
	Prompt  []interface{} `json:"prompt"`
	Outputs map[string]struct {
		Images []struct {
			Filename  string `json:"filename"`
			Subfolder string `json:"subfolder"`
			Type      string `json:"type"`
		} `json:"images,omitempty"`
	} `json:"outputs"`
	Status struct {
		StatusStr string          `json:"status_str"`
		Completed bool            `json:"completed"`
		Messages  [][]interface{} `json:"messages"`
	} `json:"status"`
}

// ComfyUIQueueResponse ComfyUI 队列状态响应
type ComfyUIQueueResponse struct {
	QueueRunning [][]interface{} `json:"queue_running"`
	QueuePending [][]interface{} `json:"queue_pending"`
}

// NewComfyUIImageClient 创建 ComfyUI 图片客户端
func NewComfyUIImageClient(baseURL, apiKey, model, endpoint, queryEndpoint string) *ComfyUIImageClient {
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

	return &ComfyUIImageClient{
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

// GenerateImage 生成图片
func (c *ComfyUIImageClient) GenerateImage(prompt string, opts ...ImageOption) (*ImageResult, error) {
	options := &ImageOptions{
		Width:    512,
		Height:   512,
		Steps:    20,
		CfgScale: 7.0,
	}

	for _, opt := range opts {
		opt(options)
	}

	// 构建 ComfyUI workflow
	// 如果 model 是 JSON 文件路径，则加载外部 workflow
	// 否则使用内置的 SDXL workflow
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
		
		workflow, err = c.loadWorkflowFromFile(workflowPath, prompt, options)
		if err != nil {
			return nil, fmt.Errorf("load workflow from file: %w", err)
		}
	} else {
		// 使用内置 workflow
		workflow = c.buildWorkflow(prompt, options)
	}

	reqBody := ComfyUIPromptRequest{
		Prompt:     workflow,
		ClientID:   fmt.Sprintf("huobao-drama-img-%d", time.Now().Unix()),
		FrontQueue: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.BaseURL + c.Endpoint
	fmt.Printf("[ComfyUI-Image] Sending request to: %s\n", url)

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

	fmt.Printf("[ComfyUI-Image] Task created - PromptID: %s\n", promptResp.PromptID)

	result := &ImageResult{
		TaskID:    promptResp.PromptID,
		Status:    "processing",
		Completed: false,
		Width:     options.Width,
		Height:    options.Height,
	}

	return result, nil
}

// GetTaskStatus 查询任务状态
func (c *ComfyUIImageClient) GetTaskStatus(taskID string) (*ImageResult, error) {
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
		return &ImageResult{
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
		return &ImageResult{
			TaskID:    taskID,
			Status:    "processing",
			Completed: false,
		}, nil
	}

	// 检查任务状态
	if !taskHistory.Status.Completed {
		return &ImageResult{
			TaskID:    taskID,
			Status:    "processing",
			Completed: false,
		}, nil
	}

	// 查找输出图片
	var imageURL string
	for _, output := range taskHistory.Outputs {
		if len(output.Images) > 0 {
			img := output.Images[0]
			// 构建图片 URL
			imageURL = fmt.Sprintf("%s/view?filename=%s&subfolder=%s&type=%s",
				c.BaseURL, img.Filename, img.Subfolder, img.Type)
			break
		}
	}

	if imageURL == "" {
		return nil, fmt.Errorf("no image found in outputs")
	}

	return &ImageResult{
		TaskID:    taskID,
		Status:    "completed",
		ImageURL:  imageURL,
		Completed: true,
	}, nil
}

// loadWorkflowFromFile 从 JSON 文件加载 workflow 并替换占位符
func (c *ComfyUIImageClient) loadWorkflowFromFile(filePath, prompt string, options *ImageOptions) (map[string]interface{}, error) {
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
	workflowStr = strings.ReplaceAll(workflowStr, "%negative_prompt%", escapeJSON(options.NegativePrompt))
	workflowStr = strings.ReplaceAll(workflowStr, "%width%", fmt.Sprintf("%d", options.Width))
	workflowStr = strings.ReplaceAll(workflowStr, "%height%", fmt.Sprintf("%d", options.Height))
	workflowStr = strings.ReplaceAll(workflowStr, "%steps%", fmt.Sprintf("%d", options.Steps))
	workflowStr = strings.ReplaceAll(workflowStr, "%cfg_scale%", fmt.Sprintf("%.1f", options.CfgScale))

	if options.Seed > 0 {
		workflowStr = strings.ReplaceAll(workflowStr, "%seed%", fmt.Sprintf("%d", options.Seed))
	}

	// 解析为 map
	var workflow map[string]interface{}
	if err := json.Unmarshal([]byte(workflowStr), &workflow); err != nil {
		return nil, fmt.Errorf("parse workflow JSON: %w", err)
	}

	fmt.Printf("[ComfyUI-Image] Loaded workflow from: %s\n", filePath)
	fmt.Printf("[ComfyUI-Image] Replaced %%prompt%% with: %s\n", prompt)

	return workflow, nil
}

// buildWorkflow 构建内置的 SDXL workflow
func (c *ComfyUIImageClient) buildWorkflow(prompt string, options *ImageOptions) map[string]interface{} {
	model := c.Model
	if model == "" {
		model = "sdxl"
	}

	// 根据模型选择 checkpoint
	var checkpointName string
	switch model {
	case "sdxl":
		checkpointName = "sd_xl_base_1.0.safetensors"
	case "sd15":
		checkpointName = "v1-5-pruned-emaonly.safetensors"
	case "flux":
		checkpointName = "flux1-dev.safetensors"
	default:
		checkpointName = "sd_xl_base_1.0.safetensors"
	}

	seed := options.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	return map[string]interface{}{
		"3": map[string]interface{}{
			"class_type": "KSampler",
			"inputs": map[string]interface{}{
				"seed":       seed,
				"steps":      options.Steps,
				"cfg":        options.CfgScale,
				"sampler_name": "euler",
				"scheduler":    "normal",
				"denoise":      1.0,
				"model":        []interface{}{"4", 0},
				"positive":     []interface{}{"6", 0},
				"negative":     []interface{}{"7", 0},
				"latent_image": []interface{}{"5", 0},
			},
		},
		"4": map[string]interface{}{
			"class_type": "CheckpointLoaderSimple",
			"inputs": map[string]interface{}{
				"ckpt_name": checkpointName,
			},
		},
		"5": map[string]interface{}{
			"class_type": "EmptyLatentImage",
			"inputs": map[string]interface{}{
				"width":       options.Width,
				"height":      options.Height,
				"batch_size":  1,
			},
		},
		"6": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": prompt,
				"clip": []interface{}{"4", 1},
			},
		},
		"7": map[string]interface{}{
			"class_type": "CLIPTextEncode",
			"inputs": map[string]interface{}{
				"text": options.NegativePrompt,
				"clip": []interface{}{"4", 1},
			},
		},
		"8": map[string]interface{}{
			"class_type": "VAEDecode",
			"inputs": map[string]interface{}{
				"samples": []interface{}{"3", 0},
				"vae":     []interface{}{"4", 2},
			},
		},
		"9": map[string]interface{}{
			"class_type": "SaveImage",
			"inputs": map[string]interface{}{
				"images":          []interface{}{"8", 0},
				"filename_prefix": "ComfyUI",
			},
		},
	}
}
