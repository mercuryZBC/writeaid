#!/bin/bash

# 定义日志目录
LOG_DIR="./log"
mkdir -p "$LOG_DIR"

# 启动服务并将输出记录到日志文件
nohup go run ./app/api/cmd > "$LOG_DIR/api.log" 2>&1 &
echo "API service started with PID $!"

nohup go run ./app/knowledgeBaseService/cmd > "$LOG_DIR/knowledge_base_service.log" 2>&1 &
echo "Knowledge Base Service started with PID $!"

nohup go run ./app/notificationService/cmd > "$LOG_DIR/notification_service.log" 2>&1 &
echo "Notification Service started with PID $!"

nohup go run ./app/userService/cmd > "$LOG_DIR/user_service.log" 2>&1 &
echo "User Service started with PID $!"