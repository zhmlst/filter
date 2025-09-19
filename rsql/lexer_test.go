package rsql

import "testing"

func lexAll(input string) []token {
	lxr := newLexer(input)
	var toks []token
	for {
		tok := lxr.Next()
		toks = append(toks, tok)
		if tok.Type == tokEOF || tok.Type == tokInvalid {
			break
		}
	}
	return toks
}

func compareTokens(a, b []token) bool {
	for i, t := range a {
		if t.Type != b[i].Type || t.Literal != b[i].Literal {
			return false
		}
	}
	return true
}

func TestLexer_Next(t *testing.T) {
	cases := []struct {
		input  string
		expect []token
	}{
		{
			`name==John`,
			[]token{
				{tokIdent, "name"},
				{tokEq, "=="},
				{tokIdent, "John"},
				{tokEOF, ""},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens := lexAll(c.input)
			if !compareTokens(tokens, c.expect) {
				t.Errorf("has:%v exp:%v", tokens, c.expect)
			} else {
				t.Logf("has:%v exp:%v", tokens, c.expect)
			}
		})
	}
}
