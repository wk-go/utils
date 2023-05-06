package utils

import (
	"archive/zip"
	"bytes"
	"container/list"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// zip path
func Zip(dst string, files ...string) error {
	rootPath := files[0]
	files = files[1:]
	dstFile, _ := os.Create(dst)
	defer dstFile.Close()
	zipWriter := zip.NewWriter(dstFile)
	defer zipWriter.Close()
	stack := list.New()
	if len(files) > 0 {
		for _, fileString := range files {
			file, err := os.Open(rootPath + "/" + fileString)
			if err != nil {
				log.Fatal(err)
			}
			stack.PushFront(file)
		}
	} else {
		file, err := os.Open(rootPath)
		if err != nil {
			log.Fatal(err)
		}
		stack.PushFront(file)
	}
	for fileEl := stack.Back(); fileEl != nil; fileEl = stack.Back() {
		stack.Remove(fileEl)
		file, ok := fileEl.Value.(*os.File)
		if ok {

		}
		info, err := file.Stat()
		if err != nil {
			continue
		}
		if info.IsDir() {
			dirList, err := file.Readdir(-1)
			if err != nil {
				log.Println(err)
			}

			p := file.Name()
			for _, subInfo := range dirList {
				subFile, err := os.Open(p + "/" + subInfo.Name())
				if err != nil {
					continue
				}
				stack.PushBack(subFile)
			}
			continue
		}
		header, err := zip.FileInfoHeader(info)
		if pathZip := strings.Replace(path.Dir(file.Name()), rootPath, "", 1); !(pathZip == "/" || pathZip == "") {
			header.Name = pathZip[1:] + "/" + header.Name
		}
		if err != nil {
			return err
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func Unzip(zipFile, dst string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	return UnzipFromReader(&reader.Reader, dst)
}

func UnzipFromBytes(b []byte, dst string) error {
	reader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}
	return UnzipFromReader(reader, dst)
}

func UnzipFromReader(reader *zip.Reader, dst string) error {
	copy := func(file *zip.File, dst string) error {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		decodeName := file.Name
		if file.Flags == 0 {
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		}

		filename := dst + "/" + decodeName
		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			var _err *fs.PathError
			if errors.As(err, &_err) {
				return nil
			}
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		return nil
	}
	for _, file := range reader.File {
		err := copy(file, dst)
		if err != nil {
			return err
		}
	}
	return nil
}
