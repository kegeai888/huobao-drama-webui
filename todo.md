# 火宝短剧 WebUI 二次开发记录

> 基于 huobao-drama-webui 项目，记录所有定制化修改、部署经验和注意事项。
> 仓库地址：https://github.com/kegeai888/huobao-drama-webui

---

## 一、环境部署

### 系统环境
- OS: Ubuntu 22.04.3 LTS (x86_64)
- Node.js: v24.11.0 / npm: 11.6.1
- FFmpeg: 4.4.2（系统自带）
- Go: 1.23.6（手动安装至 `/usr/local/go`）

### 部署步骤

1. **安装 Go 1.23**
   ```bash
   curl -fsSL https://go.dev/dl/go1.23.6.linux-amd64.tar.gz -o /tmp/go1.23.6.linux-amd64.tar.gz
   tar -C /usr/local -xzf /tmp/go1.23.6.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **复制配置文件**
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   ```

3. **安装前端依赖并构建**
   ```bash
   cd web && npm install && npm run build && cd ..
   ```

4. **下载 Go 依赖并编译**
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   go mod download
   go build -o huobao-drama .
   ```

5. **启动服务**
   ```bash
   ./start_app.sh
   ```

### 注意事项
- 每次新 shell 需要 `export PATH=$PATH:/usr/local/go/bin`，或写入 `/etc/profile`
- 前端修改后必须重新 `npm run build`，后端才能提供最新静态文件
- 数据库文件：`data/drama_generator.db`，存储路径：`data/storage/`

---

## 二、配置修改

### 访问端口改为 7860
- 文件：`configs/config.yaml`
- 修改：`server.port: 5678` → `server.port: 7860`
- 访问地址：`http://localhost:7860`

---

## 三、启动脚本

### start_app.sh
- 路径：`/root/huobao-drama-webui/start_app.sh`
- 功能：
  1. 检测 7860 端口是否被占用
  2. 若占用则 `kill -9` 强制终止旧进程（无需确认）
  3. `sleep 2` 等待端口释放
  4. 启动 `./huobao-drama` 主程序

```bash
#!/bin/bash
export PATH=$PATH:/usr/local/go/bin
cd "$(dirname "$0")"
PORT=7860
PID=$(lsof -ti tcp:$PORT 2>/dev/null)
if [ -n "$PID" ]; then
    kill -9 $PID 2>/dev/null
    sleep 2
fi
./huobao-drama
```

---

## 四、UI 定制修改

### 4.1 页面标题和品牌信息

| 位置 | 原内容 | 修改后 |
|------|--------|--------|
| AppHeader.vue / AppLayout.vue | `ComfyUI 贡献者: 你们喜爱的老王` | `构建by科哥` |
| AppHeader.vue / AppLayout.vue | B站链接 `https://space.bilibili.com/97727630` | `https://space.bilibili.com/54364653` |
| web/index.html | favicon.ico | favicon.png |

**涉及文件：**
- `web/src/components/common/AppHeader.vue`
- `web/src/components/common/AppLayout.vue`
- `web/index.html`

### 4.2 导航菜单文字

| 位置 | 原内容 | 修改后 |
|------|--------|--------|
| zh-CN.ts `settings.aiConfig` | `AI配置` | `API服务配置` |
| zh-CN.ts `aiConfig.title` | `AI 服务配置` | `API服务配置` |

**涉及文件：**
- `web/src/locales/zh-CN.ts`

---

## 五、API服务配置页面定制

### 5.1 注册链接样式

- 文件：`web/src/components/common/AIConfigDialog.vue`
- 原链接：`https://api.chatfire.site/login?inviteCode=C4453345`
- 新链接：`https://ai.kegeai.top/register?aff=78Gs`
- 样式：颜色改为橘色 `#f97316`，字体大小 `14px`

### 5.2 默认厂商和 Base URL

- 原默认厂商：`chatfire`，base_url：`https://api.chatfire.site/v1`
- 新默认厂商：`kegeai`（API国际站），base_url：`https://api.kegeai.top`
- 涉及文件：
  - `web/src/components/common/AIConfigDialog.vue`
  - `web/src/views/settings/AIConfig.vue`

### 5.3 厂商列表新增"API国际站"

三个 tab 均在首位新增 `kegeai`（API国际站）厂商。

### 5.4 各 Tab 默认模型

| Tab | 默认模型 |
|-----|---------|
| 文本生成 | `deepseek-v3.2` |
| 图片生成 | `gemini-3-pro-image-preview` |
| 视频生成 | `sora-2-all` |

### 5.5 API国际站模型列表

**文本生成（14个）：**
`gemini-3.1-pro-preview`, `gemini-3-flash-preview`, `gemini-2.5-flash`, `claude-sonnet-4-6`, `gpt-5.2`, `gpt-5.1`, `deepseek-v3.2`, `deepseek-v3.2-fast`, `deepseek-v3.1`, `deepseek-v3.1-fast`, `doubao-seed-1-8-251228`, `doubao-seed-1-6-250615`, `glm-4.7`, `kimi-k2.5`

**图片生成（6个）：**
`gemini-2.5-flash-image-preview`, `gemini-3-pro-image-preview`, `doubao-seedream-4-0-250828`, `doubao-seedream-4-5-251128`, `grok-3-image`, `grok-4-image`

**视频生成（10个）：**
`sora-2-all`, `sora-2-vip-all`, `sora-2-pro-all`, `veo_3_1-fast`, `veo3.1`, `grok-video-3-15s`, `grok-video-3-10s`, `grok-video-3`, `wan2.6-i2v-flash`, `kling-video`

---

## 六、经验与注意事项

1. **前端修改必须重新构建**：所有 `.vue` / `.ts` 修改后需在 `web/` 目录执行 `npm run build`，否则后端服务的静态文件不会更新。

2. **两个文件需同步修改**：`AIConfigDialog.vue`（弹窗组件）和 `AIConfig.vue`（设置页面）都有独立的 `providerConfigs` 和 `showCreateDialog`，修改厂商/模型/默认值时两个文件必须同步。

3. **Go 路径问题**：系统未将 `/usr/local/go/bin` 加入 PATH，每次执行 Go 命令需手动 export，建议写入 `/etc/profile` 永久生效：
   ```bash
   echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
   source /etc/profile
   ```

4. **Edit 工具使用前必须先 Read**：使用 Edit 工具修改文件前，必须先用 Read 工具读取该文件，否则会报错。

5. **favicon 路径**：Vite 构建时 `public/` 目录下的文件会直接复制到 `dist/` 根目录，`index.html` 中引用路径为 `/favicon.png`（不含 `/public/`）。
