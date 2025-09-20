package ast

import (
	"fmt"
	"reflect"
	"strings"
)

/* --Operators-- */

//go:generate stringer -type=BoolOp
type BoolOp int

//go:generate stringer -type=CompOp
type CompOp int

const (
	And BoolOp = iota
	Or
)

const (
	Eq CompOp = iota
	Ne
	In
	Out
	Lt
	Le
	Gt
	Ge
)

/* --Nodes-- */

type Node interface {
	node()
	fmt.Stringer
}

type BoolNode struct {
	Op   BoolOp
	Args []Node
}

type CompNode struct {
	Field string
	Op    CompOp
	Arg   any
}

/* --BoolNode-- */
func (BoolNode) node()            {}
func (b BoolNode) String() string { return b.formatIndented(0) }
func (b BoolNode) formatIndented(indent int) string {
	ind := strings.Repeat("\t", indent)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%s [\n", ind, b.Op.String()))
	for i, node := range b.Args {
		switch n := node.(type) {
		case interface{ formatIndented(int) string }:
			sb.WriteString(n.formatIndented(indent + 1))
		default:
			lines := strings.Split(node.String(), "\n")
			for j, ln := range lines {
				if ln == "" && j == len(lines)-1 {
					continue
				}
				sb.WriteString(strings.Repeat("\t", indent+1) + ln)
				if j < len(lines)-1 {
					sb.WriteString("\n")
				}
			}
		}
		if i < len(b.Args)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(ind + "]")
	return sb.String()
}

/* --CompNode-- */
func (CompNode) node()            {}
func (c CompNode) String() string { return c.formatIndented(0) }
func (c CompNode) formatIndented(indent int) string {
	ind := strings.Repeat("\t", indent)
	val := func() string {
		v := reflect.ValueOf(c.Arg)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			var elems []string
			for i := 0; i < v.Len(); i++ {
				elems = append(elems, fmt.Sprintf("%#v", v.Index(i).Interface()))
			}
			return "[" + strings.Join(elems, ", ") + "]"
		default:
			return fmt.Sprintf("%#v", c.Arg)
		}
	}()
	return fmt.Sprintf("%s%s %s %s", ind, c.Field, c.Op.String(), val)
}
