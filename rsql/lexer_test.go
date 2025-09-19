package rsql

import (
	"reflect"
	"testing"
)

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

	tokens := func(src string) (tokens []token) {
		lxr := newLexer(src)
		for t := lxr.Next(); ; t = lxr.Next() {
			tokens = append(tokens, t)
			if t.Type == tokEOF || t.Type == tokInvalid {
				break
			}
		}
		return
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens := tokens(c.input)
			if !reflect.DeepEqual(tokens, c.expect) {
				t.Errorf("\ngot: %v\nexp: %v", tokens, c.expect)
			}
		})
	}
}
