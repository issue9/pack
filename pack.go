// SPDX-License-Identifier: MIT

package pack

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"go/format"
	"io/ioutil"
	"os"

	"github.com/issue9/errwrap"
)

var base64Encoding = base64.StdEncoding

// File 将对象 v 打包成一个 Go 文件内容中
func File(v interface{}, pkgName, varName, fileHeader, tag, path string) error {
	src, err := Bytes(v, pkgName, varName, fileHeader, tag)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, src, os.ModePerm)
}

// Bytes 将 v 打包成一份合法的 Go 格式文件并以 []byte 内容返回
//
// pkgName 和 varName 分别指定包名和变量名；
// fileHeader 指定了文件头，如果为空，则不会输出文件内容；
// tag 指定了 // +build 指令下的标签，如果为空则不生成该指令；
func Bytes(v interface{}, pkgName, varName, fileHeader, tag string) ([]byte, error) {
	buf := new(bytes.Buffer)
	g := gob.NewEncoder(buf)

	if err := g.Encode(v); err != nil {
		return nil, err
	}

	w := errwrap.Buffer{}

	if fileHeader != "" {
		w.Printf("// %s \n\n", fileHeader)
	}

	if tag != "" {
		w.Printf("// +build %s \n\n", tag)
	}

	content := base64Encoding.EncodeToString(buf.Bytes())
	w.Printf("package %s \n\n", pkgName).
		Printf("const %s = `%s`\n", varName, content)

	if w.Err != nil {
		return nil, w.Err
	}

	return format.Source(w.Bytes())
}

// Unpack 用于解压由 Pack 输出的内容
func Unpack(buffer string, v interface{}) error {
	buf, err := base64Encoding.DecodeString(buffer)
	if err != nil {
		return err
	}

	g := gob.NewDecoder(bytes.NewReader(buf))
	return g.Decode(v)
}
