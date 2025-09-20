package rsql

import (
	"reflect"
	"testing"

	"github.com/zhmlst/filter/ast"
)

func TestRsql_Parse(t *testing.T) {
	cases := []struct {
		input string
		exp   ast.Node
	}{
		{
			`name!=John`,
			ast.CompNode{
				Field: "name",
				Op:    ast.Ne,
				Arg:   "John",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got, err := Parse(c.input)
			if err != nil {
				t.Fatal(err.Error())
				return
			}
			if !reflect.DeepEqual(got, c.exp) {
				t.Errorf("\ngot: %v, \nexp: %v", got, c.exp)
			}
		})
	}
}
