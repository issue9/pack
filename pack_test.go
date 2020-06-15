// SPDX-License-Identifier: MIT

package pack

import (
	"testing"

	"github.com/issue9/assert"
)

func TestString(t *testing.T) {
	a := assert.New(t)

	str := "这是字符串的测试`"
	err := File(str, "testdir", "String", "请勿修改", "pack", "./testdir/string.go")
	a.NotError(err).FileExists("./testdir/string.go")

	// String

	v, err := String(&str)
	a.NotError(err).NotEmpty(v)

	var v2 string
	a.NotError(Unpack(v, &v2))
	a.Equal(v2, str)

	// Bytes

	bs, err := Bytes([]byte(str))
	a.NotError(err).NotNil(bs)

	var bs2 []byte
	a.NotError(Unpack(string(bs), &bs2))
	a.Equal(v2, []byte(str))
}

type obj struct {
	ID   int
	Name string
}

func TestObjet(t *testing.T) {
	a := assert.New(t)

	o := &obj{ID: 1, Name: "111"}
	err := File(o, "testdir", "Object", "请勿修改", "pack", "./testdir/obj.go")
	a.NotError(err).FileExists("./testdir/obj.go")

	objs := []*obj{
		{ID: 1, Name: "name"},
		{ID: 2, Name: "name2"},
		{ID: 3, Name: "name3"},
	}
	err = File(objs, "testdir", "Objects", "请勿修改", "pack", "./testdir/objs.go")
	a.NotError(err).FileExists("./testdir/objs.go")
}
