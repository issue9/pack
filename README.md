pack
[![Test Status](https://github.com/issue9/pack/workflows/Test/badge.svg?branch=master)](https://github.com/issue9/pack/actions?query=workflow%3ATest)
[![codecov](https://codecov.io/gh/issue9/pack/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/pack)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/issue9/pack)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
=====

用于将数据对象打包成一个 Go 文件

```go
type Object{
    // fields
}

o := &Object{}

// 在 ./static.go 文件生成 Obj 常量，其内容为编码后的 o 内容。
// 该内容可通过 Unpack 解码。
pack.File(o, "static", "Obj", "NOT EDIT", "","./static.go")
```

安装
----

```shell
go get github.com/issue9/pack
```

版权
----

本项目采用 [MIT](http://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
