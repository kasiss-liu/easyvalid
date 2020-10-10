package lexer

import "testing"

func TestLexer(t *testing.T) {
	l := &Lexer{
		message:   "len:eq(1)|notEmpty|msg:st;",
		delimiter: ";",
	}
	l.splitExpr()
	t.Log(l.exprStrings)
	l.splitValid()
	t.Log(l.exprs)
	l.parseExpress()
	for _, v := range l.exprs {
		t.Logf("%#v\n", v.Checkers[0])
		t.Logf("%#v\n", v.Checkers[1])
	}
}
