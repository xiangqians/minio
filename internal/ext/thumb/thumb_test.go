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

// 测试生成图片缩略图
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

// 测试生成视频缩略图
func TestVidGen(t *testing.T) {
	// 视频文件名
	var vidName = "D:\\tmp\\minio\\tmp\\test-26.mp4"
	vidName = "D:\\tmp\\minio\\tmp\\test-174.mp4"

	// 打开文件
	file, err := os.Open(vidName)
	if err != nil {
		fmt.Sprintf("open file failed: %v", err)
		return
	}
	defer file.Close()

	// 创建字节缓冲区，用于存储将要传递给 ffmpeg 的视频数据
	var b = make([]byte, 1024*1024)
	var r = bytes.NewBuffer(nil)

	for {
		n, err := file.Read(b)
		if n > 0 {
			r.Write(b[:n])

			var seekSecond = 1

			// 创建字节缓冲区，用于接收 ffmpeg 输出的图片数据
			var w = bytes.NewBuffer(nil)

			err = VidGen(bytes.NewReader(r.Bytes()), w, seekSecond, 200, 200, Fill)
			if err != nil {
				continue
			}

			// 直接将字节数据写入文件
			if w.Len() > 0 {
				var index = strings.LastIndex(vidName, ".")
				var imgName = fmt.Sprintf("%s_%ds.jpeg", vidName[:index], seekSecond)
				err = os.WriteFile(imgName, w.Bytes(), 0644)
				if err != nil {
					log.Fatalf("write file failed: %v", err)
				}
			}
			log.Printf("read %s\n", formatBytes(int64(w.Len())))
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

func formatBytes(b int64) string {
	switch {
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
