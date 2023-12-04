package filesystem

import (
	"github.com/minio/minio-go/v7"
)

type Filesystem interface {
	List() ([]string, error)
	Read(path string) ([]byte, error)
	Create(path string, data []byte) error
	Remove(path string) error
}

type FS struct {
	Filesystem  Filesystem
	Driver      string
	Dir         string
	Client      *minio.Client
	StorageName string
}

const (
	Local = "local"
	Cloud = "cloud"
)
