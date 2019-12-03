package fs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const TestDir = "TEMP"

func init() {
	err := CheckDir(TestDir)
	if err != nil {
		panic(err)
	}
}

// Test Mkdir
func Test_MkDir(t *testing.T) {
	var (
		dirPath = path.Join(TestDir, "mkdir_test_dir")
		err     error
	)

	if err = Mkdir(dirPath); err != nil {
		t.Errorf("Create mkdir fial")
		return
	}

	defer func() {
		Remove(dirPath)
	}()

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Test Mkdir fail")
		return
	}
}

// Check if Directory Exist
func Test_CheckDir(t *testing.T) {
	var (
		dirPath = path.Join(TestDir, "check_test_dir")
		err     error
	)

	if err = CheckDir(dirPath); err != nil {
		t.Errorf("Check dir fail")
	}

	defer func() {
		Remove(dirPath)
	}()

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Check dir fail")
		return
	}
}

// Chec If Parent Directory Not Exist Directory
func Test_CheckDirIfParentDirNotExist(t *testing.T) {
	var (
		dirPath = path.Join(TestDir, "check_parent_test_dir/parent/nest")
		err     error
	)

	if err = CheckDir(dirPath); err != nil {
		t.Errorf("Check dir fail")
		return
	}

	defer func() {
		Remove(path.Join(TestDir, "check_parent_test_dir"))
	}()

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("Check dir fail")
		return
	}
}

// Read Directory testing
func Test_ReadDir(t *testing.T) {
	var (
		distDir = path.Join("./")
	)

	if files, err := ReadDir(distDir); err != nil {
		t.Error("ReadDir Fail.")
	} else {
		if len(files) != 3 {
			t.Error("ReadDir Fail.")
			return
		}
	}
}

// Creating Temporary Directory testing
func Test_TempDir(t *testing.T) {
	if dir, err := TempDir("TEMP", "test_temp_dir"); err != nil {
		t.Errorf("Create temp dir fail %v", err.Error())
		return
	} else {
		defer os.RemoveAll(dir)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Errorf("Read temp dir fail %v", err.Error())
			return
		}

		if len(files) != 0 {
			t.Errorf("Temp dir should be empty")
			return
		}

	}
}

// Check if file exist in path
func Test_PathExistsFile(t *testing.T) {
	var (
		filePath = path.Join(TestDir, "pathExists_test.file")
		err      error
	)

	if isExist := PathExists(filePath); isExist {
		t.Errorf("Path shoud not exist before create")
		return
	}

	if err = ioutil.WriteFile(filePath, []byte("Hello May CMF"), os.ModePerm); err != nil {
		t.Errorf("Create file fail")
		return
	}

	defer func() {
		os.RemoveAll(filePath)
	}()

	if isExist := PathExists(filePath); !isExist {
		t.Errorf("Path shoud exist after create")
		return
	}
}

// Check if Directory exist with path
func Test_PathExistsDir(t *testing.T) {
	var (
		filePath = path.Join(TestDir, "pathExists_test_dir")
		err      error
	)

	if isExist := PathExists(filePath); isExist {
		t.Errorf("Path shoud not be exist before create")
		return
	}

	if err = os.Mkdir(filePath, os.ModePerm); err != nil {
		t.Errorf("Creating directory is fail")
		return
	}

	defer func() {
		os.RemoveAll(filePath)
	}()

	if isExist := PathExists(filePath); !isExist {
		t.Errorf("Path shoud exist after create")
		return
	}
}

func Test_PathCopyFile(t *testing.T) {
	var (
		srcPath  = path.Join(TestDir, "copy_test.file1")
		distPath = path.Join(TestDir, "copy_test.file2")
		content  = "Hello MyCMF"
		err      error
	)

	// before copy, the new file should not be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) == false {
		t.Errorf("Path shoud not exist before create")
		return
	}

	if err = ioutil.WriteFile(srcPath, []byte(content), os.ModePerm); err != nil {
		t.Errorf("Creating file is failed")
		return
	}

	defer func() {
		os.RemoveAll(srcPath)
	}()

	if err = Copy(srcPath, distPath); err != nil {
		t.Errorf("Copying file is failed")
		return
	}

	defer func() {
		os.RemoveAll(distPath)
	}()

	//after copy the new file should be exist
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		t.Errorf("src file should be there")
		return
	}

	//after copy the new file should be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		t.Errorf("dist file should be exist")
		return
	}

	// Content should exist
	if d, err := ioutil.ReadFile(distPath); err != nil {
		t.Errorf("dist file can not read")
		return
	} else {
		if string(d[:]) != content {
			t.Errorf("dist file shoud have same content with src file")
			return
		}
	}
}

