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
//
// pkgName 和 varName 分别指定包名和变量名；
// fileHeader 指定了文件头，如果为空，则不会输出文件内容；
// tag 指定了 // +build 指令下的标签，如果为空则不生成该指令；
func File(v interface{}, pkgName, varName, fileHeader, tag, path string) error {
	w := errwrap.Buffer{}

	if fileHeader != "" {
		w.Printf("// %s \n\n", fileHeader)
	}

	if tag != "" {
		w.Printf("// +build %s \n\n", tag)
	}

	bs, err := Bytes(v)
	if err != nil {
		return err
	}
	w.Printf("package %s \n\n", pkgName).
		Printf("const %s = `%s`\n", varName, string(bs))

	if w.Err != nil {
		return w.Err
	}

	src, err := format.Source(w.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, src, os.ModePerm)
}

// Bytes 将 v 打包成一份合法的 Go []byte 类型
func Bytes(v interface{}) ([]byte, error) {
	str, err := String(v)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// String 将 v 打包成一份合法的 Go 字符串
func String(v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	g := gob.NewEncoder(buf)

	if err := g.Encode(v); err != nil {
		return "", err
	}

	return base64Encoding.EncodeToString(buf.Bytes()), nil
}

// Unpack 用于解压由 Bytes 或是 String 打包的内容
func Unpack(buffer string, v interface{}) error {
	buf, err := base64Encoding.DecodeString(buffer)
	if err != nil {
		return err
	}

	g := gob.NewDecoder(bytes.NewReader(buf))
	return g.Decode(v)
}
