// @author xiangqian
// @date 2026/07/29 22:44
package thumb

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
)

// Mode 缩略图模式
type Mode int

const (
	Fit  Mode = iota // 适配模式：保持宽高比，完整显示图片，可能留白
	Fill             // 填充模式：保持宽高比，完全覆盖目标区域，可能裁剪
)

func (mode Mode) String() string {
	switch mode {
	case Fit:
		return "fit"
	case Fill:
		return "fill"
	default:
		return fmt.Sprintf("unknown(%d)", mode)
	}
}

// Startup 初始化缩略图模块
func Startup() error {
	// 启动 libvips（在整个程序生命周期中只需执行一次）
	//config := &vips.Config{
	//	MaxCacheSize:  1024 * 1024 * 1024, // 1GB 缓存
	//	MaxCacheFiles: 100,
	//})
	return vips.Startup(nil)
}

// Shutdown 关闭缩略图模块，释放所有占用的资源
func Shutdown() {
	// 关闭 libvips 库连接、清理内存池、停止后台协程
	vips.Shutdown()
}

// Gen 生成图片缩略图
// r      图片数据流
// w      缩略图数据流
// width  宽度
// height 高度
// mode   模式
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

// VidGen 生成视频缩略图
// r          视频数据流
// w          缩略图数据流
// seekSecond 要截取视频的时间点（单位：秒）
// width      宽度
// height     高度
// mode       模式
func VidGen(r io.Reader, w io.Writer, seekSecond, width, height int, mode Mode) error {
	// 链式调用 ffmpeg 命令从视频数据中截取指定时间点的画面
	return ffmpeg.Input("pipe:", // 输入源为标准输入（stdin）
		ffmpeg.KwArgs{}).
		WithInput(r). // 将缓冲区作为输入数据写入 ffmpeg 的标准输入（stdin）
		Output("pipe:", // 输出到标准输出（stdout）
			ffmpeg.KwArgs{
				"vframes": 1,          // 只输出 1 帧
				"format":  "image2",   // 输出格式为图片
				"vcodec":  "mjpeg",    // 编码为 JPEG
				"ss":      seekSecond, // 截取指定时间点
			},
		).
		WithOutput(w).
		Run()
}

func c() {
	// 调用 C 的 malloc 和 free 函数
	ptr := C.malloc(100)
	defer C.free(ptr)
	log.Println("CGO works! Memory allocated:", ptr)
}
