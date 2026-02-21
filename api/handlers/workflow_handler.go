package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type WorkflowHandler struct {
	workflowDir string
	log         *logger.Logger
}

func NewWorkflowHandler(workflowDir string, log *logger.Logger) *WorkflowHandler {
	// 确保 workflow 目录存在
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		log.Errorw("Failed to create workflow directory", "error", err)
	}

	return &WorkflowHandler{
		workflowDir: workflowDir,
		log:         log,
	}
}

// ListWorkflows 列出所有 workflow 文件
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	workflowType := c.Query("type") // image 或 video

	files, err := os.ReadDir(h.workflowDir)
	if err != nil {
		h.log.Errorw("Failed to read workflow directory", "error", err)
		response.InternalError(c, "Failed to read workflows")
		return
	}

	var workflows []map[string]interface{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 只返回 .json 文件
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			continue
		}

		// 根据文件名判断类型（可选过滤）
		if workflowType != "" {
			if workflowType == "image" && !strings.Contains(strings.ToLower(file.Name()), "image") && !strings.Contains(strings.ToLower(file.Name()), "img") && !strings.Contains(strings.ToLower(file.Name()), "text2img") {
				continue
			}
			if workflowType == "video" && !strings.Contains(strings.ToLower(file.Name()), "video") && !strings.Contains(strings.ToLower(file.Name()), "i2v") && !strings.Contains(strings.ToLower(file.Name()), "t2v") {
				continue
			}
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		workflows = append(workflows, map[string]interface{}{
			"name":        file.Name(),
			"path":        filepath.Join("workflows", file.Name()),
			"size":        info.Size(),
			"modified_at": info.ModTime(),
		})
	}

	response.Success(c, gin.H{
		"workflows": workflows,
		"count":     len(workflows),
	})
}

// UploadWorkflow 上传 workflow 文件
func (h *WorkflowHandler) UploadWorkflow(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file uploaded")
		return
	}

	// 验证文件扩展名
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".json") {
		response.BadRequest(c, "Only JSON files are allowed")
		return
	}

	// 验证文件大小（最大 10MB）
	if file.Size > 10*1024*1024 {
		response.BadRequest(c, "File size exceeds 10MB limit")
		return
	}

	// 清理文件名，防止路径遍历攻击
	filename := filepath.Base(file.Filename)
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	// 保存文件
	dst := filepath.Join(h.workflowDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		h.log.Errorw("Failed to save workflow file", "error", err)
		response.InternalError(c, "Failed to save file")
		return
	}

	// 验证 JSON 格式
	if err := h.validateWorkflowJSON(dst); err != nil {
		os.Remove(dst) // 删除无效文件
		response.BadRequest(c, fmt.Sprintf("Invalid JSON format: %v", err))
		return
	}

	h.log.Infow("Workflow uploaded", "filename", filename, "size", file.Size)

	response.Success(c, gin.H{
		"message":  "Workflow uploaded successfully",
		"filename": filename,
		"path":     filepath.Join("workflows", filename),
	})
}

// DeleteWorkflow 删除 workflow 文件
func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	filename := c.Param("filename")

	// 清理文件名，防止路径遍历攻击
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")

	filePath := filepath.Join(h.workflowDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "Workflow not found")
		return
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		h.log.Errorw("Failed to delete workflow", "error", err, "filename", filename)
		response.InternalError(c, "Failed to delete workflow")
		return
	}

	h.log.Infow("Workflow deleted", "filename", filename)

	response.Success(c, gin.H{
		"message": "Workflow deleted successfully",
	})
}

// GetWorkflow 获取 workflow 文件内容
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	filename := c.Param("filename")

	// 清理文件名，防止路径遍历攻击
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")

	filePath := filepath.Join(h.workflowDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "Workflow not found")
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.log.Errorw("Failed to read workflow", "error", err, "filename", filename)
		response.InternalError(c, "Failed to read workflow")
		return
	}

	c.Data(http.StatusOK, "application/json", content)
}

// validateWorkflowJSON 验证 JSON 格式
func (h *WorkflowHandler) validateWorkflowJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// 简单验证是否是有效的 JSON
	if len(content) == 0 {
		return fmt.Errorf("empty file")
	}

	// 检查是否以 { 或 [ 开头（JSON 对象或数组）
	trimmed := strings.TrimSpace(string(content))
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return fmt.Errorf("not a valid JSON file")
	}

	return nil
}
