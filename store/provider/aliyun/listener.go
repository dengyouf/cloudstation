package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/schollz/progressbar/v3"
)

func NewListener() oss.ProgressListener {
	return &listener{}
}

type listener struct {
	bar *progressbar.ProgressBar
}

func (l *listener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		l.bar = progressbar.DefaultBytes(
			event.TotalBytes,
			"文件上传中",
		)
	case oss.TransferDataEvent:
		l.bar.Add64(event.RwBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\n上传完成\n")
	case oss.TransferFailedEvent:
		fmt.Printf("\n上传失败\n")
	default:
	}
	// fmt.Println(event.EventType, event.TotalBytes, event.RwBytes)
}
