// SPDX-License-Identifier: MIT

package pack

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DirFile 打包 dir 下所有的文件内容
func DirFile(root string, pkgName, varName, fileHeader, tag, path string) error {
	value, err := DirString(root)
	if err != nil {
		return err
	}
	return writeToFile(value, pkgName, varName, fileHeader, tag, path)
}

// DirBytes 打包 dir 下所有的文件内容
func DirBytes(root string) ([]byte, error) {
	str, err := DirString(root)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// DirString 打包 dir 下所有的文件内容
func DirString(root string) (string, error) {
	fs := &fileSystem{}

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		f := &file{
			FileName: filepath.ToSlash(rel),
			Created:  info.ModTime(),
		}

		if !info.IsDir() {
			f.Content, err = ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		} else if f.FileName[len(f.FileName)-1] != '/' {
			f.FileName += "/" // 如果是路径，保证以 / 结尾
		}

		fs.Files = append(fs.Files, f)
		return nil
	}

	if err := filepath.Walk(root, walk); err != nil {
		return "", err
	}

	return String(fs)
}

type fileSystem struct {
	Files []*file

	// 以下这些默认为 nil，由 NewFileSystem 加载时初始化
	indexes map[string]int
}

type file struct {
	FileName string
	Created  time.Time
	Content  []byte // 如果是文件，也有可能 content 长度为 0

	// 以下这些默认为 nil，由 NewFileSystem 加载时初始化

	reader *bytes.Reader
	fs     *fileSystem
}

// NewFileSystem 将 data 解包并生成符合 http.FileSystem 接口的实例
//
// data 必须是由 DirBytes 或是 DirString 打包后的数据。
func NewFileSystem(data string) (http.FileSystem, error) {
	fs := &fileSystem{indexes: map[string]int{}}
	if err := Unpack(data, fs); err != nil {
		return nil, err
	}

	for _, file := range fs.Files {
		file.reader = bytes.NewReader(file.Content)
		file.fs = fs
	}

	return fs, nil
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	if name != "" && name[0] == '/' {
		name = name[1:]
	}

	if name == "" {
		name = "./"
	}

	for _, file := range fs.Files {
		if file.Name() == name || (file.IsDir() == true && file.Name() == name+"/") {
			return file, nil
		}
	}
	return nil, os.ErrNotExist
}

func (f *file) Close() error {
	_, err := f.Seek(0, io.SeekStart)
	return err
}

func (f *file) Read(p []byte) (int, error) {
	if f.IsDir() || f.reader == nil {
		return 0, io.EOF
	}
	return f.reader.Read(p)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	if f.IsDir() || f.reader == nil {
		return 0, io.EOF
	}
	return f.reader.Seek(offset, whence)
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	if !f.IsDir() {
		return nil, nil
	}

	var fis []os.FileInfo

	index := f.fs.indexes[f.FileName]
	var i int
	var item *file
	for i, item = range f.fs.Files {
		if !strings.HasPrefix(item.FileName, f.FileName) || i < index {
			continue
		}

		if count > 0 && len(fis) >= count {
			break
		}
		fis = append(fis, item)
	}

	if count <= 0 || i >= len(f.fs.Files)-1 {
		f.fs.indexes[f.FileName] = 0
	}

	return fis, nil
}

func (f *file) Stat() (os.FileInfo, error) { return f, nil }

func (f *file) Name() string { return f.FileName }

func (f *file) Size() int64 { return int64(len(f.Content)) }

func (f *file) Mode() os.FileMode { return 0755 }

func (f *file) ModTime() time.Time { return f.Created }

func (f *file) IsDir() bool { return f.Content == nil }

func (f *file) Sys() interface{} { return nil }
