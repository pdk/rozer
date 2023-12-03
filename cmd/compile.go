package main

import (
	"fmt"
	"log"
)

type ExecutionEnvironment map[string]any
type ExecutionResult any

type Executable interface {
	Execute(ExecutionEnvironment) ExecutionResult
	Type() Type
}

type Type uint

const (
	TypeUnknown Type = iota
	TypeBool
	TypeFloat
	TypeInteger
	TypeString
	TypeList
	TypeKeyValue
)

func (t Type) String() string {
	switch t {
	case TypeUnknown:
		return "unknown"
	case TypeBool:
		return "bool"
	case TypeFloat:
		return "float"
	case TypeInteger:
		return "integer"
	case TypeString:
		return "string"
	case TypeList:
		return "list"
	case TypeKeyValue:
		return "keyvalue"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

type CompileErrors struct {
	errs *[]error
}

var (
	NoErrors = CompileErrors{}
)

func (ce *CompileErrors) Append(err ...error) CompileErrors {
	if ce.errs == nil {
		ce.errs = &[]error{}
	}

	*ce.errs = append(*ce.errs, err...)

	return *ce
}

func (ce *CompileErrors) Collect(e Executable, errs CompileErrors) Executable {
	if errs.errs != nil {
		ce.Append(*errs.errs...)
	}
	return e
}

func NewError(err error) CompileErrors {
	return CompileErrors{&[]error{err}}
}

func (ce CompileErrors) Len() int {
	if ce.errs == nil {
		return 0
	}
	return len(*ce.errs)
}

type Block struct {
	Commands []Executable
}

func (b Block) Execute(ee ExecutionEnvironment) ExecutionResult {
	var lastResult ExecutionResult
	for _, c := range b.Commands {
		lastResult = c.Execute(ee)
		if lastResult != nil {
			log.Printf("execution returned %v", lastResult)
		} else {
			log.Printf("execution returned nil")
		}
	}
	return lastResult
}

func (b Block) Type() Type {
	return TypeUnknown
}

type BoolValue bool

func (b BoolValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("bool: %t", b)
	return b
}

func (b BoolValue) Type() Type {
	return TypeBool
}

type FloatValue float64

func (f FloatValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("float: %.20f", f)
	return f
}

func (f FloatValue) Type() Type {
	return TypeFloat
}

type IntegerValue int64

func (i IntegerValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("integer: %d", i)
	return i
}

func (i IntegerValue) Type() Type {
	return TypeInteger
}

type StringValue string

func (s StringValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("string: %#v", s)
	return s
}

func (s StringValue) Type() Type {
	return TypeString
}

type ListValue []Executable

func (l ListValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("list: %v", l)
	return l
}

func (l ListValue) Type() Type {
	return TypeList
}

func (b Base) Compile() (Executable, CompileErrors) {
	switch {
	case b.Bool != nil:
		switch *b.Bool {
		case "true":
			return BoolValue(true), NoErrors
		case "false":
			return BoolValue(false), NoErrors
		default:
			log.Fatalf("invalid bool at %s: %s", b.Pos, *b.Bool)
		}
	case b.Float != nil:
		return FloatValue(*b.Float), NoErrors
	case b.Integer != nil:
		return IntegerValue(*b.Integer), NoErrors
	case b.Ident != nil:
		return StringValue(*b.Ident), NoErrors // TODO lookup in environment
	case b.StringValue != nil:
		return StringValue(*b.StringValue), NoErrors
	case b.Subexpression != nil:
		return b.Subexpression.Compile()
	case b.List != nil:
		return b.List.Compile()
	}

	return nil, NewError(fmt.Errorf("cannot compile base %#v", b))
}

type UnaryExecuteNot struct {
	Unary   *Unary
	Operand Executable
}

func (uen UnaryExecuteNot) Execute(ee ExecutionEnvironment) ExecutionResult {
	return BoolValue(!uen.Operand.Execute(ee).(BoolValue))
}

func (uen UnaryExecuteNot) Type() Type {
	return TypeBool
}

type UnaryExecuteSubtractFloat struct {
	Unary   *Unary
	Operand Executable
}

func (uemf UnaryExecuteSubtractFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	return FloatValue(-uemf.Operand.Execute(ee).(FloatValue))
}

func (uemf UnaryExecuteSubtractFloat) Type() Type {
	return TypeFloat
}

type UnaryExecuteMinusInteger struct {
	Unary   *Unary
	Operand Executable
}

func (uemi UnaryExecuteMinusInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(-uemi.Operand.Execute(ee).(IntegerValue))
}

func (uemi UnaryExecuteMinusInteger) Type() Type {
	return TypeInteger
}

func (u *Unary) Compile() (Executable, CompileErrors) {

	if u.Base != nil {
		return u.Base.Compile()
	}

	if u.Unary == nil {
		return nil, NewError(fmt.Errorf("cannot compile unary %s", u))
	}

	operand, errs := u.Unary.Compile()

	switch *u.Op {
	case "!":
		if operand.Type() != TypeBool {
			errs.Append(fmt.Errorf("invalid unary operation %s at %s", *u.Op, u.Pos))
		}
		return UnaryExecuteNot{u, operand}, errs
	case "-":
		switch operand.Type() {
		case TypeFloat:
			return UnaryExecuteSubtractFloat{u, operand}, errs
		case TypeInteger:
			return UnaryExecuteMinusInteger{u, operand}, errs
		default:
			errs.Append(fmt.Errorf("invalid unary operation %s at %s", *u.Op, u.Pos))
			return nil, errs
		}
	}

	return nil, errs.Append(fmt.Errorf("invalid unary operation %s at %s", *u.Op, u.Pos))
}

type MultiplyFloat struct {
	Multiplication *Multiplication
	Left, Right    Executable
}

func (mf MultiplyFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	return FloatValue(mf.Left.Execute(ee).(FloatValue) * mf.Right.Execute(ee).(FloatValue))
}

func (mf MultiplyFloat) Type() Type {
	return TypeFloat
}

type MultiplyInteger struct {
	Multiplication *Multiplication
	Left, Right    Executable
}

func (mi MultiplyInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(mi.Left.Execute(ee).(IntegerValue) * mi.Right.Execute(ee).(IntegerValue))
}

func (mi MultiplyInteger) Type() Type {
	return TypeInteger
}

type DivideFloat struct {
	Multiplication *Multiplication
	Left, Right    Executable
}

func (df DivideFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	left := df.Left.Execute(ee).(FloatValue)
	right := df.Right.Execute(ee).(FloatValue)
	if right == 0.0 {
		log.Printf("float divide by zero at %s (returning 0.0)", df.Multiplication.Pos)
		return FloatValue(0.0)
	}
	return FloatValue(left / right)
}

func (df DivideFloat) Type() Type {
	return TypeFloat
}

type DivideInteger struct {
	Multiplication *Multiplication
	Left, Right    Executable
}

func (di DivideInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	left := di.Left.Execute(ee).(IntegerValue)
	right := di.Right.Execute(ee).(IntegerValue)
	if right == 0 {
		log.Printf("integer divide by zero at %s (returning 0)", di.Multiplication.Pos)
		return IntegerValue(0.0)
	}
	return IntegerValue(left / right)
}

func (di DivideInteger) Type() Type {
	return TypeInteger
}

type ModuloInteger struct {
	Multiplication *Multiplication
	Left, Right    Executable
}

func (mi ModuloInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(mi.Left.Execute(ee).(IntegerValue) % mi.Right.Execute(ee).(IntegerValue))
}

func (mi ModuloInteger) Type() Type {
	return TypeInteger
}

func (m *Multiplication) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(m.Unary.Compile())
	if len(m.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range m.Operations {
		operand := errs.Collect(opExpr.Operand.Compile())
		if ex.Type() != operand.Type() {
			errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(), operand.Type(), m.Pos))
		}
		operands = append(operands, operand)
	}

	switch ex.Type() {
	case TypeFloat:
		for i, operand := range operands {
			switch m.Operations[i].Op {
			case "*":
				ex = MultiplyFloat{m, ex, operand}
			case "/":
				ex = DivideFloat{m, ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid type %s for *, /, %% at %s", ex.Type(), m.Pos))
				return nil, errs
			}
		}

	case TypeInteger:
		for i, operand := range operands {
			switch m.Operations[i].Op {
			case "*":
				ex = MultiplyInteger{m, ex, operand}
			case "/":
				ex = DivideInteger{m, ex, operand}
			case "%":
				ex = ModuloInteger{m, ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid type %s for *, /, %% at %s", ex.Type(), m.Pos))
				return nil, errs
			}
		}

	default:
		errs.Append(fmt.Errorf("invalid type %s for *, /, %% at %s", ex.Type(), m.Pos))
		return nil, errs
	}

	return ex, errs
}

func (kv *KeyValue) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(kv.Addition.Compile())
	if kv.RightValue == nil {
		return ex, errs
	}

	if ex.Type() != TypeString {
		errs.Append(fmt.Errorf("invalid type %s for key (should be string) at %s", ex.Type(), kv.Pos))
	}

	rex := errs.Collect(kv.RightValue.Compile())
	return KeyValueExecute{kv, ex, rex}, errs
}

type KeyValueExecute struct {
	KeyValue *KeyValue

	Key   Executable
	Value Executable
}

type KeyValueResult struct {
	Key   any
	Value any
}

func (kve KeyValueExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	return KeyValueResult{kve.Key.Execute(ee), kve.Value.Execute(ee)}
}

func (kve KeyValueExecute) Type() Type {
	return TypeKeyValue
}

func (l *Logical) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(l.Comparison.Compile())
	if len(l.Operations) == 0 {
		return ex, errs
	}

	if ex.Type() != TypeBool {
		errs.Append(fmt.Errorf("invalid type %s for logical operation at %s", ex.Type(), l.Pos))
	}

	operands := []Executable{}
	for _, opExpr := range l.Operations {
		operands = append(operands, errs.Collect(opExpr.Operand.Compile()))
	}

	for i, operand := range operands {
		switch ex.Type() {
		case TypeBool:
			switch l.Operations[i].Op {
			case "&&":
				ex = Comparator{nil, ex, operand,
					func(left, right ExecutionResult) bool {
						return bool(left.(BoolValue)) && bool(right.(BoolValue))
					}}
			case "||":
				ex = Comparator{nil, ex, operand,
					func(left, right ExecutionResult) bool {
						return bool(left.(BoolValue)) || bool(right.(BoolValue))
					}}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", l.Operations[i].Op, ex.Type(), l.Pos))
				return nil, errs
			}
		default:
			errs.Append(fmt.Errorf("invalid type %s for logical operation at %s", ex.Type(), l.Pos))
			return nil, errs
		}
	}

	return ex, errs
}

