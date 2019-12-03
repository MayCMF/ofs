package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MyCMF/ofs"
)

// FileSystem file system storage
type FileSystem struct {
	Base string
}

// New initialize FileSystem storage
func New(base string) *FileSystem {
	absbase, err := filepath.Abs(base)
	if err != nil {
		fmt.Println("FileSystem storage's directory haven't been initialized")
	}
	return &FileSystem{Base: absbase}
}

// GetFullPath get full path from absolute/relative path
func (fileSystem FileSystem) GetFullPath(path string) string {
	fullpath := path
	if !strings.HasPrefix(path, fileSystem.Base) {
		fullpath, _ = filepath.Abs(filepath.Join(fileSystem.Base, path))
	}
	return fullpath
}

// Get receive file with given path
func (fileSystem FileSystem) Get(path string) (*os.File, error) {
	return os.Open(fileSystem.GetFullPath(path))
}

// GetStream get file as stream
func (fileSystem FileSystem) GetStream(path string) {
	return ReadStream(fileSystem.GetFullPath(path))
}

// Put store a reader into given path
func (fileSystem FileSystem) Put(path string, reader io.Reader) (*ofs.Object, error) {
	var (
		fullpath = fileSystem.GetFullPath(path)
		err = CheckDir(filepath.Dir(fullpath))
	)

	if err != nil {
		return nil, err
	}

	dst, err := CheckFile(fullpath)

	if err == nil {
		if seeker, ok := reader.(io.ReadSeeker); ok {
			seeker.Seek(0, 0)
		}
		_, err = Copy(dst, reader)
	}

	return &ofs.Object{Path: path, Name: filepath.Base(path), StorageInterface: fileSystem}, err
}

// Delete delete file
func (fileSystem FileSystem) Delete(path string) error {
	return Remove(fileSystem.GetFullPath(path))
}

// List of all objects under current path
func (fileSystem FileSystem) List(path string) ([]*ofs.Object, error) {
	var (
		objects  []*ofs.Object
		fullpath = fileSystem.GetFullPath(path)
	)

	filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if path == fullpath {
			return nil
		}

		if err == nil && !info.IsDir() {
			modTime := info.ModTime()
			objects = append(objects, &ofs.Object{
				Path:             strings.TrimPrefix(path, fileSystem.Base),
				Name:             info.Name(),
				LastModified:     &modTime,
				StorageInterface: fileSystem,
			})
		}
		return nil
	})

	return objects, nil
}
