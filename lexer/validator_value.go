package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

//ValueCheckerName 值验证方法名token
const ValueCheckerName = "value"

func init() {
	SetCheckMethods(ValueCheckerName, NewValueValidator)
}

type valueValidator struct {
	*AbstractValidator
	valueType  string
	exprPrefix string
}

func (v *valueValidator) GenCode() string {
	expressPrefix := v.exprPrefix
	msg := v.Checker.Exp.Msg
	comp := CompToSymbol(v.Checker.Comp)
	value := v.Checker.Value

	return v.Indent + "if !(" + expressPrefix + " " + comp + " " + value + ")" + v.GetReturn(msg)
}

func (v *valueValidator) checkParsable() (err error) {
	v.exprPrefix = v.Checker.Exp.Lexer.PointerAlias + "." + v.Checker.Exp.Field.Name()

	//数字型 需要验证是否能转为数字型值
	if strings.HasPrefix(v.valueType, "int") || strings.HasPrefix(v.valueType, "uint") {
		if _, err = strconv.ParseInt(v.Checker.Value, 10, 64); err != nil {
			err = fmt.Errorf("%s can not parse to int", v.Checker.Value)
			return
		}
		return
	}
	//浮点型需要校验值是否是浮点数
	if strings.HasPrefix(v.valueType, "float") {
		if _, err = strconv.ParseFloat(v.Checker.Value, 64); err != nil {
			err = fmt.Errorf("%s can not parse to float", v.Checker.Value)
			return
		}
		return
	}
	//切片和map不能用此方法校验
	if strings.HasPrefix(v.valueType, "[") || strings.HasPrefix(v.valueType, "map") {
		err = fmt.Errorf("%s can not valid slice(array) or map", v.Checker.Value)
		return
	}
	//字符型 需要加双引号包裹
	if strings.HasPrefix(v.valueType, "string") {
		v.Checker.Value = "\"" + v.Checker.Value + "\""
		return
	}
	//字节型 需要单引号包裹
	if strings.HasPrefix(v.valueType, "byte") || strings.HasPrefix(v.valueType, "rune") {
		v.Checker.Value = "'" + v.Checker.Value + "'"
		return
	}

	err = fmt.Errorf("%s is unsupported type", v.valueType)
	return
}

//NewValueValidator 返回值验证器
func NewValueValidator(checker *Checker) (IValidator, error) {
	v := &AbstractValidator{Checker: checker}
	valueType := v.Checker.Exp.Field.Type().String()
	vv := &valueValidator{v, valueType, ""}
	if err := vv.checkParsable(); err != nil {
		return nil, err
	}
	return vv, nil
}