func (c *Comparison) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(c.KeyValue.Compile())
	if len(c.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range c.Operations {
		operands = append(operands, errs.Collect(opExpr.Operand.Compile()))
	}

	for i, operand := range operands {

		if ex.Type() != operand.Type() {
			errs.Append(fmt.Errorf("type mismatch %s / %s at %s", ex.Type(), operand.Type(), c.Operations[i].Pos))
		}

		switch ex.Type() {
		case TypeBool:
			switch c.Operations[i].Op {
			case "==":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(BoolValue) == right.(BoolValue)
					}}
			case "!=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(BoolValue) != right.(BoolValue)
					}}
			}

		case TypeFloat:
			switch c.Operations[i].Op {
			case "<":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) < right.(FloatValue)
					}}
			case "<=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) <= right.(FloatValue)
					}}
			case ">":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) > right.(FloatValue)
					}}
			case ">=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) >= right.(FloatValue)
					}}
			case "==":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) == right.(FloatValue)
					}}
			case "!=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(FloatValue) != right.(FloatValue)
					}}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(), c.Pos))
				return nil, errs
			}

		case TypeInteger:
			switch c.Operations[i].Op {
			case "<":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) < right.(IntegerValue)
					}}
			case "<=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) <= right.(IntegerValue)
					}}
			case ">":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) > right.(IntegerValue)
					}}
			case ">=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) >= right.(IntegerValue)
					}}
			case "==":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) == right.(IntegerValue)
					}}
			case "!=":
				ex = Comparator{c, ex, operand,
					func(left, right ExecutionResult) bool {
						return left.(IntegerValue) != right.(IntegerValue)
					}}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(), c.Pos))
				return nil, errs
			}

		case TypeString:
			for i, operand := range operands {
				switch c.Operations[i].Op {
				case "<":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) < right.(StringValue)
						},
					}
				case "<=":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) <= right.(StringValue)
						}}
				case ">":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) > right.(StringValue)
						}}
				case ">=":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) >= right.(StringValue)
						}}
				case "==":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) == right.(StringValue)
						}}
				case "!=":
					ex = Comparator{c, ex, operand,
						func(left, right ExecutionResult) bool {
							return left.(StringValue) != right.(StringValue)
						}}
				default:
					errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(), c.Pos))
					return nil, errs
				}
			}
		default:
			errs.Append(fmt.Errorf("invalid type %s for comparison at %s", ex.Type(), c.Pos))
		}
	}
	return ex, errs
}

