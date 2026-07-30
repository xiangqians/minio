// @author xiangqian
// @date 2026/07/29 21:17
package thumb

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// IDEA
// Environment: CGO_ENABLED=1;PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%

func TestC(t *testing.T) {
	c()
}

func TestGen(t *testing.T) {
	err := Startup()
	if err != nil {
		log.Fatalf("vips startup failed: %v", err)
	}
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
