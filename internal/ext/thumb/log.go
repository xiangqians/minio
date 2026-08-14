// @author xiangqian
// @date 2026/08/10 21:02
package thumb

import (
	"github.com/minio/minio/internal/logger"
	"time"
)

// cmd/object-handlers.go
// func (api objectAPIHandlers) PutObjectHandler(w http.ResponseWriter, r *http.Request) {}

// cmd/object-multipart-handlers.go
// func (api objectAPIHandlers) NewMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {}
// func (api objectAPIHandlers) PutObjectPartHandler(w http.ResponseWriter, r *http.Request) {}
// func (api objectAPIHandlers) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {}
// func (api objectAPIHandlers) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {}

func Info(format string, a ...any) {
	Log("INFO", format, a...)
}

func Warn(format string, a ...any) {
	Log("WARN", format, a...)
}

func Error(format string, a ...any) {
	Log("ERROR", format, a...)
}

func Log(level string, format string, a ...any) {
	format = "%s [ext/thumb] " + format

	// len=0, cap=1+len(a)
	na := make([]any, 0, 1+len(a))
	na = append(na, time.Now().Format("2006/01/02 15:04:05.000"))
	na = append(na, a...)

	if level == "INFO" {
		logger.Info(format, na...)
	} else if level == "WARN" {
		logger.Warning(format, na...)
	} else if level == "ERROR" {
		logger.Error(format, na...)
	} else {
		logger.Info(format, na...)
	}
}