type Comparator struct {
	Comparison  *Comparison
	Left, Right Executable
	Executor    func(ExecutionResult, ExecutionResult) bool
}

func (c Comparator) Execute(ee ExecutionEnvironment) ExecutionResult {
	left := c.Left.Execute(ee)
	right := c.Right.Execute(ee)
	result := c.Executor(left, right)
	return BoolValue(result)
}

func (c Comparator) Type() Type {
	return TypeBool
}

func (p *Pipe) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(p.Logical.Compile())
	if len(p.Operations) == 0 {
		return ex, errs
	}

	pe := PipeExecute{p, []Executable{ex}}
	for _, op := range p.Operations {
		pe.Commands = append(pe.Commands, errs.Collect(op.Operand.Compile()))
	}

	return pe, errs
}

type PipeExecute struct {
	Pipe *Pipe

	Commands []Executable
}

func (pe PipeExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	var lastResult ExecutionResult
	for _, c := range pe.Commands {
		lastResult = c.Execute(ee)
	}
	return lastResult
}

func (pe PipeExecute) Type() Type {
	return TypeUnknown
}

func (a *Assignment) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(a.Pipe.Compile())
	if len(a.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range a.Operations {
		operands = append(operands, errs.Collect(opExpr.Operand.Compile()))
	}

	for i, operand := range operands {
		switch a.Operations[i].Op {
		case ":=":
			ex = AssignmentExecute{a, ex, operand}
		case "+=":
			ex = PlusAssignmentExecute{a, ex, operand}
		}
	}

	return ex, errs
}

