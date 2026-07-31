# 版本要求

- Go 1.23+
- libvips 8.14+

# Windows

## 下载并安装 MSYS2

https://www.msys2.org/

`msys2-x86_64-20260611.exe` 下载安装程序并安装（建议默认路径 `C:\msys64`）

## 打开 "MSYS2 MINGW64" 终端

- 更新包数据库

```shell
$ pacman -Syu
```

- 安装 pkg-config 工具，用于向 MinGW 编译器提供库的编译和链接参数

```shell
$ pacman -S mingw-w64-x86_64-pkg-config
```

- 安装 gcc

```shell
$ pacman -S mingw-w64-x86_64-gcc
$ gcc --version
```

- 安装 libvips

```shell
# 搜索 libvips
$ pacman -Ss libvips

# 安装 libvips
$ pacman -S mingw-w64-x86_64-libvips

# 安装图像处理库依赖（libjpeg, libpng, libtiff, libwebp, librsvg）
$ pacman -S mingw-w64-x86_64-{libjpeg-turbo,libpng,libtiff,libwebp,librsvg}

# 检查 libvips 是否安装成功
$ pkg-config --modversion vips

# 检查 libvips 的库文件位置
$ pkg-config --libs vips

$ pkg-config --cflags vips
# 如果输出中，头文件搜索路径 -I/mingw64/include 使用的是 Unix 风格路径（/mingw64/...），
# 而不是 Windows 绝对路径（如 C:/msys64/mingw64/...）时，Go 的 CGO 编译器可能无法正确解析。
# 解决方案：强制 pkg-config 使用完整的 Windows 路径
# $ export PKG_CONFIG_SYSROOT_DIR=C:/msys64/mingw64
```

- 安装 ffmpeg

```shell
$ pacman -S mingw-w64-x86_64-ffmpeg
$ ffmpeg -version
```

## 环境变量配置

将 `C:\msys64\mingw64\bin` 和 `C:\msys64\usr\bin` 添加到系统 `PATH`

# Linux

