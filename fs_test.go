// SPDX-License-Identifier: MIT

package pack

import (
	"net/http"
	"os"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"
)

var (
	_ http.FileSystem = &fileSystem{}
	_ http.File       = &file{}
)

func TestFileSystem(t *testing.T) {
	a := assert.New(t)

	data, err := DirString("./testdata")
	a.NotError(err).NotEmpty(data)

	fs, err := NewFileSystem(data)
	a.NotError(err).NotNil(fs)

	f, err := fs.Open("index.html")
	a.NotError(err).NotNil(f)
	stat, err := f.Stat()
	a.NotError(err).NotNil(stat).
		False(stat.IsDir()).
		True(stat.Size() > 0).
		True(stat.ModTime().Unix() > 0)

	f, err = fs.Open("assets/index.css")
	a.NotError(err).NotNil(f)

	f, err = fs.Open("/assets/index.css")
	a.NotError(err).NotNil(f)

	f, err = fs.Open("assets/not.exists.file")
	a.ErrorType(err, os.ErrNotExist).Nil(f)

	f, err = fs.Open("assets/")
	a.NotError(err).NotNil(f)
	ff, ok := f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("assets")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("/assets")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("/assets/")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("/")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("emptyDir")
	a.NotError(err).NotNil(f)
	ff, ok = f.(*file)
	a.True(ok).NotNil(ff).True(ff.IsDir())

	f, err = fs.Open("not-exists-dir/")
	a.ErrorType(err, os.ErrNotExist).Nil(f)
}

func TestDir(t *testing.T) {
	a := assert.New(t)

	data, err := DirString("./testdata")
	a.NotError(err).NotEmpty(data)

	fs, err := NewFileSystem(data)
	a.NotError(err).NotNil(fs)

	srv := rest.NewServer(t, http.FileServer(fs), nil)
	srv.Get("/index.html").
		Do().
		Status(http.StatusOK)
	srv.Get("/").
		Do().
		Status(http.StatusOK)

	srv.Get("/not-exists").
		Do().
		Status(http.StatusNotFound)

	srv.Get("/assets/index.css").
		Do().
		Status(http.StatusOK).
		BodyNotEmpty()
	srv.Get("/assets/index.css"). // 二次读取
					Do().
					Status(http.StatusOK).
					BodyNotEmpty()

	srv.Get("/assets").
		Do().
		Status(http.StatusOK).
		BodyNotEmpty()

	srv.Get("/emptyDir").
		Do().
		Status(http.StatusOK)

	a.NotError(DirFile("./testdir", "testdir", "Object", "", "", "./testdir/dir.go"))
}
