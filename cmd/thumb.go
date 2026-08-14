// @author xiangqian
// @date 2026/08/11 09:46
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio/internal/ext/thumb"
	"github.com/minio/minio/internal/hash"
	"github.com/minio/pkg/v3/mimedb"
	"io"
	"path"
	"strings"
	"time"
)

var suffixes = map[string]struct{}{
	// 图片
	"jpg": {}, "jpeg": {}, "jfif": {}, "png": {}, "gif": {},
	"webp": {}, "bmp": {}, "tif": {}, "tiff": {}, "svg": {},
	"avif": {}, "arw": {}, "cr2": {}, "cr3": {}, "nef": {},
	"orf": {}, "rw2": {}, "pef": {}, "dng": {}, "raf": {},
	// 视频
	"mp4": {}, "m4v": {}, "m4a": {}, "avi": {}, "mov": {},
	"qt": {}, "wmv": {}, "asf": {}, "flv": {}, "swf": {},
	"mkv": {}, "rm": {}, "rmvb": {}, "3gp": {}, "3g2": {},
	"webm": {}, "mpeg": {}, "mpg": {}, "ts": {},
}

func genThumb(ctx context.Context, bucket string, object string, r io.Reader, written, readable thumb.Byte, duration time.Duration) (string, error) {
	// Get current object layer instance.
	objAPI := newObjectLayerFn()
	if objAPI == nil {
		return "newObjectLayerFn", errors.New("current object layer instance is nil")
	}

	// 生成缩略图
	var start = time.Now()
	var imgBuf = &bytes.Buffer{}
	err := thumb.GenImg(r, imgBuf, 200, 200, thumb.Fit)
	if err != nil {
		return "GenImg", err
	}
	duration += time.Since(start)

	// 去除原对象的扩展名，统一转换为 .webp 格式
	var index = strings.LastIndex(object, ".")
	if index != -1 && index > strings.LastIndex(object, "/") {
		var suffix = strings.ToLower(object[index+1:])
		if _, ok := suffixes[suffix]; ok {
			object = object[:index]
		}
	}
	object += ".webp"

	// 元数据
	// 标准 S3 元数据前缀："x-amz-meta-"
	var metadata = make(map[string]string)
	metadata["x-amz-meta-written"] = written.String()
	metadata["x-amz-meta-readable"] = readable.String()
	metadata["x-amz-meta-duration"] = duration.String()

	// 创建缩略图对象
	bucket = fmt.Sprintf("%s-thumb", bucket)
	err = putThumbObject(ctx, objAPI, bucket, object, imgBuf.Bytes(), metadata)
	if err != nil {
		// 存储桶不存在
		var bnf BucketNotFound
		if errors.As(err, &bnf) {
			// 创建存储桶
			err = objAPI.MakeBucket(ctx, bucket, MakeBucketOptions{})
			if err != nil {
				var be BucketExists
				if !errors.As(err, &be) {
					return "MakeBucket", err
				}
			}

			// 创建对象
			err = putThumbObject(ctx, objAPI, bucket, object, imgBuf.Bytes(), metadata)
			if err == nil {
				return "", nil
			}
		}

		return "putThumbObject", err
	}

	return "", nil
}

func putThumbObject(ctx context.Context, objAPI ObjectLayer, bucket string, object string, data []byte, metadata map[string]string) error {
	var size = int64(len(data))
	rawReader, err := hash.NewReader(ctx, bytes.NewReader(data), size, getMD5Hash(data), getSHA256Hash(data), size)
	if err != nil {
		return err
	}
	var putObjReader = NewPutObjReader(rawReader)
	_, err = objAPI.PutObject(ctx, bucket, object, putObjReader, ObjectOptions{
		UserDefined: metadata,
	})
	return err
}

// newThumbTeeBuffer 创建缩略图分流缓存器
func newThumbTeeBuffer(bucket, object string, data *PutObjReader, opts ObjectOptions) thumb.Buffer {
	//var thumbEnabled = true

	// MinIO 系统内部存储桶（存储元数据）
	if bucket == ".minio.sys" {
		return nil
	}

	// 缩略图存储桶（避免递归生成）
	if strings.HasSuffix(bucket, "-thumb") {
		return nil
	}

	// 获取对象内容类型
	// 优先使用用户通过元数据指定的 Content-Type
	// 若未指定，则根据文件扩展名自动推断
	var contentType = opts.UserDefined["content-type"]
	if contentType == "" {
		contentType = mimedb.TypeByExtension(path.Ext(object))
	}

	// 图片类型对象
	if strings.HasPrefix(contentType, "image/") {
		// 对象最大限制
		var maxSize int64 = 50 * 1024 * 1024 // 50 MB

		// 超过最大限制的对象
		if data.Size() > maxSize {
			return nil
		}

		return thumb.NewImgBuffer(bucket, object, maxSize)
	}

	// 视频类型对象
	if strings.HasPrefix(contentType, "video/") {
		return thumb.NewVidBuffer(bucket, object)
	}

	return nil
}
