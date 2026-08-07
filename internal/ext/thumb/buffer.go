// @author xiangqian
// @date 2026/08/04 22:43
package thumb

import (
	"bytes"
	"fmt"
	"github.com/minio/minio/internal/logger"
	"io"
	"time"
)

type Byte int64

func (b Byte) String() string {
	switch {
	case b >= 1<<50:
		return fmt.Sprintf("%.2fPB", float64(b)/(1<<50))
	case b >= 1<<40:
		return fmt.Sprintf("%.2fTB", float64(b)/(1<<40))
	case b >= 1<<30:
		return fmt.Sprintf("%.2fGB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.2fMB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.2fKB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

type Buffer interface {
	io.Writer
	io.Reader
	IsSuccess() bool         // 处理是否成功
	Written() Byte           // 已写入字节数
	Readable() Byte          // 可读取字节数
	Duration() time.Duration // 耗时
}

const (
	OffsetSuccess int64 = -1 // 处理成功
	OffsetError         = -2 // 处理失败
)

func NewImgBuffer(bucket, object string, maxSize int64) *ImgBuffer {
	return &ImgBuffer{
		bucket:   bucket,
		object:   object,
		buf:      &bytes.Buffer{},
		offset:   0,
		maxSize:  maxSize,
		written:  0,
		duration: time.Duration(0),
	}
}

type ImgBuffer struct {
	bucket   string        // 存储桶名称
	object   string        // 对象键名
	buf      *bytes.Buffer // 图像数据缓冲区
	offset   int64         // 视频数据偏移量
	maxSize  int64         // 图像数据最大容量
	written  int64         // 已写入字节数
	duration time.Duration // 耗时
}

func (b *ImgBuffer) Write(p []byte) (int, error) {
	if b.offset == OffsetError {
		return len(p), nil
	}

	var start = time.Now()
	var n, _ = b.buf.Write(p)
	b.written += int64(n)
	if b.written > b.maxSize {
		b.offset = OffsetError
		logger.Warning("[ext/thumb] IMG-MAX-%s bucket=%s, object=%s", Byte(b.written), b.bucket, b.object)
	}
	b.duration += time.Since(start)

	return len(p), nil
}

func (b *ImgBuffer) Read(p []byte) (int, error) {
	return b.buf.Read(p)
}

func (b *ImgBuffer) IsSuccess() bool {
	return b.offset != OffsetError
}

func (b *ImgBuffer) Written() Byte {
	return Byte(b.written)
}

func (b *ImgBuffer) Readable() Byte {
	return Byte(b.written)
}

func (b *ImgBuffer) Duration() time.Duration {
	return b.duration
}

func NewVidBuffer(bucket, object string, maxSize int64) *VidBuffer {
	return &VidBuffer{
		bucket:   bucket,
		object:   object,
		vidBuf:   &bytes.Buffer{},
		imgBuf:   &bytes.Buffer{},
		offset:   0,
		maxSize:  maxSize,
		written:  0,
		readable: 0,
		duration: time.Duration(0),
	}
}

type VidBuffer struct {
	bucket   string        // 存储桶名称
	object   string        // 对象键名
	vidBuf   *bytes.Buffer // 视频数据缓冲区
	imgBuf   *bytes.Buffer // 图像数据缓冲区
	offset   int64         // 视频数据偏移量
	maxSize  int64         // 视频数据最大容量
	written  int64         // 已写入字节数
	readable int64         // 可读取字节数
	duration time.Duration // 耗时
}

func (b *VidBuffer) Write(p []byte) (int, error) {
	if b.offset == OffsetSuccess || b.offset == OffsetError {
		return len(p), nil
	}

	var start = time.Now()
	var n, _ = b.vidBuf.Write(p)
	b.written += int64(n)
	if b.written > b.maxSize {
		err := b.capture()
		if err != nil {
			b.offset = OffsetError
			logger.Warning("[ext/thumb] VID-MAX-%s bucket=%s, object=%s, %v", Byte(b.written), b.bucket, b.object, err)
		} else {
			b.offset = OffsetSuccess
			b.readable = int64(b.imgBuf.Len())
			logger.Info("[ext/thumb] VID-%s bucket=%s, object=%s", Byte(b.written), b.bucket, b.object)
		}
	} else if b.written-b.offset >= 1*1024*1024 { // 1MB
		err := b.capture()
		if err != nil {
			b.offset = b.written
			logger.Warning("[ext/thumb] VID-%s - bucket=%s, object=%s, %v", Byte(b.written), b.bucket, b.object, err)
		} else {
			b.offset = OffsetSuccess
			b.readable = int64(b.imgBuf.Len())
			logger.Info("[ext/thumb] VID-%s bucket=%s, object=%s", Byte(b.written), b.bucket, b.object)
		}
	}
	b.duration += time.Since(start)

	return len(p), nil
}

func (b *VidBuffer) Read(p []byte) (int, error) {
	return b.imgBuf.Read(p)
}

func (b *VidBuffer) isSuccess() bool {
	return b.offset == OffsetSuccess
}

func (b *VidBuffer) IsSuccess() bool {
	if b.offset == OffsetSuccess {
		return true
	}

	if b.offset == OffsetError {
		return false
	}

	if b.offset == b.written {
		b.offset = OffsetError
		return false
	}

	var start = time.Now()
	err := b.capture()
	if err != nil {
		b.offset = OffsetError
		logger.Warning("[ext/thumb] VID-LAST-%s bucket=%s, object=%s, %v", Byte(b.written), b.bucket, b.object, err)
	} else {
		b.offset = OffsetSuccess
		b.readable = int64(b.imgBuf.Len())
		logger.Info("[ext/thumb] VID-%s bucket=%s, object=%s", Byte(b.written), b.bucket, b.object)
	}
	b.duration += time.Since(start)
	
	return b.offset == OffsetSuccess
}

func (b *VidBuffer) Written() Byte {
	return Byte(b.written)
}

func (b *VidBuffer) Readable() Byte {
	return Byte(b.readable)
}

func (b *VidBuffer) Duration() time.Duration {
	return b.duration
}

func (b *VidBuffer) capture() error {
	// 视频数据流
	var r = bytes.NewReader(b.vidBuf.Bytes())

	// 清空图像缓冲区
	b.imgBuf.Reset()

	// 从视频数据中截取指定时间点的画面
	err := Capture(r, b.imgBuf)
	if err != nil {
		b.imgBuf.Reset()
		return err
	}

	return nil
}