type AssignmentExecute struct {
	Assignment *Assignment

	Left, Right Executable
}

func (ae AssignmentExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	// ee[ae.Left.Execute(ee).(StringValue)] = right
	right := ae.Right.Execute(ee)
	return right
}

func (ae AssignmentExecute) Type() Type {
	return ae.Right.Type()
}

type PlusAssignmentExecute struct {
	Assignment *Assignment

	Left, Right Executable
}

func (pae PlusAssignmentExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	// left := pae.Left.Execute(ee)
	right := pae.Right.Execute(ee)
	return right
}

func (pae PlusAssignmentExecute) Type() Type {
	return pae.Right.Type()
}

func (l List) Compile() (Executable, CompileErrors) {
	log.Printf("compliling list")

	var errs CompileErrors
	var list ListExecute

	for _, item := range l.Items {
		list.Items = append(list.Items, errs.Collect(item.Compile()))
	}

	return list, errs
}

type ListExecute struct {
	List *List

	Items []Executable
}

type ListResult struct {
	Items []any
}

func (le ListExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	result := ListResult{}

	for _, Item := range le.Items {
		result.Items = append(result.Items, Item.Execute(ee))
	}

	return result
}

func (le ListExecute) Type() Type {
	return TypeList
}

func (e Expression) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	if e.Assignment == nil {
		return nil, errs.Append(fmt.Errorf("expression is empty at %s", e.Pos))
	}

	return e.Assignment.Compile()
}

