package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type UploadFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Parent  string `json:"parent"`
	IsDir   bool   `json:"is_dir"`
	Mode    string `json:"mode"`
	ModTime string `json:"mod_time"`
	Size    string `json:"size"`
}

// The IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

// The IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func FileGetContents(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(content)
}

func ReadFile(filePath string) string {
	b, e := ioutil.ReadFile(filePath)
	if e != nil {
		fmt.Println("read file error")
		return e.Error()
	}
	return string(b)
}

// BasicFileInfo 收集文件基本信息
type BasicFileInfo struct {
	Filepath string
	ext      string
	mime     string
	hash     string
	Info     fs.FileInfo
}

func (bfi *BasicFileInfo) Name() string       { return bfi.Info.Name() }
func (bfi *BasicFileInfo) Size() int64        { return bfi.Info.Size() }
func (bfi *BasicFileInfo) Mode() fs.FileMode  { return bfi.Info.Mode() }
func (bfi *BasicFileInfo) ModTime() time.Time { return bfi.Info.ModTime() }
func (bfi *BasicFileInfo) IsDir() bool        { return bfi.Info.IsDir() }
func (bfi *BasicFileInfo) Sys() any           { return bfi.Info.Sys() }
func (bfi *BasicFileInfo) Ext() string        { return bfi.ext }
func (bfi *BasicFileInfo) MIME() string       { return bfi.mime }
func (bfi *BasicFileInfo) Hash() string       { return bfi.hash }
func (bfi *BasicFileInfo) Path() string       { return bfi.Filepath }

// GetBasicFileInfo 获取文件基本信息
func GetBasicFileInfo(pathToFile string) (*BasicFileInfo, error) {
	fi, err := os.Stat(pathToFile)
	if err != nil {
		return nil, err
	}
	// 处理基本属性
	fileExt := path.Ext(pathToFile)
	hash, _ := GenerateFileHash(pathToFile)

	bfi := &BasicFileInfo{
		Filepath: pathToFile,
		ext:      fileExt,
		mime:     mime.TypeByExtension(fileExt),
		hash:     hash,
		Info:     fi,
	}

	return bfi, nil
}

func ZipFilesToWrite(w1 io.Writer, fileList *[]string, pathNames map[string]string) error {
	if len(*fileList) < 1 {
		return fmt.Errorf("将要压缩的文件列表不能为空")
	}
	zw := zip.NewWriter(w1)
	defer zw.Close()

	for _, fileName := range *fileList {
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}

		stat, _ := fr.Stat()
		if stat.IsDir() {
			return errors.New("仅支持文件打包。")
		}

		// 写入文件的头信息
		var w io.Writer

		fileName = pathNames[filepath.Base(fileName)]
		w, err = zw.Create(filepath.Base(fileName))

		if err != nil {
			return err
		}

		// 写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return zw.Flush()
}

// JoinDir 使用系统分隔符连接目录
func JoinDir(dirs ...string) string {
	if len(dirs) == 0 {
		return ""
	}
	return strings.Join(dirs, string(os.PathSeparator))
}
