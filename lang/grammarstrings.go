package lang

import "fmt"

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

func (c Command) String() string {
	switch {
	case c.Scriptor != nil:
		return *c.Scriptor
	case c.Comment != nil:
		return c.Comment.String()
	case c.EOL != nil:
		return ""
	case c.Expression != nil:
		return c.Expression.String()
	case c.NamedFunction != nil:
		return c.NamedFunction.String()
	default:
		return fmt.Sprintf("*error in Command.String with %#v *", c)
	}

}

func (c Comment) String() string {
	return c.Comment
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
	case b.Tag != nil:
		return *b.Tag
	case b.Ident != nil:
		return *b.Ident
	case b.StringValue != nil:
		return fmt.Sprintf("%#v", *b.StringValue)
	case b.Subexpression != nil:
		return "(" + b.Subexpression.String() + ")"
	case b.List != nil:
		return b.List.String()
	case b.Invocation != nil:
		return b.Invocation.String()
	// case b.StatementBlock != nil:
	// 	return b.StatementBlock.String()
	case b.UnnamedFunction != nil:
		return b.UnnamedFunction.String()
	default:
		return fmt.Sprintf("*error in Base.String with %#v", b)
	}
}

func (f UnnamedFunction) String() string {
	s := "fn("
	for i, param := range f.Params {
		if i > 0 {
			s += ", "
		}
		s += param
	}
	s += ") " + f.Body.String()
	return s
}

func (f NamedFunction) String() string {
	s := "fn "
	if f.Name != nil {
		s += *f.Name
	}
	s += "("
	for i, param := range f.Params {
		if i > 0 {
			s += ", "
		}
		s += param
	}
	s += ") " + f.Body.String()
	return s
}

func (s StatementBlock) String() string {
	s1 := "{\n"
	for _, stmt := range s.Statements {
		s1 += "    " + stmt.String() + "\n"
	}
	return s1 + "}"
}

func (rb RequiredBlock) String() string {
	s1 := "{\n"
	for _, stmt := range rb.Statements {
		s1 += "    " + stmt.String() + "\n"
	}
	return s1 + "}"
}

func (e Expression) String() string {
	if e.Assignment == nil {
		return ""
	}
	return e.Assignment.String()
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

func (a Assignment) String() string {
	if a.Operation == nil {
		return a.Pipe.String()
	}
	return "(" + a.Pipe.String() + " " + a.Operation.Op + " " + a.Operation.Operand.String() + ")"
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

func (c Comparison) String() string {
	if len(c.Operations) == 0 {
		return c.Series.String()
	}
	s := "(" + c.Series.String()
	for _, op := range c.Operations {
		s += " " + op.Op + " " + op.Operand.String()
	}
	return s + ")"
}

func (s Series) String() string {
	str := ""
	if s.FromValue != nil {
		str += s.FromValue.String()
	}
	if s.ToValue != nil {
		str += " .. " + s.ToValue.String()
	}
	return str
}

func (kv KeyValue) String() string {
	s := kv.Addition.String()
	if kv.RightValue != nil {
		s += ": " + kv.RightValue.String()
	}
	return s
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

func (i Invocation) String() string {
	s := *i.Name
	s += "("
	for j, arg := range i.Arguments {
		if j > 0 {
			s += ", "
		}
		s += arg.String()
	}
	s += ")"
	return s
}
