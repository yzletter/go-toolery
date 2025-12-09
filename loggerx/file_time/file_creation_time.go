// file_creation_time.go
package file_time

import "time"

// FileCreationTime 返回给定路径文件的创建时间。
// 不同平台由各自的实现完成，Linux 默认返回不支持错误。
func FileCreationTime(path string) (time.Time, error) {
	return fileCreationTime(path)
}
