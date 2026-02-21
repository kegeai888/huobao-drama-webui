# 🚀 推送到 GitHub 指南

## 📊 当前状态

### ✅ 已完成

- [x] 项目整理完成
- [x] 文档重组完成
- [x] 安全检查通过
- [x] 代码提交完成
- [x] 远程仓库配置完成

### 📝 Git 配置

**当前远程仓库**:
- `origin`: https://github.com/ops120/huobao-drama.git
- `comfyui`: https://github.com/chatfire-AI/huobao-drama-comfyui.git

**最新提交**:
```
0a0ccfb - feat: 整理项目结构，添加 ComfyUI 集成功能
```

---

## 🎯 推送选项

### 选项 1: 推送到原仓库 (ops120/huobao-drama)

```bash
cd "火宝一键漫剧/huobao-drama"

# 推送到 origin (ops120/huobao-drama)
git push origin master
```

### 选项 2: 推送到新仓库 (chatfire-AI/huobao-drama-comfyui)

**步骤 1: 在 GitHub 创建新仓库**
1. 访问 https://github.com/new
2. 仓库名: `huobao-drama-comfyui`
3. 描述: `AI-powered short drama production platform with ComfyUI integration`
4. 可见性: Public
5. 不要勾选 "Initialize this repository"
6. 创建仓库

**步骤 2: 推送代码**
```bash
cd "火宝一键漫剧/huobao-drama"

# 推送到 comfyui 远程仓库
git push comfyui master:main
```

### 选项 3: 同时推送到两个仓库

```bash
cd "火宝一键漫剧/huobao-drama"

# 推送到原仓库
git push origin master

# 推送到新仓库
git push comfyui master:main
```

---

## 📋 推送前检查

### 验证文件状态

```bash
# 查看提交历史
git log --oneline -5

# 查看远程仓库
git remote -v

# 查看当前分支
git branch
```

### 验证 .gitignore

确认以下文件/目录已被忽略：
- [ ] `data/drama_generator.db`
- [ ] `data/storage/*.png`
- [ ] `data/storage/*.mp4`
- [ ] `configs/config.yaml`
- [ ] `workflows/*.json`
- [ ] `.env`

### 验证提交内容

```bash
# 查看最新提交的文件
git show --name-only
```

---

## 🔧 推送后配置

### 配置新仓库 (chatfire-AI/huobao-drama-comfyui)

**1. 仓库设置**
- Description: `AI-powered short drama production platform with ComfyUI integration`
- Website: (可选)
- Topics: `ai`, `drama`, `video-generation`, `go`, `vue3`, `comfyui`, `image-generation`

**2. 启用功能**
- [x] Issues
- [x] Discussions (可选)
- [ ] Wiki (使用 docs/ 目录)
- [ ] Projects (可选)

**3. 分支保护**
- 设置 `main` 为默认分支
- 启用分支保护规则（可选）

**4. 创建 Release**
- Tag: `v1.0.6`
- Title: `v1.0.6 - ComfyUI Integration & Project Cleanup`
- Description: 参考 README 中的 Changelog

---

## 📝 推送命令总结

### 推荐方案：推送到新仓库

```bash
# 1. 进入项目目录
cd "火宝一键漫剧/huobao-drama"

# 2. 确认远程仓库
git remote -v

# 3. 推送到新仓库（首次推送）
git push comfyui master:main

# 4. 设置上游分支（可选）
git branch --set-upstream-to=comfyui/main master
```

### 后续推送

```bash
# 推送到新仓库
git push comfyui master:main

# 或推送到原仓库
git push origin master
```

---

## ⚠️ 注意事项

### 1. 分支名称

- 原仓库使用 `master` 分支
- 新仓库建议使用 `main` 分支
- 推送时使用 `master:main` 映射

### 2. 敏感信息

确认以下内容已被 .gitignore 保护：
- 数据库文件
- 配置文件
- 环境变量
- 用户数据
- 工作流文件

### 3. 大文件

如果有大文件，考虑使用 Git LFS：
```bash
git lfs install
git lfs track "*.mp4"
git lfs track "*.png"
```

---

## 🎉 推送完成后

### 验证推送

1. 访问 GitHub 仓库页面
2. 检查文件是否完整
3. 验证 README 显示正常
4. 测试文档链接

### 更新文档

在主 README 中添加新仓库链接：
```markdown
## 🔗 相关仓库

- [huobao-drama](https://github.com/ops120/huobao-drama) - 原始仓库
- [huobao-drama-comfyui](https://github.com/chatfire-AI/huobao-drama-comfyui) - ComfyUI 集成版本
```

### 创建 Release

1. 进入 Releases 页面
2. 创建新 Release
3. 填写版本信息和更新日志
4. 发布

---

## 📧 问题反馈

如遇到问题：
- 检查网络连接
- 验证 GitHub 权限
- 查看错误信息
- 联系项目维护者

---

<div align="center">

**准备推送**: ✅  
**下一步**: 选择推送选项并执行  

Made with ❤️ by Huobao Team

</div>
