#!/bin/bash

# 容器名称
NAME="minio"

# 端口映射（宿主机:容器）
PORT_API="9000:9000"     # MinIO API 端口
PORT_CONSOLE="9001:9001" # MinIO Web 控制台端口

# 数据持久化路径（宿主机:容器）
VOLUME_DATA="/opt/minio/data:/app/data"

# 管理员账号/密码
ROOT_USER="minioadmin"
ROOT_PASS="minioadmin"

# 日志配置
LOG_MAX_SIZE="10m" # 单个日志文件最大 10MB
LOG_MAX_FILE="1"   # 最多保留 1 个日志文件

# 启动容器 — 后台运行模式
docker run -d \
  --name ${NAME} \
  --restart always \
  -p ${PORT_API} \
  -p ${PORT_CONSOLE} \
  -v ${VOLUME_DATA} \
  -e MINIO_ROOT_USER=${ROOT_USER} \
  -e MINIO_ROOT_PASSWORD=${ROOT_PASS} \
  --log-opt max-size=${LOG_MAX_SIZE} \
  --log-opt max-file=${LOG_MAX_FILE} \
  minio:2026-08-beta