// Copy File If It Not Exist
func Test_PathCopyFileIfItNotExist(t *testing.T) {
	var (
		srcPath  = path.Join(TestDir, "copy_not_exist_test.file1")
		distPath = path.Join(TestDir, "copy_not_exist_test.file2")
		err      error
	)
	if err = Copy(srcPath, distPath); err == nil {
		t.Error("Copy an unexist file, it should throw an error")
		return
	}
}

// Copy Dir testing
func Test_PathCopyDir(t *testing.T) {
	var (
		srcPath   = path.Join(TestDir, "copy_test_dir1")
		distPath  = path.Join(TestDir, "copy_test_dir2")
		file1     = path.Join(srcPath, "copy_test_dir_child.file1")
		file2     = path.Join(distPath, "copy_test_dir_child.file1")
		file1Byte []byte
		file2Byte []byte
		content   = "hello MayCMF"
		err       error
	)

	// before copy, the new file should not be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) == false {
		t.Errorf("path shoud not exist before create")
		return
	}

	if err = os.Mkdir(srcPath, os.ModePerm); err != nil {
		t.Errorf("create dir fail %v", err.Error())
		return
	}

	// remove src dir
	defer func() {
		os.RemoveAll(srcPath)
	}()

	// write file in the dir
	if err = ioutil.WriteFile(file1, []byte(content), os.ModePerm); err != nil {
		t.Errorf("create file fail %v", err.Error())
		return
	}

	// remove src file
	defer func() {
		os.RemoveAll(file1)
	}()

	if err = Copy(srcPath, distPath); err != nil {
		t.Errorf("copy file fail %v", err.Error())
		return
	}

	// remove dist dir
	defer func() {
		os.RemoveAll(distPath)
	}()
	// remove dist file
	defer func() {
		os.RemoveAll(file2)
	}()

	//after copy the new file should be exist
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		t.Errorf("src file should still there")
		return
	}

	//after copy the new file should be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		t.Errorf("dist file should create")
		return
	}

	if file1Byte, err = ioutil.ReadFile(file1); err != nil {
		t.Errorf("read file1 fail %v", err.Error())
		return
	}

	if file2Byte, err = ioutil.ReadFile(file2); err != nil {
		t.Errorf("read file2 fail %v", err.Error())
		return
	}

	file1Content := string(file1Byte[:])
	file2Content := string(file2Byte[:])

	if file1Content != file2Content {
		t.Errorf("the copy file should be the same")
		return
	}
}

// Test Write File
func Test_WriteFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "write_test.file")
		err      error
		content  = "hello MyCMF"
		data     []byte
	)

	if err = WriteFile(filepath, []byte(content)); err != nil {
		t.Errorf("Write file failed")
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	if data, err = ioutil.ReadFile(filepath); err != nil {
		t.Errorf("Write file failed")
		return
	}

	s := string(data[:])

	if s != content {
		t.Errorf("Write file failed")
		return
	}

}

// Test Read file
func Test_ReadFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "read_test.file")
		err      error
		content  = "Hello MyCMF"
		data     []byte
	)

	if err = ioutil.WriteFile(filepath, []byte(content), os.ModePerm); err != nil {
		panic(err)
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	if data, err = ReadFile(filepath); err != nil {
		t.Errorf("Write file failed")
		return
	}

	s := string(data[:])

	if s != content {
		t.Errorf("Write file failed")
		return
	}

}

// Check File Exits, if not create It
func Test_CheckFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "check_test.file")
		err      error
		content  string
		data     []byte
	)

	if err = CheckFile(filepath); err != nil {
		t.Errorf("Check file fail")
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	// Read file
	if data, err = ioutil.ReadFile(filepath); err != nil {
		panic(err)
	}

	content = string(data[:])

	if content != "" {
		t.Errorf("Check file failed")
		return
	}

}

