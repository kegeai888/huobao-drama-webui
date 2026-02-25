#!/bin/bash

export PATH=$PATH:/usr/local/go/bin

cd "$(dirname "$0")"

PORT=7860

# 检测并终止占用 7860 端口的进程
PID=$(lsof -ti tcp:$PORT 2>/dev/null)
if [ -n "$PID" ]; then
    echo "端口 $PORT 被进程 $PID 占用，正在终止..."
    kill -9 $PID 2>/dev/null
    echo "已终止进程 $PID"
    sleep 2
fi

echo "启动 Huobao Drama (端口 $PORT)..."
./huobao-drama
