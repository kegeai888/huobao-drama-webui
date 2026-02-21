#!/bin/bash

# ComfyUI 视频生成使用示例
# 本脚本演示如何通过 API 使用 ComfyUI 生成视频

# 配置
API_BASE="http://localhost:5678/api/v1"
STORYBOARD_ID=1
IMAGE_URL="http://localhost:5678/static/images/example.jpg"
PROMPT="A beautiful sunset scene with gentle waves"
MODEL="svd"

echo "==================================="
echo "ComfyUI 视频生成示例"
echo "==================================="
echo ""

# 1. 创建视频生成任务
echo "1. 创建视频生成任务..."
RESPONSE=$(curl -s -X POST "${API_BASE}/video-generations" \
  -H "Content-Type: application/json" \
  -d "{
    \"storyboard_id\": ${STORYBOARD_ID},
    \"image_url\": \"${IMAGE_URL}\",
    \"prompt\": \"${PROMPT}\",
    \"model\": \"${MODEL}\",
    \"duration\": 5,
    \"fps\": 8
  }")

echo "响应: ${RESPONSE}"
echo ""

# 提取任务 ID
VIDEO_GEN_ID=$(echo ${RESPONSE} | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)

if [ -z "$VIDEO_GEN_ID" ]; then
  echo "❌ 创建任务失败"
  exit 1
fi

echo "✅ 任务创建成功，ID: ${VIDEO_GEN_ID}"
echo ""

# 2. 查询任务状态
echo "2. 查询任务状态..."
MAX_ATTEMPTS=60
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
  ATTEMPT=$((ATTEMPT + 1))
  
  STATUS_RESPONSE=$(curl -s "${API_BASE}/video-generations/${VIDEO_GEN_ID}")
  STATUS=$(echo ${STATUS_RESPONSE} | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
  
  echo "尝试 ${ATTEMPT}/${MAX_ATTEMPTS} - 状态: ${STATUS}"
  
  if [ "$STATUS" = "completed" ]; then
    echo ""
    echo "✅ 视频生成完成！"
    echo ""
    echo "完整响应:"
    echo ${STATUS_RESPONSE} | python3 -m json.tool 2>/dev/null || echo ${STATUS_RESPONSE}
    exit 0
  elif [ "$STATUS" = "failed" ]; then
    echo ""
    echo "❌ 视频生成失败"
    echo ${STATUS_RESPONSE} | python3 -m json.tool 2>/dev/null || echo ${STATUS_RESPONSE}
    exit 1
  fi
  
  sleep 5
done

echo ""
echo "⏱️ 任务超时"
exit 1
