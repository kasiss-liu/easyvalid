package lexer

import "testing"

func TestInValidator(t *testing.T) {
	l := &Lexer{
		message:   "in:(1,2,3)|msg:st;",
		delimiter: ";",
	}
	l.splitExpr()
	t.Log(l.exprStrings)
	l.splitValid()
	t.Log(l.exprs)
	l.parseExpress()
	for _, v := range l.exprs {
		t.Logf("%#v\n", v.Checkers[0])
	}
}
