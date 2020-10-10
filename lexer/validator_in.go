package lexer

import (
	"fmt"
	"strings"
)

//InCheckerName 包含验证器
const InCheckerName = "in"

//InValueDelimiter in值分割字符
const InValueDelimiter = ","

func init() {
	SetCheckMethods(InCheckerName, NewInValidator)
}

type inValidator struct {
	*AbstractValidator
	valueType  string
	exprPrefix string
	warpSymbol string
}

func (in *inValidator) GenCode() string {
	msg := in.Checker.Exp.Msg
	comp := CompToSymbol(in.Checker.Comp)
	value := in.Checker.Value

	str := in.Indent
	//拆分value
	valuesIn := strings.Split(value, InValueDelimiter)
	if len(valuesIn) == 0 {
		return str + "return nil"
	}

	logicToken := ""
	if comp == tkEq {
		logicToken = and
	} else {
		logicToken = or
	}

	str += "if !("
	for i, val := range valuesIn {
		str += in.exprPrefix + " " + comp + " " + in.warpSymbol + val + in.warpSymbol
		if i < len(valuesIn)-1 {
			str += " " + logicToken + " "
		}
	}
	str += " ) " + in.GetReturn(msg)

	return str
}

func (in *inValidator) checkParsable() (err error) {

	in.exprPrefix = in.Checker.Exp.Lexer.PointerAlias + "." + in.Checker.Exp.Field.Name()

	if strings.HasPrefix(in.valueType, "int") || strings.HasPrefix(in.valueType, "uint") || strings.HasPrefix(in.valueType, "float") {
		return
	}

	if strings.HasPrefix(in.valueType, "string") {
		in.warpSymbol = "\""
		return
	}
	if strings.HasPrefix(in.valueType, "byte") || strings.HasPrefix(in.valueType, "rune") {
		in.warpSymbol = "'"
		return
	}
	if strings.HasPrefix(in.valueType, "*") {
		in.Checker.Value = ""
		return
	}
	err = fmt.Errorf("%s is unsupported type, number/string/byte/rune is available", in.valueType)
	return
}

//NewInValidator 返回存在验证器
func NewInValidator(checker *Checker) (IValidator, error) {
	v := &AbstractValidator{Checker: checker}
	valueType := v.Checker.Exp.Field.Type().String()
	vv := &inValidator{v, valueType, "", ""}
	if err := vv.checkParsable(); err != nil {
		return nil, err
	}
	return vv, nil
}
