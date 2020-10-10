package lexer

import "sync"

//IValidator 验证器接口
type IValidator interface {
	GenCode() string
	SetReturn(s string)
	SetIndent(s string)
}

//可用的验证方法
var allowMethods = make(map[string]func(*Checker) (IValidator, error))
var allowMethodLock sync.Mutex

//SetCheckMethods 注册验证器到方法池
func SetCheckMethods(method string, factory func(*Checker) (IValidator, error)) {
	allowMethodLock.Lock()
	allowMethods[method] = factory
	allowMethodLock.Unlock()
}

//AbstractValidator 验证器抽象类
type AbstractValidator struct {
	Checker *Checker
	Indent  string
	Rt      string
}

//SetReturn 设置换行符
func (vr *AbstractValidator) SetReturn(s string) {
	vr.Rt = s
}

//SetIndent 设置缩进
func (vr *AbstractValidator) SetIndent(s string) {
	vr.Indent = s
}

//GetReturn 返回return字符串
func (vr *AbstractValidator) GetReturn(msg string) string {
	return " { return errors.New(`" + msg + "`) }" + vr.Rt
}
