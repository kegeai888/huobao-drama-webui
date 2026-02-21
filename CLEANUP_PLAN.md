# Huobao Drama 项目整理计划

## 📋 整理目标
准备项目上传到 GitHub，清理开发过程中的临时文件和测试代码，整理文档结构。

## 🗑️ 需要删除的文件

### 1. 测试文件（根目录）
- [ ] `check_ai_config.go` - AI配置检查测试
- [ ] `check_comfyui_queue.go` - ComfyUI队列检查测试
- [ ] `check_task_29.go` - 任务29检查测试
- [ ] `check_task_31.go` - 任务31检查测试
- [ ] `check_task_error.go` - 任务错误检查测试
- [ ] `clear_old_tasks.go` - 清理旧任务测试
- [ ] `test_comfyui_local.go` - ComfyUI本地测试
- [ ] `test_comfyui_simple.go` - ComfyUI简单测试
- [ ] `test_comfyui_upload.go` - ComfyUI上传测试
- [ ] `test_comfyui.go` - ComfyUI测试
- [ ] `test_exact_request.go` - 精确请求测试
- [ ] `test_full_workflow.go` - 完整工作流测试
- [ ] `test_new_task_simple.go` - 新任务简单测试
- [ ] `test_new_video_task.go` - 新视频任务测试
- [ ] `test_task_7.go` - 任务7测试
- [ ] `test_video_id7.go` - 视频ID7测试

### 2. SQL脚本文件（临时）
- [ ] `clear_old_tasks.sql` - 清理旧任务SQL
- [ ] `fix-comfyui-endpoint.sql` - 修复ComfyUI端点SQL
- [ ] `reset_task_31.sql` - 重置任务31 SQL

### 3. 临时JSON文件
- [ ] `comfyui_response.json` - ComfyUI响应示例
- [ ] `video_wan2_2_14B_i2v.json` - 视频工作流示例（已在workflows目录）

### 4. 可执行文件
- [ ] `drama-server.exe` - 编译的可执行文件（不应提交）

### 5. 临时批处理文件
- [ ] `test-comfyui-connection.cmd` - ComfyUI连接测试脚本

## 📁 需要整理的文档

### 1. ComfyUI相关文档（移动到 docs/comfyui/）
- [ ] `COMFYUI_CHANGES.md` → `docs/comfyui/CHANGES.md`
- [ ] `COMFYUI_FILES.md` → `docs/comfyui/FILES.md`
- [ ] `COMFYUI_IMAGE_UPLOAD_FIX.md` → `docs/comfyui/IMAGE_UPLOAD_FIX.md`
- [ ] `COMFYUI_README.md` → `docs/comfyui/README.md`
- [ ] `COMFYUI_UI_FIX.md` → `docs/comfyui/UI_FIX.md`
- [ ] `COMFYUI_URL_FIX.md` → `docs/comfyui/URL_FIX.md`
- [ ] `COMFYUI_WORKFLOW_GUIDE.md` → `docs/comfyui/WORKFLOW_GUIDE.md`
- [ ] `COMFYUI_WORKFLOW_PATH_FIX.md` → `docs/comfyui/WORKFLOW_PATH_FIX.md`
- [ ] `COMFYUI_WORKFLOW_UPLOAD.md` → `docs/comfyui/WORKFLOW_UPLOAD.md`

### 2. 部署相关文档（移动到 docs/deployment/）
- [ ] `DOCKER_HOST_ACCESS.md` → `docs/deployment/DOCKER_HOST_ACCESS.md`
- [ ] `MIGRATE_README.md` → `docs/deployment/MIGRATE_README.md`
- [ ] `REBUILD_INSTRUCTIONS.md` → `docs/deployment/REBUILD_INSTRUCTIONS.md`

### 3. 启动脚本文档（移动到 docs/）
- [ ] `启动脚本说明.md` → `docs/STARTUP_SCRIPTS.md`

## 📝 需要更新的文件

### 1. .gitignore
添加更多忽略规则：
```
# Test files
*_test.go
test_*.go
check_*.go

# Temporary SQL scripts
*.sql

# Temporary JSON files
*_response.json

# Compiled binaries
*.exe
drama-server
```

### 2. README.md
- 更新 ComfyUI 文档链接
- 添加文档目录结构说明
- 更新版本号到 v1.0.6

## 📂 建议的文档结构

```
huobao-drama/
├── README.md                    # 主文档（英文）
├── README-CN.md                 # 中文文档
├── README-JA.md                 # 日文文档
├── LICENSE                      # 许可证
├── docs/                        # 文档目录
│   ├── COMFYUI_INTEGRATION.md   # ComfyUI集成指南
│   ├── COMFYUI_QUICKSTART.md    # ComfyUI快速开始
│   ├── COMFYUI_SUMMARY.md       # ComfyUI总结
│   ├── DATA_MIGRATION.md        # 数据迁移指南
│   ├── STARTUP_SCRIPTS.md       # 启动脚本说明
│   ├── comfyui/                 # ComfyUI详细文档
│   │   ├── README.md
│   │   ├── CHANGES.md
│   │   ├── FILES.md
│   │   ├── WORKFLOW_GUIDE.md
│   │   └── ...
│   ├── deployment/              # 部署相关文档
│   │   ├── DOCKER_HOST_ACCESS.md
│   │   ├── MIGRATE_README.md
│   │   └── REBUILD_INSTRUCTIONS.md
│   └── comfyui-config-example.yaml
├── workflows/                   # ComfyUI工作流
└── ...
```

## ✅ 执行步骤

1. 创建文档子目录
2. 移动文档文件
3. 删除测试和临时文件
4. 更新 .gitignore
5. 更新 README 中的文档链接
6. 提交更改

## 🎯 最终目标

- ✅ 代码库干净整洁
- ✅ 文档结构清晰
- ✅ 易于新用户理解和使用
- ✅ 符合开源项目规范
