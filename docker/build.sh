#!/bin/bash

# 清理所有未使用的构建缓存
#docker builder prune -a

# beta
docker build \
  -f docker/Dockerfile \
  -t minio:2026-08-beta .

# release
