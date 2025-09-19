package filter

import (
	"fmt"
	"reflect"
	"strings"
)

//go:generate stringer -type=LogicalOp
type LogicalOp int

//go:generate stringer -type=ComparisonOp -trimprefix=ComparisonOp
type ComparisonOp int

const (
	LogicalOpAnd LogicalOp = iota
	LogicalOpOr
)

const (
	ComparisonOpEq ComparisonOp = iota
	ComparisonOpNe
	ComparisonOpIn
	ComparisonOpOut
	ComparisonOpLt
	ComparisonOpLe
	ComparisonOpGt
	ComparisonOpGe
)

type Node interface {
	node()
	fmt.Stringer
}

type Logical struct {
	Operator LogicalOp
	Nodes    []Node
}

func (Logical) node() {}

func (l Logical) String() string { return l.formatIndented(0) }

func (l Logical) formatIndented(indent int) string {
	ind := strings.Repeat("\t", indent)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s%s [\n", ind, l.Operator.String()))
	for i, node := range l.Nodes {
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
		if i < len(l.Nodes)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString(ind + "]")
	return sb.String()
}

type Constraint struct {
	Field    string
	Operator ComparisonOp
	Value    any
}

func (Constraint) node() {}

func (c Constraint) String() string { return c.formatIndented(0) }

func (c Constraint) formatIndented(indent int) string {
	ind := strings.Repeat("\t", indent)
	val := func() string {
		v := reflect.ValueOf(c.Value)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			var elems []string
			for i := 0; i < v.Len(); i++ {
				elems = append(elems, fmt.Sprintf("%#v", v.Index(i).Interface()))
			}
			return "[" + strings.Join(elems, ", ") + "]"
		default:
			return fmt.Sprintf("%#v", c.Value)
		}
	}()
	return fmt.Sprintf("%s%s %s %s", ind, c.Field, c.Operator.String(), val)
}

func Same(n1, n2 Node) bool {
	if v1, ok := n1.(Logical); ok {
		n1 = &v1
	}
	if v2, ok := n2.(Logical); ok {
		n2 = &v2
	}
	if n1 == nil || n2 == nil {
		return n1 == n2
	}
	switch v1 := n1.(type) {
	case *Constraint:
		if v2, ok := n2.(*Constraint); ok {
			return v1.Field == v2.Field && v1.Operator == v2.Operator && v1.Value == v2.Value
		}
		return false
	case *Logical:
		if v2, ok := n2.(*Logical); ok {
			if len(v1.Nodes) != len(v2.Nodes) {
				return false
			}
			for i := range v1.Nodes {
				if !Same(v1.Nodes[i], v2.Nodes[i]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		return false
	}
}
