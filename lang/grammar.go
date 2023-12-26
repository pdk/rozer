package lang

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Operator string

type Program struct {
	Pos lexer.Position

	Commands []*Command `parser:"@@*"`
}

// Command is basically a line, or a multi-line.
type Command struct {
	Pos lexer.Position

	Scriptor      *string        `parser:"  @Scriptor "`
	Comment       *Comment       `parser:"| @@ "`
	EOL           *string        `parser:"| @EOL "`
	NamedFunction *NamedFunction `parser:"| @@ "` // named functions only allowed at top level
	Expression    *Expression    `parser:"| @@ (EOL|EOF|(Comment (EOL|EOF))) "`
}

type Comment struct {
	Pos lexer.Position

	Comment string `parser:" @Comment (EOL|EOF) "`
}

type NamedFunction struct {
	Pos lexer.Position

	Func   *string        `parser:" @Function "`
	Name   *string        `parser:" @Ident "`
	Params []string       `parser:" '(' @Ident? (',' @Ident)* ')' "`
	Body   *RequiredBlock `parser:"@@"`
}

type Expression struct {
	Pos lexer.Position

	Assignment *Assignment `parser:"@@"`
}

type Assignment struct {
	Pos lexer.Position

	Pipe      *Pipe         `parser:"@@"`
	Operation *OpAssignment `parser:"@@?"`
}

type OpAssignment struct {
	Pos lexer.Position

	Op      string `parser:"@( ColonEqual | PlusEqual ) (EOL|Comment EOL)*"`
	Operand *Pipe  `parser:"@@"`
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

type Comparison struct {
	Pos lexer.Position

	Series     *Series         `parser:"@@"`
	Operations []*OpComparison `parser:"@@*"`
}

type OpComparison struct {
	Pos lexer.Position

	Op      string  `parser:"@( EqualEqual | BangEqual | LessEqual | MoreEqual | Less | More ) (EOL|Comment EOL)*"`
	Operand *Series `parser:"@@"`
}

// Series parses a n..m series.
type Series struct {
	Pos lexer.Position

	FromValue *KeyValue `parser:" @@ "`
	ToValue   *KeyValue `parser:" ( DotDot @@ )? "`
}

// KeyValue parses a key-value pair. If there is no colon, then it's just a value.
type KeyValue struct {
	Pos lexer.Position

	Addition   *Addition `parser:"@@"`
	RightValue *Addition `parser:"( Colon (EOL|Comment EOL)* @@ )?"`
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

type Unary struct {
	Pos lexer.Position

	Op    *string `parser:"  ( @( Bang | Minus )"`
	Unary *Unary  `parser:"    @@ )"`
	Base  *Base   `parser:"| @@"`
}

type Base struct {
	Pos lexer.Position

	Subexpression   *Expression      `parser:"  '(' @@ ')' "`
	List            *List            `parser:"| @@"`
	UnnamedFunction *UnnamedFunction `parser:"| @@ "`
	Invocation      *Invocation      `parser:"| @@ "`
	StringValue     *string          `parser:"| @String "`
	Bool            *string          `parser:"| @('true' | 'false')"`
	Tag             *string          `parser:"| @Tag"`
	Ident           *string          `parser:"| @Ident "`
	// DateTime      *string      `parser:"| @DateTime"`
	// Date          *string      `parser:"| @Date"`
	// Time          *string      `parser:"| @Time"`
	Float   *float64 `parser:"| @Float"`
	Integer *int64   `parser:"| @Integer"`
	// TimeSpan      *string      `parser:"| @TimeSpan"`
	// DottedIdent   *DottedIdent `parser:"| @@"`
	// StatementBlock  *StatementBlock  `parser:"| '{' (Comment EOL|EOL)* @@ (Comment EOL|EOL)* '}' (EOF|EOL|Comment EOL)* "`
}

type UnnamedFunction struct {
	Pos lexer.Position

	Func   *string        `parser:" Function "`
	Params []string       `parser:" '(' @Ident? (',' @Ident)* ')' "`
	Body   *RequiredBlock `parser:"@@"`
}

type Invocation struct {
	Pos lexer.Position

	Name      *string       `parser:" @Ident "`
	Arguments []*Expression `parser:" '(' @@? ( ',' @@ )* ')' "`
}

type RequiredBlock struct {
	Pos lexer.Position

	LeftBrace  *string       `parser:" @'{' "`
	Statements []*Expression `parser:" (@@ | Comment EOL | EOL)* "`
	RightBrace *string       `parser:" @'}' "`
}

type StatementBlock struct {
	Pos lexer.Position

	Statements []*Expression `parser:" (@@ | Comment EOL | EOL)* "`
}

type List struct {
	Pos lexer.Position

	Items []*Expression `parser:"'[' (EOL|Comment EOL)* @@? ( ',' (EOL|Comment EOL)* @@ )* (EOL|Comment EOL)* ']' "`
}
