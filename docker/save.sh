#!/bin/bash

# 操作系统
os="linux"

# CPU 架构
arch=$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/; s/armv7l/armv7/')

docker save -o minio-2026-08-beta-${os}-${arch}.tar minio:2026-08-beta
