package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type localFilesystem struct {
	dir string
}

func CreateLocalFilesystem(dir string) Filesystem {
	if !strings.HasPrefix(dir, "/") && !strings.HasPrefix(dir, "./") {
		dir = "/" + dir
	}
	s := &localFilesystem{
		dir: dir,
	}
	if err := s.validateDir(dir); err != nil {
		panic(err)
	}
	return s
}

func (s *localFilesystem) List() ([]string, error) {
	result := make([]string, 0)
	err := filepath.Walk(
		s.dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			result = append(result, path)
			return nil
		},
	)
	return result, err
}

func (s *localFilesystem) Read(path string) ([]byte, error) {
	return os.ReadFile(s.createPath(path))
}

func (s *localFilesystem) Create(path string, data []byte) error {
	path = s.createPath(path)
	if err := s.validateDir(path[:strings.LastIndex(path, "/")]); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = f.Chmod(0600); err != nil {
		return err
	}
	return f.Close()
}

func (s *localFilesystem) Remove(path string) error {
	return os.Remove(s.createPath(path))
}

func (s *localFilesystem) createPath(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	if strings.HasPrefix(path, "./") {
		path = strings.TrimPrefix(path, "./")
	}
	return fmt.Sprintf("%s/%s", s.dir, path)
}

func (s *localFilesystem) validateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}
