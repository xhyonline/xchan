package server

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FormatFileSizeAndUnit 格式化
func FormatFileSizeAndUnit(fileSize int64) (string, string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1)), "B"
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024)), "KB"
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024)), "MB"
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024)), "GB"
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024)), "TB"
	} else {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024*1024)), "EB"
	}
}

// GetCurrentPath 获取当前路径
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, `\`)
	}
	if i < 0 {
		return "", errors.New(`error: Can't find / or \.`)
	}
	return path[0 : i+1], nil
}
