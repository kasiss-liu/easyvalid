package lexer

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

//定义一些比较运算符token
var (
	tkEq  = "eq"
	tkLt  = "lt"
	tkLte = "lte"
	tkGt  = "gt"
	tkGte = "gte"
	tkNEq = "neq"
)

//逻辑运算符
var (
	and = "&&"
	or  = "||"
)

//有效的比较运算符列表
var allowComp = map[string]bool{
	tkEq:  true,
	tkLt:  true,
	tkLte: true,
	tkGt:  true,
	tkGte: true,
}

//比较运算符转换map
var compToSymbolMap = map[string]string{
	tkEq:  "==",
	tkLt:  "<",
	tkLte: "<=",
	tkGt:  ">",
	tkGte: ">=",
	tkNEq: "!=",
}

//CompToSymbol 根据定义的比较运算符字符串 获取 比较运算发symbol
func CompToSymbol(s string) string {
	return compToSymbolMap[s]
}

//常用的字节
var (
	parenthesesLeft  = '('
	parenthesesRight = ')'
	space            = ' '
)

//tag中表示错误信息的token
var tkMsg = "msg"

//len:eq(1)|msg:st;
const sep = "|"   //信息分隔符
const colon = ":" //键值分隔符

//Lexer 字段tag分析器
type Lexer struct {
	field        *types.Var
	delimiter    string
	message      string
	exprStrings  []string
	exprs        []*express
	PointerAlias string
	imports      []*packages.Package
	usedImports  []*packages.Package
}

//SetImports 设置解析后的依赖
func (l *Lexer) SetImports(imt []*packages.Package) {
	l.imports = imt
}

//GetUsedImports 返回生成文件中使用过的package
func (l *Lexer) GetUsedImports() []*packages.Package {
	return l.usedImports
}

//Parse 分析tag
//先拆解表达式 拆解为多个验证器
//后拆解验证器，拆解为多个验证方法
//最后组合
func (l *Lexer) Parse() error {
	l.splitExpr()
	l.splitValid()
	if err := l.parseExpress(); err != nil {
		return err
	}
	return nil
}

//GetValidatorList  获取字段的验证器列表
func (l *Lexer) GetValidatorList() ([]IValidator, error) {
	IV := make([]IValidator, 0)
	for _, expr := range l.exprs {
		for _, c := range expr.Checkers {
			iValidator, ok := allowMethods[c.Method]
			if !ok {
				return nil, fmt.Errorf("struct \"%s\" field \"%s\" validator \"%s\" is invalid", "%s", "%s", c.Method)
			}
			validator, err := iValidator(c)
			if err != nil {
				return nil, fmt.Errorf("struct \"%s\" field \"%s\" validator \"%s\" parsed error: %s", "%s", "%s", c.Method, err.Error())
			}
			IV = append(IV, validator)
		}
	}
	return IV, nil
}

//splitExpr 表达式拆解
func (l *Lexer) splitExpr() {
	ss := strings.Split(l.message, l.delimiter)
	for _, s := range ss {
		if ps := strings.TrimSpace(s); ps != "" {
			l.exprStrings = append(l.exprStrings, ps)
		}
	}

}

//splitValid 将每个表达式拆解为多个验证方法
func (l *Lexer) splitValid() {
	for _, s := range l.exprStrings {
		ss := strings.Split(s, sep)
		exps := make([]string, 0)
		msg := "data invalid"
		for _, e := range ss {
			ee := strings.TrimSpace(e)
			if strings.HasPrefix(ee, tkMsg) {
				prefix := tkMsg + colon
				msg = strings.TrimSpace(strings.TrimPrefix(ee, prefix))
			} else {
				exps = append(exps, ee)
			}
		}
		l.exprs = append(l.exprs, &express{exps: exps, Msg: msg, Field: l.field, Lexer: l})
	}
}

//parseExpress 将每个表达式字符串 结构化为express
func (l *Lexer) parseExpress() error {
	for _, exp := range l.exprs {
		err := exp.parse()
		if err != nil {
			return err
		}
	}
	return nil
}

//express 字段tag解析后的表达式单元
type express struct {
	StName   string
	Field    *types.Var
	exps     []string
	Msg      string
	Checkers []*Checker
	Lexer    *Lexer
}

//Checker 每个字段表达式单元有多个checker
type Checker struct {
	Exp    *express
	Method string
	Comp   string
	Value  string
}

//parse 表达式单元解析 将表达式解析为checker
func (e *express) parse() error {
	for _, exp := range e.exps {
		ss := strings.Split(exp, colon)
		if len(ss) == 0 {
			continue
		}
		method := strings.TrimSpace(ss[0])
		if _, ok := allowMethods[method]; !ok {
			continue
		}

		comp := ""
		value := ""

		valString := ""
		if len(ss) > 1 {
			valString = strings.TrimSpace(ss[1])
			tokens := make([]string, 0)
			for compTk := range allowComp {
				if strings.HasPrefix(valString, compTk) {
					tokens = append(tokens, compTk)
				}
			}
			if len(tokens) == 1 {
				comp = tokens[0]
			}
			if len(tokens) > 0 {
				for _, tk := range tokens {
					if len(tk) > len(comp) {
						comp = tk
					}
				}
			}

			if comp != "" {
				valbytes := make([]byte, 0)
				record := false
				lastIndex := 0
				startIndex := 0
				for i, r := range valString {
					if !record && r == space {
						continue
					}
					if r == parenthesesLeft {
						record = true
						startIndex = i
						continue
					}
					if r == parenthesesRight {
						lastIndex = i
					}
					if record {
						valbytes = append(valbytes, byte(r))
					}
				}
				cutIndex := lastIndex - startIndex - 1
				if cutIndex < 0 {
					return fmt.Errorf("Struct[%s] Field[%s] evalid[%s] is invalid", "%s", "%s", exp)
				}
				valbytes = valbytes[:cutIndex]
				value = strings.TrimSpace(string(valbytes))
			}
		}

		e.Checkers = append(e.Checkers, &Checker{Exp: e, Method: method, Comp: comp, Value: value})
	}
	return nil
}

//NewLexer 获取一个表达式分词器
func NewLexer(msg, deli, pointer string, field *types.Var) *Lexer {
	return &Lexer{
		field:        field,
		message:      msg,
		delimiter:    deli,
		PointerAlias: pointer,
	}
}
