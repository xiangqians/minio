// @author xiangqian
// @date 2026/07/29 21:17
package thumb

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// IDEA
// Environment: CGO_ENABLED=1;PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%

// 测试 C
func TestC(t *testing.T) {
	c()
}

// 测试生成图片缩略图（WebP 格式）
func TestGenImg(t *testing.T) {
	// 初始化缩略图模块
	err := Startup()
	if err != nil {
		log.Fatalf("vips startup failed: %v", err)
	}
	// 关闭缩略图模块，释放所有占用的资源
	defer Shutdown()

	var src = "D:\\tmp\\minio\\tmp\\test.jpg"
	var index = strings.LastIndex(src, ".")
	var dst = fmt.Sprintf("%s-thumb.webp", src[:index])

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
	err = GenImg(r, w, 200, 200, Fit)
	if err != nil {
		log.Fatalf("thumb gen failed: %v", err)
	}
	log.Printf("Thumbnail created: %s\n", dst)
}

// 测试从视频数据中截取指定时间点的画面（WebP 格式）
func TestCapture(t *testing.T) {
	// 视频文件名
	var vidName = "D:\\tmp\\minio\\tmp\\test.mp4"
	vidName = "D:\\tmp\\minio\\tmp\\test-faststart.mp4"

	// 打开文件
	vidFile, err := os.Open(vidName)
	if err != nil {
		fmt.Sprintf("open file failed: %v", err)
		return
	}
	defer vidFile.Close()

	// 创建字节缓冲区，用于存储将要传递给 ffmpeg 的视频数据
	var b = make([]byte, 1*1024*1024)

	var buf = NewVidBuffer("", "")

	for {
		n, err := vidFile.Read(b)
		if n > 0 {
			buf.Write(b[:n])
			if buf.isSuccess() {
				saveWebp(vidName, buf)
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read file failed: %v\n", err)
			return
		}
	}

	if buf.IsSuccess() {
		saveWebp(vidName, buf)
		return
	}
}

func saveWebp(name string, r io.Reader) {
	var index = strings.LastIndex(name, ".")
	name = fmt.Sprintf("%s.webp", name[:index])
	file, err := os.Create(name)
	if err != nil {
		log.Printf("create WebP file failed: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		log.Printf("write WebP file failed: %v\n", err)
		return
	}
	log.Printf("WebP saved: %s\n", name)
}
