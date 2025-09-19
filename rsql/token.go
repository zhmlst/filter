package rsql

import (
	"fmt"
)

//go:generate stringer -type=tokenType -trimprefix=tok
type tokenType int32

const (
	_ tokenType = 1 << iota
	tokIdent
	tokString
	tokInteger
	tokFloat
	tokFalse
	tokTrue
	tokNull

	tokEq  // ==
	tokNe  // !=
	tokIn  // =in=
	tokOut // =out=
	tokLt  // =lt= or <
	tokLe  // =le= or <=
	tokGt  // =gt= or >
	tokGe  // =ge= or >=

	tokAnd // ;
	tokOr  // ,

	tokLparen
	tokRparen

	tokInvalid
	tokEOF
)

const ( // TokenType groups
	comparison = tokEq | tokNe | tokIn | tokOut | tokLt | tokLe | tokGt | tokGe
	equality   = tokEq | tokNe | tokIn | tokOut
	membership = tokIn | tokOut
	relation   = tokLt | tokLe | tokGt | tokGe
	argument   = tokIdent | tokString | tokInteger | tokFloat | tokFalse | tokTrue | tokNull
	logical    = tokAnd | tokOr
)

func (a tokenType) match(b tokenType) bool { return a&b != 0 }

type token struct {
	Type    tokenType
	Literal string
}

func (t token) String() string {
	return fmt.Sprintf("%s%#v", t.Type.String(), t.Literal)
}
