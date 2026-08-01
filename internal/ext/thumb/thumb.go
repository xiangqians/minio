// @author xiangqian
// @date 2026/07/29 22:44
package thumb

// #include <stdlib.h>
import "C"

import (
	"bytes"
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

type Byte int64

func (b Byte) String() string {
	switch {
	case b >= 1<<50:
		return fmt.Sprintf("%.2f PB", float64(b)/(1<<50))
	case b >= 1<<40:
		return fmt.Sprintf("%.2f TB", float64(b)/(1<<40))
	case b >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
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

// ImgGen 生成图片缩略图（WebP 格式）
// r      图片数据流
// w      缩略图数据流
// width  宽度
// height 高度
// mode   模式
func ImgGen(r io.Reader, w io.Writer, width, height int, mode Mode) error {
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
		crop = vips.InterestingNone // 不裁剪
		size = vips.SizeDown        // 等比例缩小，不超过指定尺寸
	case Fill:
		crop = vips.InterestingCentre // 居中裁剪
		size = vips.SizeBoth          // 等比例缩放，取宽高中较大的比例，确保完全覆盖指定尺寸
	default:
		return fmt.Errorf("unsupported mode: %s", mode)
	}
	thumb, err := vips.NewThumbnailWithSizeFromBuffer(data, width, height, crop, size)
	if err != nil {
		return fmt.Errorf("create thumbnail failed: %v", err)
	}
	defer thumb.Close()

	// 导出为 WebP 格式
	data, _, err = thumb.ExportWebp(&vips.WebpExportParams{
		StripMetadata:   true,  // 删除元数据（Exif等）
		Quality:         80,    // 75-85 最佳平衡
		Lossless:        false, // 必须显式设置为 false（有损压缩）
		NearLossless:    false, // 近无损模式，通常不需要
		ReductionEffort: 6,     // 最大压缩努力（0-6），6 最慢但文件最小
		IccProfile:      "",    // 删除 ICC 配置文件
		MinSize:         true,  // 优先减小文件大小
		MinKeyFrames:    0,     // 最小关键帧间隔
		MaxKeyFrames:    0,     // 最大关键帧间隔
	})
	if err != nil {
		return fmt.Errorf("export webp failed: %v", err)
	}

	// 写入数据
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("write data failed: %v", err)
	}
	return nil
}

// VidGen 生成视频缩略图（WebP 格式）
// r      视频数据流
// w      缩略图数据流
// width  宽度
// height 高度
// mode   模式
func VidGen(r io.Reader, w io.Writer, width, height int, mode Mode) error {
	// 快速跳转到指定时间点截取
	var seekSecond = 1

	// 创建字节缓冲区，用于接收 ffmpeg 输出的图片数据
	var buf bytes.Buffer

	// 链式调用 ffmpeg 命令从视频数据中截取指定时间点的画面
	err := ffmpeg.Input("pipe:", // 输入源为标准输入（stdin）
		ffmpeg.KwArgs{}).
		WithInput(r). // 将缓冲区作为输入数据写入 ffmpeg 的标准输入（stdin）
		Output("pipe:", // 输出到标准输出（stdout）
			ffmpeg.KwArgs{
				"vframes":           1,          // 只输出 1 帧
				"format":            "image2",   // 输出格式为图片
				"vcodec":            "webp",     // 编码为 WebP，文件标准后缀为 .webp
				"ss":                seekSecond, // 快速跳转到指定时间点截取
				"compression_level": 6,          // 0-6，6 压缩率最高
				"quality":           80,         // 0-100，推荐 75-85
				"lossless":          0,          // 0=有损，1=无损
			},
		).
		WithOutput(&buf).
		Run()
	if err != nil {
		return fmt.Errorf("video frame capture failed: %w", err)
	}
	if buf.Len() == 0 {
		return fmt.Errorf("video frame capture failed: output is empty")
	}

	// 生成图片缩略图
	return ImgGen(&buf, w, width, height, mode)
}

func c() {
	// 调用 C 的 malloc 和 free 函数
	ptr := C.malloc(100)
	defer C.free(ptr)
	log.Println("CGO works! Memory allocated:", ptr)
}
