package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

//LenCheckerName 长度验证器token
const LenCheckerName = "len"

func init() {
	SetCheckMethods(LenCheckerName, NewLenValidator)
}

type lenValidator struct {
	*AbstractValidator
}

func (v *lenValidator) GenCode() string {
	fieldName := v.Checker.Exp.Field.Name()
	msg := v.Checker.Exp.Msg
	comp := CompToSymbol(v.Checker.Comp)
	value := v.Checker.Value
	alias := v.Checker.Exp.Lexer.PointerAlias

	return v.Indent +
		"if !(len(" + alias + "." + fieldName + ") " + comp + " " + value + ")" +
		v.GetReturn(msg)
}

func (v *lenValidator) checkParsable() error {
	if _, err := strconv.ParseInt(v.Checker.Value, 10, 64); err != nil {
		return fmt.Errorf("%s can not parse to int", v.Checker.Value)
	}
	valueType := v.Checker.Exp.Field.Type().String()
	if !(strings.HasPrefix(valueType, "[") || strings.HasPrefix(valueType, "map") || strings.HasPrefix(valueType, "string")) {
		return fmt.Errorf("%s is unsupported type", valueType)
	}

	if ok := CompToSymbol(v.Checker.Comp); ok == "" {
		return fmt.Errorf("compare symbol invalid %s", v.Checker.Comp)
	}

	return nil
}

//NewLenValidator 返回长度验证器
func NewLenValidator(checker *Checker) (IValidator, error) {
	v := &AbstractValidator{Checker: checker}
	vv := &lenValidator{v}
	if err := vv.checkParsable(); err != nil {
		return nil, err
	}
	return vv, nil
}
