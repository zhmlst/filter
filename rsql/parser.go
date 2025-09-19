package rsql

import (
	"fmt"
	"strconv"

	"github.com/zhmlst/filter"
)

type parser struct {
	lexer *lexer
	curr  token
	next  token
}

func (p *parser) readToken() {
	p.curr = p.next
	p.next = p.lexer.Next()
}

func (p *parser) parseLiteral() (any, error) {
	switch p.curr.Type {
	case tokIdent, tokString:
		v := p.curr.Literal
		p.readToken()
		return v, nil
	case tokInteger:
		i, err := strconv.ParseInt(p.curr.Literal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %w", p.curr.Literal, err)
		}
		p.readToken()
		return i, nil
	case tokFloat:
		f, err := strconv.ParseFloat(p.curr.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q: %w", p.curr.Literal, err)
		}
		p.readToken()
		return f, nil
	case tokTrue:
		p.readToken()
		return true, nil
	case tokFalse:
		p.readToken()
		return false, nil
	case tokNull:
		p.readToken()
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected value token: %s", p.curr.String())
	}
}

func (p *parser) expect(tt tokenType) error {
	if p.curr.Type != tt {
		return fmt.Errorf("unexpected token %v, expected %v", p.curr, tt)
	}
	return nil
}

func tokenTypeToComparisonOp(t tokenType) (filter.ComparisonOp, error) {
	switch t {
	case tokEq:
		return filter.Eq, nil
	case tokNe:
		return filter.Ne, nil
	case tokIn:
		return filter.In, nil
	case tokOut:
		return filter.Out, nil
	case tokLt:
		return filter.Lt, nil
	case tokLe:
		return filter.Le, nil
	case tokGt:
		return filter.Gt, nil
	case tokGe:
		return filter.Ge, nil
	default:
		return 0, fmt.Errorf("not a comparison token: %v", t)
	}
}

func (p *parser) parsePrimary() (filter.Node, error) {
	switch p.curr.Type {
	case tokLparen:
		p.readToken()
		node, err := p.parseExpression(precLowest)
		if err != nil {
			return nil, err
		}
		if err := p.expect(tokRparen); err != nil {
			return nil, fmt.Errorf("expected ')': %w", err)
		}
		p.readToken()
		return node, nil
	case tokIdent:
		field := p.curr.Literal
		p.readToken()
		if !p.curr.Type.match(comparison) {
			return nil, fmt.Errorf("expected comparison after field %q, got %s", field, p.curr.String())
		}
		compTok := p.curr.Type
		compOp, err := tokenTypeToComparisonOp(compTok)
		if err != nil {
			return nil, err
		}
		p.readToken()

		if compTok.match(membership) {
			if p.curr.Type != tokLparen {
				return nil, fmt.Errorf("expected '(' after membership operator for field %q, got %s", field, p.curr.String())
			}
			p.readToken()
			var vals []any
			if p.curr.Type == tokRparen {
				return nil, fmt.Errorf("empty list for membership operator on field %q", field)
			}
			for {
				val, err := p.parseLiteral()
				if err != nil {
					return nil, err
				}
				vals = append(vals, val)

				if p.curr.Type == tokRparen {
					p.readToken()
					break
				}
				if p.curr.Type != tokOr {
					return nil, fmt.Errorf("expected ',' between membership values, got %s", p.curr.String())
				}
				p.readToken()
			}
			return filter.Constraint{Field: field, Operator: compOp, Value: vals}, nil
		}

		if !p.curr.Type.match(argument) {
			return nil, fmt.Errorf("expected argument after comparison for field %q, got %s", field, p.curr.String())
		}
		val, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		return filter.Constraint{Field: field, Operator: compOp, Value: val}, nil
	default:
		return nil, fmt.Errorf("unexpected primary token: %s", p.curr.String())
	}
}

func (p *parser) parseExpression(precedence int) (filter.Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for precedence < precedenceOf(p.curr.Type) {
		opTok := p.curr.Type
		if !opTok.match(logical) {
			break
		}
		logOp, err := tokenTypeToLogicalOp(opTok)
		if err != nil {
			return nil, err
		}
		p.readToken()
		right, err := p.parseExpression(precedenceOf(opTok))
		if err != nil {
			return nil, err
		}

		var nodes []filter.Node
		if l, ok := left.(filter.Logical); ok && l.Operator == logOp {
			nodes = append(nodes, l.Nodes...)
		} else {
			nodes = append(nodes, left)
		}
		if r, ok := right.(filter.Logical); ok && r.Operator == logOp {
			nodes = append(nodes, r.Nodes...)
		} else {
			nodes = append(nodes, right)
		}
		left = filter.Logical{Operator: logOp, Nodes: nodes}
	}

	return left, nil
}

func tokenTypeToLogicalOp(t tokenType) (filter.LogicalOp, error) {
	switch t {
	case tokAnd:
		return filter.And, nil
	case tokOr:
		return filter.Or, nil
	default:
		return 0, fmt.Errorf("not a logical token: %v", t)
	}
}

const (
	precLowest = iota
	precOr
	precAnd
	precCmp
)

func precedenceOf(t tokenType) int {
	switch {
	case t == tokOr:
		return precOr
	case t == tokAnd:
		return precAnd
	case t.match(comparison):
		return precCmp
	default:
		return precLowest
	}
}

func Parse(rsql string) (filter.Node, error) {
	p := &parser{lexer: newLexer(rsql)}
	p.readToken()
	p.readToken()

	node, err := p.parseExpression(precLowest)
	if err != nil {
		return nil, err
	}

	if p.curr.Type != tokEOF {
		return nil, fmt.Errorf("unexpected token after expression: %s", p.curr.String())
	}
	return node, nil
}