func (p Program) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	block := Block{}

	for _, c := range p.Commands {
		if c != nil && c.Expression != nil {
			e := errs.Collect(c.Expression.Compile())
			block.Commands = append(block.Commands, e)
		}
	}
	return block, errs
}

func (a *Addition) Compile() (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(a.Multiplication.Compile())
	if len(a.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range a.Operations {
		operand := errs.Collect(opExpr.Operand.Compile())
		if ex.Type() != operand.Type() {
			errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(), operand.Type(), a.Pos))
		}
		operands = append(operands, operand)
	}

	switch ex.Type() {
	case TypeFloat:
		for i, operand := range operands {
			switch a.Operations[i].Op {
			case "+":
				ex = AddFloat{a, ex, operand}
			case "-":
				ex = SubtractFloat{a, ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", a.Operations[i].Op, ex.Type(), a.Pos))
				return nil, errs
			}
		}

	case TypeInteger:
		for i, operand := range operands {
			switch a.Operations[i].Op {
			case "+":
				ex = AddInteger{a, ex, operand}
			case "-":
				ex = SubtractInteger{a, ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", a.Operations[i].Op, ex.Type(), a.Pos))
				return nil, errs
			}
		}

	case TypeString:
		for i, operand := range operands {
			switch a.Operations[i].Op {
			case "+":
				ex = AddString{a, ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid type %s for - at %s", ex.Type(), a.Pos))
				return nil, errs
			}
		}

	default:
		errs.Append(fmt.Errorf("invalid type %s for +, - at %s", ex.Type(), a.Pos))
		return nil, errs
	}

	return ex, errs
}

type AddFloat struct {
	Addition    *Addition
	Left, Right Executable
}

func (af AddFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	return FloatValue(af.Left.Execute(ee).(FloatValue) + af.Right.Execute(ee).(FloatValue))
}

func (af AddFloat) Type() Type {
	return TypeFloat
}

type SubtractFloat struct {
	Addition    *Addition
	Left, Right Executable
}

func (mf SubtractFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	return FloatValue(mf.Left.Execute(ee).(FloatValue) - mf.Right.Execute(ee).(FloatValue))
}

func (mf SubtractFloat) Type() Type {
	return TypeFloat
}

type AddInteger struct {
	Addition    *Addition
	Left, Right Executable
}

func (ai AddInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(ai.Left.Execute(ee).(IntegerValue) + ai.Right.Execute(ee).(IntegerValue))
}

func (ai AddInteger) Type() Type {
	return TypeInteger
}

type SubtractInteger struct {
	Addition    *Addition
	Left, Right Executable
}

func (si SubtractInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(si.Left.Execute(ee).(IntegerValue) - si.Right.Execute(ee).(IntegerValue))
}

func (si SubtractInteger) Type() Type {
	return TypeInteger
}

type AddString struct {
	Addition    *Addition
	Left, Right Executable
}

func (as AddString) Execute(ee ExecutionEnvironment) ExecutionResult {
	return StringValue(as.Left.Execute(ee).(StringValue) + as.Right.Execute(ee).(StringValue))
}

func (as AddString) Type() Type {
	return TypeString
}
