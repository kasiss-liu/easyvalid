package lexer

import (
	"fmt"
	"strings"
)

//NotEmptyCheckerName 验证不为空的方法token
const NotEmptyCheckerName = "notEmpty"

func init() {
	SetCheckMethods(NotEmptyCheckerName, NewNotEmptyValidator)
}

type notEmptyValidator struct {
	*AbstractValidator
	valueType  string
	exprPrefix string
}

func (v *notEmptyValidator) GenCode() string {
	expressPrefix := v.exprPrefix
	msg := v.Checker.Exp.Msg
	comp := CompToSymbol(v.Checker.Comp)
	value := v.Checker.Value

	return v.Indent + "if " + expressPrefix + " " + comp + " " + value + v.GetReturn(msg)
}

func (v *notEmptyValidator) checkParsable() (err error) {
	v.exprPrefix = v.Checker.Exp.Lexer.PointerAlias + "." + v.Checker.Exp.Field.Name()
	v.Checker.Comp = tkEq

	if strings.HasPrefix(v.valueType, "int") || strings.HasPrefix(v.valueType, "uint") || strings.HasPrefix(v.valueType, "float") {
		v.Checker.Value = "0"
		return
	}
	if strings.HasPrefix(v.valueType, "[") || strings.HasPrefix(v.valueType, "map") {
		v.exprPrefix = "len(" + v.exprPrefix + ")"
		v.Checker.Value = "0"
		return
	}
	if strings.HasPrefix(v.valueType, "string") {
		v.Checker.Value = "\"\""
		return
	}
	if strings.HasPrefix(v.valueType, "byte") || strings.HasPrefix(v.valueType, "rune") {
		v.Checker.Value = "''"
		return
	}
	if strings.HasPrefix(v.valueType, "*") {
		v.Checker.Value = "nil"
		return
	}
	var structName string
	if structName, err = v.checkPackage(v.valueType); err != nil {
		return
	}
	v.Checker.Value = "(" + structName + "{})"

	return
}

func (v *notEmptyValidator) checkPackage(valueType string) (structName string, err error) {
	pkgs := v.Checker.Exp.Lexer.imports
	valueType = strings.Trim(valueType, "*")
	for _, p := range pkgs {
		if strings.Contains(valueType, p.ID) {
			v.Checker.Exp.Lexer.usedImports = append(v.Checker.Exp.Lexer.usedImports, p)
			return p.Name + strings.TrimPrefix(valueType, p.ID), nil
		}
	}
	err = fmt.Errorf("%s is unsupported type", valueType)
	return "", err
}

//NewNotEmptyValidator 构造新的非空验证器
func NewNotEmptyValidator(checker *Checker) (IValidator, error) {
	v := &AbstractValidator{Checker: checker}
	valueType := v.Checker.Exp.Field.Type().String()
	vv := &notEmptyValidator{v, valueType, ""}
	if err := vv.checkParsable(); err != nil {
		return nil, err
	}
	return vv, nil
}
