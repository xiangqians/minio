# [编译阶段]

# 基于 Alpine Linux 的官方镜像，包含完整的 Go 编译环境
FROM golang:1.26.4-alpine AS builder

# 更新软件包
RUN apk update
# 安装 gcc, g++, make, binutils - C/C++ 编译工具链
RUN apk add --no-cache build-base
# 安装 查询库的编译/链接参数
RUN apk add --no-cache pkgconfig
# 安装 vips
RUN apk add --no-cache vips-dev
# 安装图像处理库依赖
RUN apk add --no-cache libjpeg-turbo-dev libpng-dev tiff-dev libwebp-dev librsvg-dev
# 安装 ffmpeg
RUN apk add --no-cache ffmpeg-dev
# 安装 UPX 压缩工具
RUN apk add --no-cache upx

# 设置工作目录
WORKDIR /app

# 复制源码
# 将 minio 目录下的文件复制到 /app 目录下
COPY minio/ ./

# 设置代理
RUN go env -w GOPROXY=https://goproxy.cn,direct
# 下载依赖
RUN go mod download
# 启用 CGO
ENV CGO_ENABLED=1
# 构建
RUN go build -ldflags="-s -w" -o build/minio

# 压缩可执行文件
# https://upx.github.io
#RUN upx --best --backup build/minio
#RUN upx --brute --backup build/minio


# [运行阶段]

# 基于 Alpine Linux 发行版的最小化镜像，不包含非必要组件
FROM alpine:3.20 AS runner

USER root

# 更新软件包
RUN apk update
# 安装 vips
RUN apk add --no-cache vips
# 安装图像处理库依赖
RUN apk add --no-cache libjpeg-turbo libpng tiff libwebp librsvg
# 安装 ffmpeg
RUN apk add --no-cache ffmpeg
# 安装 SSL/TLS 根证书，使容器内应用能够通过 HTTPS 访问外部服务
RUN apk add --no-cache ca-certificates

# 设置默认时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从编译阶段复制文件到 /app 目录下
COPY --from=builder /app/build/minio ./

# 设置 MinIO 管理员账号/密码
ENV MINIO_ROOT_USER=minioadmin
ENV MINIO_ROOT_PASSWORD=minioadmin

# 声明容器运行时监听的端口
# MinIO API 端口
EXPOSE 9000
# MinIO Console 端口
EXPOSE 9001

# 启动命令
ENTRYPOINT ["/app/minio"]
CMD ["server", "/app/data", "--address", ":9000", "--console-address", ":9001"]
