package fs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// Check if path is exist
func PathExists(path string) (isExist bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Create a Directory
func Mkdir(dir string) (err error) {
	return os.Mkdir(dir, os.ModePerm)
}

// Check if the Directory exist
func CheckDir(dir string) (err error) {
	parent := path.Dir(dir)
	if _, err = os.Stat(parent); os.IsNotExist(err) {
		if err = CheckDir(parent); err != nil {
			return
		}
	}
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	return
}

// Read Directory and get file list
func ReadDir(dir string) (files []string, err error) {
	var (
		fileInfos []os.FileInfo
	)
	fileInfos, err = ioutil.ReadDir(dir)

	for _, f := range fileInfos {
		files = append(files, f.Name())
	}
	return
}

// Create random temporary Directory
func TempDir(dir string, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

// Check if the file exist
func CheckFile(filepath string) (err error) {
	var (
		file *os.File
	)
	// Check if the Directory exist
	if err = CheckDir(path.Dir(filepath)); err != nil {
		return err
	}
	// Check if file exist
	if _, err = os.Stat(filepath); os.IsNotExist(err) {
		file, err = os.Create(filepath)
		defer func() {
			file.Close()
		}()
	}
	return
}

// Write a file
func WriteFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

// Read a file
func ReadFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}

func AppendFile(file string, data []byte) error {
	if f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		defer func() {
			err = f.Close()
		}()
		if _, err := f.Write(data); err != nil {
			return err
		}
		return nil
	}
}

func Truncate(filepath string, len int64) error {
	return os.Truncate(filepath, len)
}

// Create a read stream
func ReadStream(filepath string) (stream io.Reader, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(filepath); err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	stream = bufio.NewReader(file)

	return
}

// Create a read stream
func WriteStream(filepath string) (stream io.Writer, err error) {
	var (
		file *os.File
	)
	if file, err = os.Open(filepath); err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	stream = bufio.NewWriter(file)

	return
}

// Almost the same as writeFile (i.e. it overwrites),
// except that if the parent directory does not exist, it's created.
func OuputFile(filepath string, data []byte) error {
	if err := CheckDir(path.Dir(filepath)); err != nil {
		return err
	}
	return WriteFile(filepath, data)
}

// change the file permission
func Chmod(filepath string, mode os.FileMode) error {
	return os.Chmod(filepath, mode)
}

func Lchown(path string, uid int, gid int) error {
	return os.Lchown(path, uid, gid)
}

// Move File/Directory(Folder)
func Move(src string, target string) (err error) {

	var (
		fileInfo os.FileInfo
		files    []os.FileInfo
	)

	if fileInfo, err = os.Stat(src); err != nil {
		return
	}

	if fileInfo.IsDir() {

		// read dir and move one by one
		if files, err = ioutil.ReadDir(src); err != nil {
			return
		}

		if err = CheckDir(target); err != nil {
			return
		}

		for _, file := range files {
			filename := file.Name()
			srcFile := path.Join(src, filename)
			targetFile := path.Join(target, filename)
			if err = Move(srcFile, targetFile); err != nil {
				return err
			}
		}

		// copy all done, should remove the source directory
		err = os.RemoveAll(src)

		return
	}
	return os.Rename(src, target)
}

// Copy a File or Directory
func Copy(src string, target string) (err error) {

	var (
		srcFile    *os.File
		targetFile *os.File
		fileInfo   os.FileInfo
		files      []os.FileInfo
	)

	if fileInfo, err = os.Stat(src); err != nil {
		return
	}

	if fileInfo.IsDir() {
		// Read Directory and copy one by one
		if files, err = ioutil.ReadDir(src); err != nil {
			return
		}

		if err = CheckDir(target); err != nil {
			return
		}

		for _, file := range files {
			filename := file.Name()
			src = path.Join(src, filename)
			target = path.Join(target, filename)
			if err = Copy(src, target); err != nil {
				return err
			}
		}

	} else {
		if srcFile, err = os.Open(src); err != nil {
			return
		}

		defer func() {
			srcFile.Close()
		}()

		if targetFile, err = os.Create(target); err != nil {
			return
		}

		defer func() {
			targetFile.Close()
		}()

		_, err = io.Copy(targetFile, srcFile)
	}
	return
}

// Remove File/Directory
func Remove(path string) error {
	return os.RemoveAll(path)
}

// Rename File/Directory
func Rename(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// File Statistic (Stat)
func Stat(path string) (info os.FileInfo, err error) {
	return os.Stat(path)
}

// File Statistic (Stat)
func LStat(path string) (info os.FileInfo, err error) {
	return os.Lstat(path)
}
