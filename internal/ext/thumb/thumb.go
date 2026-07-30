// @author xiangqian
// @date 2026/07/29 22:44
package thumb

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"io"
	"log"
)

// 版本要求：Go 1.23+，libvips 8.14+

// Windows
// 1、下载并安装 MSYS2
// https://www.msys2.org/
// msys2-x86_64-20260611.exe
// 下载安装程序并安装（建议默认路径 C:\msys64）
// 2、打开 "MSYS2 MINGW64" 终端
// 更新包数据库
// $ pacman -Syu
// 安装 pkg-config 工具，用于向 MinGW 编译器提供库的编译和链接参数
// $ pacman -S mingw-w64-x86_64-pkg-config
// 安装 gcc
// $ pacman -S mingw-w64-x86_64-gcc
// 搜索并安装 libvips
// $ pacman -Ss libvips
// $ pacman -S mingw-w64-x86_64-libvips
// 安装图像处理库依赖（libjpeg, libpng, libtiff, libwebp, librsvg）
// $ pacman -S mingw-w64-x86_64-{libjpeg-turbo,libpng,libtiff,libwebp,librsvg}
// 3、验证安装
// 检查 gcc 是否可用
// $ gcc --version
// 检查 libvips 是否安装成功
// $ pkg-config --modversion vips
// 检查 libvips 的库文件位置
// $ pkg-config --libs vips
// $ pkg-config --cflags vips
// 如果输出中，头文件搜索路径 -I/mingw64/include 使用的是 Unix 风格路径（/mingw64/...），
// 而不是 Windows 绝对路径（如 C:/msys64/mingw64/...）时，Go 的 CGO 编译器可能无法正确解析。
// 解决方案：强制 pkg-config 使用完整的 Windows 路径
// $ export PKG_CONFIG_SYSROOT_DIR=C:/msys64/mingw64
// 4、将 C:\msys64\mingw64\bin 和 C:\msys64\usr\bin 添加到系统 PATH

// Linux
//

// Mode 定义缩略图生成模式
type Mode int

const (
	Fit  Mode = iota // 适配模式：保持宽高比，完整显示图片，可能留白
	Fill             // 填充模式：保持宽高比，完全覆盖目标区域，可能裁剪
)

func (mode Mode) String() string {
	switch mode {
	case Fit:
		return fmt.Sprintf("fit(%d)", mode)
	case Fill:
		return fmt.Sprintf("fill(%d)", mode)
	default:
		return fmt.Sprintf("unknown(%d)", mode)
	}
}

func Startup() error {
	// 启动 libvips（在整个程序生命周期中只需执行一次）
	//config := &vips.Config{
	//	MaxCacheSize:  1024 * 1024 * 1024, // 1GB 缓存
	//	MaxCacheFiles: 100,
	//})
	return vips.Startup(nil)
}

func Shutdown() {
	vips.Shutdown()
}

// Gen 生成缩略图
// width  宽度
// height 高度
func Gen(r io.Reader, w io.Writer, width, height int, mode Mode) error {
	// 读取原始图片
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read file failed: %v", err)
	}

	// 生成缩略图
	var crop vips.Interesting
	var size vips.Size
	switch mode {
	case Fit:
		crop = vips.InterestingNone
		size = vips.SizeDown
	case Fill:
		// 裁剪策略 (InterestingCentre 表示居中裁剪)
		crop = vips.InterestingCentre
		size = vips.SizeBoth
	default:
		return fmt.Errorf("unsupported mode: %s", mode)
	}
	thumb, err := vips.NewThumbnailWithSizeFromBuffer(data, width, height, crop, size)
	if err != nil {
		return fmt.Errorf("create thumbnail failed: %v", err)
	}
	defer thumb.Close()

	// 导出为 JPEG 格式
	data, _, err = thumb.ExportJpeg(&vips.JpegExportParams{
		Quality:   85,   // JPEG 图片质量（1-100）
		Interlace: true, // 是否隔行扫描
	})
	if err != nil {
		return fmt.Errorf("export jpeg failed: %v", err)
	}

	// 写入数据
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write data failed: %v", err)
	}
	return nil
}

func c() {
	// 调用 C 的 malloc 和 free 函数
	ptr := C.malloc(100)
	defer C.free(ptr)
	log.Println("CGO works! Memory allocated:", ptr)
}
