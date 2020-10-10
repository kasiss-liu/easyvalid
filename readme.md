## easyvalid 

easyvalid 是一个简单的结构体校验代码生成工具。可根据结构体内tag字段的预定义，生成预定义的字段校验代码。

#### Usage
```shell script
#install
go get -u github.com/kasiss-liu/easyvalid

#run
easyvalid -p <file>.go
```

执行命令后，会在对应目录生成 ```<package_name>_easyvalid.go``` 的文件。

#### Options
```shell script
$ easyvalid  -h
Usage of easyvalid:
  -build_tags string
        build flags to attach
  -export
        attr valid funcs if exported
  -f    force rm file suffix with _easyvalid.go
  -h    show help info
  -note string
        note info to attach
  -o string
        output filename
  -p string
        package_path: a filepath or dir which package set
  -structs string
        struct valid need be generate. example: struct1,struct2
```
声明使用方式

```go
package main
//Example
// easyvalid:valid:{alias}
// alias 是结构体方法种的形参别名 example中是 e 即 easyvalid:valid:e
type Example struct {
    field1 string `evalid:"notEmpty|msg:field1不能为空"`
}
func (e Example) GetField1() string {
    return e.field1
}
```
支持的校验方法

```text
len 
    `evalid:"len:{比较运算符别名}({比较值})|{msg:{错误信息}}`
example:
    `evalid:"len:eq(1)|msg:长度必须为1"`

值类型必须为int
```
```text
value
    `evalid:"value:{比较运算符}({值})|{msg:"{错误信息}"}`
example:
    `evalid:"value:eq(hello)|msg:"不能为空"`
支持的数据类型
    int、int64、int32、uint、uint64、uint32
    float64、float32
    string
    byte、rune
```

```text
支持的比较运算符别名
    eq : ==
    neq: !=
    lt : <
    lte: <=
    gt : >
    gte: >=
```

```text
notEmpty
    `evalid:"notEmpty|{msg:"{错误信息}"}`
example:
    `evalid:"notEmpty|msg:"不能为空"`
支持的数据类型
    int、int64、int32、uint、uint64、uint32 会判断是否等于0
    float64、float32 会判断是否等于0
    string 会判断是否等于字符串空默认值("")
    byte、rune 会判断是否等于字节空默认值('')
    map,slice 会判断长度是否等于0
    其他类型 会判断是否等于nil
```