// Check File If Parent Dir Not Exist
func Test_CheckFileIfParentDirNotExist(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "CheckFile_nest_dir", "check_test.file")
		err      error
		content  string
		data     []byte
	)

	if err = CheckFile(filepath); err != nil {
		t.Errorf("Check file failed")
		return
	}

	defer func() {
		os.RemoveAll(path.Dir(filepath))
	}()

	// read file
	if data, err = ioutil.ReadFile(filepath); err != nil {
		panic(err)
	}

	content = string(data[:])

	if content != "" {
		t.Errorf("Check file faile")
		return
	}

}

// Append to sFile
func Test_AppendFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "append_test.file")
		err      error
		content  = "Hello"
		data     []byte
	)

	if err = ioutil.WriteFile(filepath, []byte(content), os.ModePerm); err != nil {
		panic(err)
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	if err = AppendFile(filepath, []byte(" MyCMF")); err != nil {
		panic(err)
	}

	// read file
	if data, err = ioutil.ReadFile(filepath); err != nil {
		panic(err)
	}

	content = string(data[:])

	if content != "Hello MyCMF" {
		t.Errorf("append file fail")
		return
	}

}


// Test Chmod of a file
func Test_ChmodFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "chmod_test.file")
		err      error
		content  = "Hello MayCMF"
		fileInfo os.FileInfo
	)

	if err = ioutil.WriteFile(filepath, []byte(content), 0777); err != nil {
		t.Errorf("Write file fail")
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	if err = Chmod(filepath, 0666); err != nil {
		panic(err)
		return
	}

	if fileInfo, err = os.Stat(filepath); err != nil {
		panic(err)
		return
	}

	if fileInfo.Mode() != 0666 {
		t.Errorf("Change mode of a file fail")
		return
	}

}


// Test Chmod of a file
func Test_ChmodDir(t *testing.T) {
	t.Skip()
	var (
		dirPath  = path.Join(TestDir, "chmod_test_dir")
		err      error
		fileInfo os.FileInfo
	)

	// create dir
	if err = os.Mkdir(dirPath, 0666); err != nil {
		panic(err)
		return
	}

	defer func() {
		os.RemoveAll(dirPath)
	}()

	// change mode
	if err = Chmod(dirPath, 0777); err != nil {
		t.Errorf("Test chmod is failed")
		return
	}

	// get mode
	if fileInfo, err = os.Stat(dirPath); err != nil {
		panic(err)
		return
	}

	if fileInfo.Mode() != 0777 {
		t.Errorf("Change directory mode is failed")
		return
	}

}

func Test_PathMoveFile(t *testing.T) {
	var (
		srcPath      = path.Join(TestDir, "move_test.file1")
		distPath     = path.Join(TestDir, "move_test.file2")
		content      = "Hello MayCMF"
		err          error
		mode         os.FileMode
		fileInfo     os.FileInfo
		distFileInfo os.FileInfo
	)

	if err = ioutil.WriteFile(srcPath, []byte(content), os.ModePerm); err != nil {
		t.Errorf("Create file fail")
		return
	}

	// before copy, the new file shouldn't be exist
	if fileInfo, err = os.Stat(srcPath); os.IsNotExist(err) {
		t.Errorf("Path shoud be exist after created")
		return
	}

	mode = fileInfo.Mode()

	if err = Move(srcPath, distPath); err != nil {
		t.Errorf("Move file fail %v", err.Error())
		return
	}

	defer func() {
		os.RemoveAll(distPath)
	}()

	//after copy the new file should be exist
	if distFileInfo, err = os.Stat(distPath); os.IsNotExist(err) {
		t.Errorf("Destination file should create")
		return
	}

	// mode same
	if distFileInfo.Mode() != mode {
		t.Errorf("Mode not same after move file")
		return
	}

	// and file should be some
	if d, err := ioutil.ReadFile(distPath); err != nil {
		t.Errorf("Destination file can not be read %v", err.Error())
		return
	} else {
		// content same
		if string(d[:]) != content {
			t.Errorf("Destination file shoud have same content with src file")
			return
		}
	}
}

