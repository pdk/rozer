package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2/lexer"
)

type Operator string

type Program struct {
	Pos lexer.Position

	Commands []*Command `parser:"@@*"`
}

func (p Program) String() string {
	s := ""
	for i, c := range p.Commands {
		s1 := c.String()
		if s1 != "" {
			s += fmt.Sprintf("%3d: %s\n", i, s1)
		}
	}
	return s
}

// Command is basically a line, or a multi-line.
type Command struct {
	Pos lexer.Position

	Comment *Comment `parser:"  @@"`
	EOL     *string  `parser:"| @EOL"`
	// Func       *Func       `parser:"| @@"`
	Expression *Expression `parser:"| @@ (EOL|EOF|(Comment (EOL|EOF)))"`
}

func (c Command) String() string {
	switch {
	case c.Comment != nil:
		return c.Comment.String()
	case c.EOL != nil:
		return ""
	case c.Expression != nil:
		return c.Expression.String()
	default:
		return fmt.Sprintf("*error in Command.String with %#v *", c)
	}

}

type Comment struct {
	Pos lexer.Position

	Comment string `parser:"@Comment"`
}

func (c Comment) String() string {
	return c.Comment
}

type Base struct {
	Pos lexer.Position

	Bool *string `parser:"  @('true' | 'false')"`
	// DateTime      *string      `parser:"| @DateTime"`
	// Date          *string      `parser:"| @Date"`
	// Time          *string      `parser:"| @Time"`
	Float   *float64 `parser:"| @Float"`
	Integer *int64   `parser:"| @Integer"`
	// TimeSpan      *string      `parser:"| @TimeSpan"`
	// Invocation    *Invocation  `parser:"| @@"`
	// DottedIdent   *DottedIdent `parser:"| @@"`
	Ident         *string     `parser:"| @Ident"`
	StringValue   *string     `parser:"| @String"`
	Subexpression *Expression `parser:"| '(' @@ ')' "`
	List          *List       `parser:"| @@"`
	// Lambda        *Lambda      `parser:"| @@"`
}

func (b Base) String() string {
	switch {
	case b.Bool != nil:
		return *b.Bool
	// case b.DateTime != nil:
	// 	return *b.DateTime
	// case b.Date != nil:
	// 	return *b.Date
	// case b.Time != nil:
	// 	return *b.Time
	case b.Float != nil:
		return fmt.Sprintf("%.20f", *b.Float)
	case b.Integer != nil:
		return fmt.Sprintf("%d", *b.Integer)
	// case b.TimeSpan != nil:
	// 	return *b.TimeSpan
	// case b.DottedIdent != nil:
	// 	return b.DottedIdent.String()
	case b.Ident != nil:
		return *b.Ident
	case b.StringValue != nil:
		return fmt.Sprintf("%#v", *b.StringValue)
	case b.Subexpression != nil:
		return "(" + b.Subexpression.String() + ")"
	case b.List != nil:
		return b.List.String()
	// case b.Lambda != nil:
	// 	return b.Lambdb.String()
	// case b.Invocation != nil:
	// 	return b.Invocation.String()
	default:
		return fmt.Sprintf("*error in Base.String with %#v", b)
	}
}

type Expression struct {
	Pos lexer.Position

	Assignment *Assignment `parser:"@@"`
}

func (e Expression) String() string {
	if e.Assignment == nil {
		return ""
	}
	return e.Assignment.String()
}

type List struct {
	Pos lexer.Position

	Items []*Expression `parser:"'[' (EOL|Comment EOL)* @@? ( ',' (EOL|Comment EOL)* @@ )* (EOL|Comment EOL)* ']' "`
}

func (l List) String() string {
	s := "["
	for i, item := range l.Items {
		if i > 0 {
			s += ", "
		}
		s += item.String()
	}
	return s + "]"
}

type Assignment struct {
	Pos lexer.Position

	Pipe       *Pipe           `parser:"@@"`
	Operations []*OpAssignment `parser:"@@*"`
}

