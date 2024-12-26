package filemanager

import "time"

// 文件管理接口
type Filemanager interface {
	Upload(filePath string) (*File, error)
	Download(fileName string) (*File, error)
	Delete(fileName string) error
	List() ([]File, error)
	Get(fileName string) (*File, error)
	Rename(oldName, newName string) error
	Copy(sourceFileName, destFileName string) error
	Move(sourceFileName, destFileName string) error
}

// 文件信息
type File struct {
	Name         string
	Size         int64
	ContentType  string    // Mime类型
	ModTime      time.Time // 最后修改时间
	ETag         string    // 文件的唯一标识符（通常是 MD5 或类似的散列值）
	Owner        string    //
	StorageClass string    // 文件的存储类型（如 "STANDARD", "REDUCED_REDUNDANCY"）
}
