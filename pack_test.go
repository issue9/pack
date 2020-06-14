// SPDX-License-Identifier: MIT

package pack

import (
	"testing"

	"github.com/issue9/assert"
)

func TestString(t *testing.T) {
	a := assert.New(t)

	str := "这是字符串的测试`"
	err := File(str, "testdata", "String", "请勿修改", "pack", "./testdata/string.go")
	a.NotError(err).FileExists("./testdata/string.go")

	v, err := String(&str)
	a.NotError(err).NotEmpty(v)

	var v2 string
	a.NotError(Unpack(v, &v2))
	a.Equal(v2, str)
}

type obj struct {
	ID   int
	Name string
}

func TestObjet(t *testing.T) {
	a := assert.New(t)

	o := &obj{ID: 1, Name: "111"}
	err := File(o, "testdata", "Object", "请勿修改", "pack", "./testdata/obj.go")
	a.NotError(err).FileExists("./testdata/obj.go")

	objs := []*obj{
		{ID: 1, Name: "name"},
		{ID: 2, Name: "name2"},
		{ID: 3, Name: "name3"},
	}
	err = File(objs, "testdata", "Objects", "请勿修改", "pack", "./testdata/objs.go")
	a.NotError(err).FileExists("./testdata/objs.go")
}