type OpAssignment struct {
	Pos lexer.Position

	Op      string `parser:"@( ColonEqual | PlusEqual ) (EOL|Comment EOL)*"`
	Operand *Pipe  `parser:"@@"`
}

func (a Assignment) String() string {
	if len(a.Operations) == 0 {
		return a.Pipe.String()
	}
	s := "(" + a.Pipe.String()
	for _, op := range a.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

type Pipe struct {
	Pos lexer.Position

	Logical    *Logical  `parser:"@@"`
	Operations []*OpPipe `parser:"@@*"`
}

type OpPipe struct {
	Pos lexer.Position

	Op      string   `parser:"@( MoreMore | MoreMoreMore ) (EOL|Comment EOL)*"`
	Operand *Logical `parser:"@@"`
}

func (p Pipe) String() string {
	if len(p.Operations) == 0 {
		return p.Logical.String()
	}
	s := "(" + p.Logical.String()
	for _, op := range p.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

type Logical struct {
	Pos lexer.Position

	Comparison *Comparison  `parser:"@@"`
	Operations []*OpLogical `parser:"@@*"`
}

type OpLogical struct {
	Pos lexer.Position

	Op      string      `parser:"@( AndAnd | OrOr ) (EOL|Comment EOL)*"`
	Operand *Comparison `parser:"@@"`
}

func (l Logical) String() string {
	if len(l.Operations) == 0 {
		return l.Comparison.String()
	}
	s := "(" + l.Comparison.String()
	for _, op := range l.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

type Comparison struct {
	Pos lexer.Position

	KeyValue   *KeyValue       `parser:"@@"`
	Operations []*OpComparison `parser:"@@*"`
}

type OpComparison struct {
	Pos lexer.Position

	Op      string    `parser:"@( EqualEqual | BangEqual | LessEqual | MoreEqual | Less | More ) (EOL|Comment EOL)*"`
	Operand *KeyValue `parser:"@@"`
}

func (c Comparison) String() string {
	if len(c.Operations) == 0 {
		return c.KeyValue.String()
	}
	s := "(" + c.KeyValue.String()
	for _, op := range c.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

// KeyValue parses a key-value pair. If there is no colon, then it's just a value.
type KeyValue struct {
	Pos lexer.Position

	Addition   *Addition `parser:"@@"`
	RightValue *Addition `parser:"( Colon (EOL|Comment EOL)* @@ )?"`
}

func (kv KeyValue) String() string {
	s := kv.Addition.String()
	if kv.RightValue != nil {
		s += ": " + kv.RightValue.String()
	}
	return s
}

type Addition struct {
	Pos lexer.Position

	Multiplication *Multiplication `parser:"@@"`
	Operations     []*OpAddition   `parser:"@@*"`
}

type OpAddition struct {
	Pos lexer.Position

	Op      string          `parser:"@( Plus | Minus ) (EOL|Comment EOL)*"`
	Operand *Multiplication `parser:"@@"`
}

func (a Addition) String() string {
	if len(a.Operations) == 0 {
		return a.Multiplication.String()
	}
	s := "(" + a.Multiplication.String()
	for _, op := range a.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

type Multiplication struct {
	Pos lexer.Position

	Unary      *Unary              `parser:"@@"`
	Operations []*OpMultiplication `parser:"@@*"`
}

type OpMultiplication struct {
	Pos lexer.Position

	Op      string `parser:"@( Star | Slash | Percent ) (EOL|Comment EOL)*"`
	Operand *Unary `parser:"@@"`
}

func (x Multiplication) String() string {
	if len(x.Operations) == 0 {
		return x.Unary.String()
	}
	s := "(" + x.Unary.String()
	for _, op := range x.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

type Unary struct {
	Pos lexer.Position

	Op    *string `parser:"  ( @( Bang | Minus )"`
	Unary *Unary  `parser:"    @@ )"`
	Base  *Base   `parser:"| @@"`
}

func (x Unary) String() string {
	s := ""
	if x.Op != nil {
		s += *x.Op
	}
	if x.Unary != nil {
		s += x.Unary.String()
	}
	if x.Base != nil {
		s += x.Base.String()
	}
	return s
}
