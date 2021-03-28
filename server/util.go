package server

import "fmt"

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
