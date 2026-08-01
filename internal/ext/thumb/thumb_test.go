// @author xiangqian
// @date 2026/07/29 21:17
package thumb

import (
	"bytes"
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

// 测试生成视频缩略图（WebP 格式）
func TestGenVid(t *testing.T) {
	// 视频文件名
	var vidName = "D:\\tmp\\minio\\tmp\\test.mp4"

	// 打开文件
	vidFile, err := os.Open(vidName)
	if err != nil {
		fmt.Sprintf("open file failed: %v", err)
		return
	}
	defer vidFile.Close()

	// 创建字节缓冲区，用于存储将要传递给 ffmpeg 的视频数据
	var b = make([]byte, 1024*1024)
	var r = bytes.NewBuffer(nil)

	for {
		n, err := vidFile.Read(b)
		if n > 0 {
			r.Write(b[:n])

			// 创建字节缓冲区，用于接收缩略图数据流
			var w = bytes.NewBuffer(nil)

			// 生成视频缩略图
			err = GenVid(bytes.NewReader(r.Bytes()), w, 200, 200, Fit)
			if err != nil {
				continue
			}

			// 直接将字节数据写入文件
			var index = strings.LastIndex(vidName, ".")
			var imgName = fmt.Sprintf("%s.webp", vidName[:index])
			err = os.WriteFile(imgName, w.Bytes(), 0644)
			if err != nil {
				log.Fatalf("write file failed: %v", err)
			}
			log.Printf("read %s\n", Byte(int64(w.Len())))
			break
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read file failed: %v\n", err)
			return
		}
	}
}