// Path Move Dir
func Test_PathMoveDir(t *testing.T) {
	var (
		srcPath      = path.Join(TestDir, "copy_test_dir1")
		distPath     = path.Join(TestDir, "copy_test_dir2")
		file1        = path.Join(srcPath, "copy_test_dir_child.file1")
		file2        = path.Join(distPath, "copy_test_dir_child.file1")
		file2Byte    []byte
		content      = "Hello MayCMF"
		err          error
		mode         os.FileMode
		srcFileInfo  os.FileInfo
		distFileInfo os.FileInfo
	)

	// before copy, the new file should not be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) == false {
		t.Errorf("Path shoud not exist before create")
		return
	}

	if err = os.Mkdir(srcPath, os.ModePerm); err != nil {
		t.Errorf("Create dir fail %v", err.Error())
		return
	}

	// write file in the dir
	if err = ioutil.WriteFile(file1, []byte(content), os.ModePerm); err != nil {
		t.Errorf("Create file fail %v", err.Error())
		return
	}

	if srcFileInfo, err = os.Stat(file1); err != nil {
		t.Errorf("Stat file1 fail %v", err.Error())
		return
	}

	mode = srcFileInfo.Mode()

	if err = Move(srcPath, distPath); err != nil {
		t.Errorf("Copy file fail %v", err.Error())
		return
	}

	// remove dist dir
	defer func() {
		os.RemoveAll(distPath)
	}()

	// after move, the src directory should be remove
	if _, err := os.Stat(srcPath); os.IsNotExist(err) == false {
		t.Errorf("src file should be move")
		return
	}

	// after move, the destination directory should be exist
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		t.Errorf("Destination file should create")
		return
	}

	// and the new file should readable
	if file2Byte, err = ioutil.ReadFile(file2); err != nil {
		t.Errorf("read file2 fail %v", err.Error())
		return
	}

	file2Content := string(file2Byte[:])

	// same content
	if file2Content != content {
		t.Errorf("Copied file should be the same")
		return
	}

	// same mode
	if distFileInfo, err = os.Stat(file2); err != nil {
		t.Errorf("Stat file2 fail %v", err.Error())
		return
	}

	if distFileInfo.Mode() != mode {
		t.Errorf("Copied file should be the same")
		return
	}

}

// Remove File
func Test_RemoveFile(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "remove_test.file")
		err      error
		content  = "Hello MayCMF"
	)

	if err = ioutil.WriteFile(filepath, []byte(content), 0777); err != nil {
		t.Errorf("Write file fail %v", err.Error())
		return
	}

	if err = Remove(filepath); err != nil {
		t.Errorf("Remove file fail %v", err.Error())
		return
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) == false {
		t.Errorf("File should be remove")
		return
	}
}

// Remove Directory
func Test_RemoveDir(t *testing.T) {
	var (
		dirpath = path.Join(TestDir, "remove_test_dir")
		err     error
	)

	if err = os.Mkdir(dirpath, 0777); err != nil {
		t.Errorf("Mkdir fail %v", err.Error())
		return
	}

	if err = Remove(dirpath); err != nil {
		t.Errorf("Remove directory fail %v", err.Error())
		return
	}

	if _, err := os.Stat(dirpath); os.IsNotExist(err) == false {
		t.Errorf("Directory should be remove")
		return
	}
}


// Test Chmod of a file
func Test_Stat(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "stat_test.file")
		err      error
		content  = "Hello MayCMF"
	)

	if err = ioutil.WriteFile(filepath, []byte(content), 0777); err != nil {
		t.Errorf("Write file fail")
		return
	}

	defer func() {
		os.RemoveAll(filepath)
	}()

	if _, err = Stat(filepath); err != nil {
		t.Errorf("Test stat fail")
		return
	}

}

// Test Chmod of a file
func Test_LStat(t *testing.T) {
	var (
		filepath = path.Join(TestDir, "lstat_test.txt")
		err      error
		content  = "Hello MayCMF"
	)

	if err = ioutil.WriteFile(filepath, []byte(content), 0777); err != nil {
		t.Errorf("Write file failed")
		return
	}

	defer func() {
		os.RemoveAll(filepath)
		os.RemoveAll(TestDir)
	}()

	if _, err = LStat(filepath); err != nil {
		t.Errorf("Test stat failed")
		return
	}
}
