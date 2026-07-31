// @author xiangqian
// @date 2026/07/29 21:17
package thumb

import (
	"bytes"
	"fmt"
	"github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"testing"
)

// IDEA
// Environment: CGO_ENABLED=1;PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%

func TestC(t *testing.T) {
	c()
}

func TestImgGen(t *testing.T) {
	// 初始化缩略图模块
	err := Startup()
	if err != nil {
		log.Fatalf("vips startup failed: %v", err)
	}
	// 关闭缩略图模块，释放所有占用的资源
	defer Shutdown()

	var dir = "D:\\tmp\\minio\\tmp"
	var src = fmt.Sprintf("%s\\%s", dir, "test.jpg")
	var dst = fmt.Sprintf("%s\\%s", dir, "test-thumb.jpg")

	// 源文件
	r, err := os.Open(src)
	if err != nil {
		log.Fatalf("src open failed: %v", err)
	}
	defer r.Close()

	// 目标文件
	w, err := os.Create(dst)
	if err != nil {
		log.Fatalf("dst open failed: %v", err)
	}
	defer w.Close()

	// 生成缩略图
	err = Gen(r, w, 200, 200, Fill)
	if err != nil {
		log.Fatalf("thumb gen failed: %v", err)
	}
	log.Printf("缩略图生成成功：%s\n", dst)
}

func TestVidGen(t *testing.T) {
	// 视频文件名
	var vidName = "D:\\tmp\\minio\\tmp\\test26.mp4"

	// 要截取的时间点（单位：秒）
	var seekSecond = 2
	// 截取图片文件名
	var imgName = fmt.Sprintf("D:\\tmp\\minio\\tmp\\test26_%d.jpeg", seekSecond)

	// 创建字节缓冲区，用于接收 ffmpeg 输出的图片数据
	buf := bytes.NewBuffer(nil)

	// 链式调用 ffmpeg 命令
	err := ffmpeg_go.Input(vidName,
		ffmpeg_go.KwArgs{"ss": seekSecond - 1}). // 粗略快速跳转到目标时间前 1 秒
		Output("pipe:", ffmpeg_go.KwArgs{
			"vframes": 1,        // 只输出 1 帧
			"format":  "image2", // 输出格式为图片
			"vcodec":  "mjpeg",  // 编码为 JPEG
			"ss":      1,        // 在上次跳转的基础上，再精确前进 1 秒
		}).
		WithOutput(buf).
		Run()
	if err != nil {
		log.Fatalf("ffmpeg run failed: %v", err)
	}

	// 直接将字节数据写入文件
	err = os.WriteFile(imgName, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("write file failed: %v", err)
	}
}
