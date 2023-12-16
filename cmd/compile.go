package main

import (
	"fmt"
	"log"
)

type ExecutionEnvironment map[string]any
type ExecutionResult any

type Executable interface {
	Execute(ExecutionEnvironment) ExecutionResult
	Type(TypeMap) Type
}

type TypeMap map[string]Type

type Type uint

const (
	TypeUnknown Type = iota
	TypeIdentifier
	TypeBool
	TypeFloat
	TypeInteger
	TypeString
	TypeList
	TypeKeyValue
	TypeNamedFunction
	TypeCount
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
	case TypeIdentifier:
		return "identifier"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

func TypeOf(x any) Type {
	switch x.(type) {
	case BoolValue:
		return TypeBool
	case FloatValue:
		return TypeFloat
	case IntegerValue:
		return TypeInteger
	case StringValue:
		return TypeString
	case ListValue:
		return TypeList
	case KeyValueResult:
		return TypeKeyValue
	case IdentifierValue:
		return TypeIdentifier
	default:
		return TypeUnknown
	}
}

var (
	PlusOpMap = [TypeCount]func(any, any) any{
		TypeFloat:   func(a, b any) any { return FloatValue(a.(FloatValue) + b.(FloatValue)) },
		TypeInteger: func(a, b any) any { return IntegerValue(a.(IntegerValue) + b.(IntegerValue)) },
		TypeString:  func(a, b any) any { return StringValue(a.(StringValue) + b.(StringValue)) },
	}
	MinusOpMap = [TypeCount]func(any, any) any{
		TypeFloat:   func(a, b any) any { return FloatValue(a.(FloatValue) - b.(FloatValue)) },
		TypeInteger: func(a, b any) any { return IntegerValue(a.(IntegerValue) - b.(IntegerValue)) },
	}
	MultOpMap = [TypeCount]func(any, any) any{
		TypeFloat:   func(a, b any) any { return FloatValue(a.(FloatValue) * b.(FloatValue)) },
		TypeInteger: func(a, b any) any { return IntegerValue(a.(IntegerValue) * b.(IntegerValue)) },
	}
	DivOpMap = [TypeCount]func(any, any) any{
		TypeFloat:   func(a, b any) any { return FloatValue(a.(FloatValue) / b.(FloatValue)) },
		TypeInteger: func(a, b any) any { return IntegerValue(a.(IntegerValue) / b.(IntegerValue)) },
	}
	ModuloOpMap = [TypeCount]func(any, any) any{
		TypeInteger: func(a, b any) any { return IntegerValue(a.(IntegerValue) % b.(IntegerValue)) },
	}
	EqualOpMap = [TypeCount]func(any, any) BoolValue{
		TypeBool:    func(a, b any) BoolValue { return BoolValue(a.(BoolValue) == b.(BoolValue)) },
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) == b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) == b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) == b.(StringValue)) },
	}
	NotEqualOpMap = [TypeCount]func(any, any) BoolValue{
		TypeBool:    func(a, b any) BoolValue { return BoolValue(a.(BoolValue) != b.(BoolValue)) },
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) != b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) != b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) != b.(StringValue)) },
	}
	LessThanOpMap = [TypeCount]func(any, any) BoolValue{
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) < b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) < b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) < b.(StringValue)) },
	}
	LessThanOrEqualOpMap = [TypeCount]func(any, any) BoolValue{
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) <= b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) <= b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) <= b.(StringValue)) },
	}
	GreaterThanOpMap = [TypeCount]func(any, any) BoolValue{
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) > b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) > b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) > b.(StringValue)) },
	}
	GreaterThanOrEqualOpMap = [TypeCount]func(any, any) BoolValue{
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) >= b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) >= b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) >= b.(StringValue)) },
	}
)

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

type ExecutableBlock struct {
	Commands []Executable
}

func (b ExecutableBlock) Execute(ee ExecutionEnvironment) ExecutionResult {
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

func (b ExecutableBlock) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

type BoolValue bool

func (b BoolValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("bool: %t", b)
	return b
}

func (b BoolValue) Type(typeMap TypeMap) Type {
	return TypeBool
}

type FloatValue float64

func (f FloatValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("float: %.20f", f)
	return f
}

func (f FloatValue) Type(typeMap TypeMap) Type {
	return TypeFloat
}

type IntegerValue int64

func (i IntegerValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("integer: %d", i)
	return i
}

func (i IntegerValue) Type(typeMap TypeMap) Type {
	return TypeInteger
}

type StringValue string

func (s StringValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("string: %#v", s)
	return s
}

func (s StringValue) Type(typeMap TypeMap) Type {
	return TypeString
}

type IdentifierValue struct {
	Base  *Base
	Value string
}

func (i IdentifierValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("identifier: %s", i)
	v, ok := ee[i.Value]
	if !ok {
		log.Fatalf("accessing variable %s before assignment at %s", i.Value, i.Base.Pos)
	}
	return v
}

func (i IdentifierValue) Type(typeMap TypeMap) Type {
	return typeMap[i.Value]
}

type ListValue []Executable

func (l ListValue) Execute(ee ExecutionEnvironment) ExecutionResult {
	// log.Printf("list: %v", l)
	return l
}

func (l ListValue) Type(typeMap TypeMap) Type {
	return TypeList
}

func (b *Base) Compile(typeMap TypeMap) (Executable, CompileErrors) {
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
		return IdentifierValue{b, *b.Ident}, NoErrors
	case b.StringValue != nil:
		return StringValue(*b.StringValue), NoErrors
	case b.Subexpression != nil:
		return b.Subexpression.Compile(typeMap)
	case b.List != nil:
		return b.List.Compile(typeMap)
	case b.StatementBlock != nil:
		return b.StatementBlock.Compile(typeMap)
	case b.NamedFunction != nil:
		return b.NamedFunction.Compile(typeMap)
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

func (uen UnaryExecuteNot) Type(typeMap TypeMap) Type {
	return TypeBool
}

type UnaryExecuteSubtractFloat struct {
	Unary   *Unary
	Operand Executable
}

func (uemf UnaryExecuteSubtractFloat) Execute(ee ExecutionEnvironment) ExecutionResult {
	return FloatValue(-uemf.Operand.Execute(ee).(FloatValue))
}

func (uemf UnaryExecuteSubtractFloat) Type(typeMap TypeMap) Type {
	return TypeFloat
}

type UnaryExecuteMinusInteger struct {
	Unary   *Unary
	Operand Executable
}

func (uemi UnaryExecuteMinusInteger) Execute(ee ExecutionEnvironment) ExecutionResult {
	return IntegerValue(-uemi.Operand.Execute(ee).(IntegerValue))
}

func (uemi UnaryExecuteMinusInteger) Type(typeMap TypeMap) Type {
	return TypeInteger
}

func (u *Unary) Compile(typeMap TypeMap) (Executable, CompileErrors) {

	if u.Base != nil {
		return u.Base.Compile(typeMap)
	}

	if u.Unary == nil {
		return nil, NewError(fmt.Errorf("cannot compile unary %s", u))
	}

	operand, errs := u.Unary.Compile(typeMap)

	switch *u.Op {
	case "!":
		if operand.Type(typeMap) != TypeBool {
			errs.Append(fmt.Errorf("invalid unary operation %s at %s", *u.Op, u.Pos))
		}
		return UnaryExecuteNot{u, operand}, errs
	case "-":
		switch operand.Type(typeMap) {
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

func (nf *NamedFunction) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	if nf.Name == nil {
		return nil, errs.Append(fmt.Errorf("invalid named function at %s", nf.Pos))
	}

	ex := errs.Collect(nf.Body.Compile(typeMap))

	return NamedFunctionExecute{nf, ex}, errs
}

type NamedFunctionExecute struct {
	NamedFunction *NamedFunction
	Block         Executable
}

func (nfe NamedFunctionExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	// todo: load into global function map
	return nil
}

func (nfe NamedFunctionExecute) Type(typeMap TypeMap) Type {
	return TypeNamedFunction
}

func (sb *StatementBlock) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	block := ExecutableBlock{}

	for _, e := range sb.Statements {
		if e != nil {
			ex := errs.Collect(e.Compile(typeMap))
			block.Commands = append(block.Commands, ex)
		}
	}

	return block, errs
}

func (m *Multiplication) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(m.Unary.Compile(typeMap))
	if len(m.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range m.Operations {
		operand := errs.Collect(opExpr.Operand.Compile(typeMap))
		if ex.Type(typeMap) != operand.Type(typeMap) {
			errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(typeMap), operand.Type(typeMap), m.Pos))
		}
		operands = append(operands, operand)
	}

	multOp := MultOpMap[ex.Type(typeMap)]
	divOp := DivOpMap[ex.Type(typeMap)]
	moduloOp := ModuloOpMap[ex.Type(typeMap)]

	for i, operand := range operands {
		switch m.Operations[i].Op {
		case "*":
			if multOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{multOp, ex, operand, ex.Type(typeMap)}
		case "/":
			if divOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{divOp, ex, operand, ex.Type(typeMap)}
		case "%":
			if moduloOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{moduloOp, ex, operand, ex.Type(typeMap)}
		default:
			errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
		}
	}

	return ex, errs
}

func (kv *KeyValue) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(kv.Addition.Compile(typeMap))
	if kv.RightValue == nil {
		return ex, errs
	}

	if ex.Type(typeMap) != TypeString {
		errs.Append(fmt.Errorf("invalid type %s for key (should be string) at %s", ex.Type(typeMap), kv.Pos))
	}

	rex := errs.Collect(kv.RightValue.Compile(typeMap))
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

func (kve KeyValueExecute) Type(typeMap TypeMap) Type {
	return TypeKeyValue
}

func (l *Logical) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(l.Comparison.Compile(typeMap))
	if len(l.Operations) == 0 {
		return ex, errs
	}

	if ex.Type(typeMap) != TypeBool {
		errs.Append(fmt.Errorf("invalid type %s for logical operation at %s", ex.Type(typeMap), l.Pos))
	}

	operands := []Executable{}
	for _, opExpr := range l.Operations {
		operands = append(operands, errs.Collect(opExpr.Operand.Compile(typeMap)))
	}

	for i, operand := range operands {
		switch ex.Type(typeMap) {
		case TypeBool:
			switch l.Operations[i].Op {
			case "&&":
				ex = ShortCircuitAnd{ex, operand}
			case "||":
				ex = ShortCircuitOr{ex, operand}
			default:
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", l.Operations[i].Op, ex.Type(typeMap), l.Pos))
				return nil, errs
			}
		default:
			errs.Append(fmt.Errorf("invalid type %s for logical operation at %s", ex.Type(typeMap), l.Pos))
			return nil, errs
		}
	}

	return ex, errs
}

func (c *Comparison) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(c.KeyValue.Compile(typeMap))
	if len(c.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range c.Operations {
		operands = append(operands, errs.Collect(opExpr.Operand.Compile(typeMap)))
	}

	for i, operand := range operands {
		switch c.Operations[i].Op {
		case "==":
			equalOp := EqualOpMap[ex.Type(typeMap)]
			if equalOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
				continue
			}
			ex = ComparisonOperation{equalOp, ex, operand}
		case "!=":
			notEqualOp := NotEqualOpMap[ex.Type(typeMap)]
			if notEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
				continue
			}
			ex = ComparisonOperation{notEqualOp, ex, operand}
		case "<":
			lessThanOp := LessThanOpMap[ex.Type(typeMap)]
			if lessThanOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
				continue
			}
			ex = ComparisonOperation{lessThanOp, ex, operand}
		case "<=":
			lessThanOrEqualOp := LessThanOrEqualOpMap[ex.Type(typeMap)]
			if lessThanOrEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{lessThanOrEqualOp, ex, operand}
		case ">":
			greaterThanOp := GreaterThanOpMap[ex.Type(typeMap)]
			if greaterThanOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{greaterThanOp, ex, operand}
		case ">=":
			greaterThanOrEqualOp := GreaterThanOrEqualOpMap[ex.Type(typeMap)]
			if greaterThanOrEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{greaterThanOrEqualOp, ex, operand}
		default:
			errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
		}
	}

	return ex, errs
}

func (p *Pipe) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(p.Logical.Compile(typeMap))
	if len(p.Operations) == 0 {
		return ex, errs
	}

	pe := PipeExecute{p, []Executable{ex}}
	for _, op := range p.Operations {
		pe.Commands = append(pe.Commands, errs.Collect(op.Operand.Compile(typeMap)))
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

func (pe PipeExecute) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

func (a *Assignment) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(a.Pipe.Compile(typeMap))
	if a.Operation == nil {
		return ex, errs
	}

	operand := errs.Collect(a.Operation.Operand.Compile(typeMap))
	if operand == nil {
		return nil, errs
	}

	if !IsIdentifier(ex) {
		errs.Append(fmt.Errorf("invalid left hand side %s for assignment at %s", ex.Type(typeMap), a.Pos))
		return nil, errs
	}
	curType, ok := typeMap[ex.(IdentifierValue).Value]
	if !ok {
		typeMap[ex.(IdentifierValue).Value] = operand.Type(typeMap)
	} else if curType != operand.Type(typeMap) {
		errs.Append(fmt.Errorf("cannot change variable %s from type %s to type %s at %s",
			ex.(IdentifierValue).Value, curType, operand.Type(typeMap), a.Pos))
	}
	switch a.Operation.Op {
	case ":=":
		ex = AssignmentExecute{a, ex.(IdentifierValue), operand}
	case "+=":
		ex = PlusAssignmentExecute{a, ex.(IdentifierValue), operand}
	default:
		errs.Append(fmt.Errorf("invalid assignment operator %s at %s", a.Operation.Op, a.Pos))
	}

	return ex, errs
}

func IsIdentifier(e Executable) bool {
	_, ok := e.(IdentifierValue)
	return ok
}

type AssignmentExecute struct {
	Assignment *Assignment
	Left       IdentifierValue
	Right      Executable
}

func (ae AssignmentExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	right := ae.Right.Execute(ee)

	leftType := TypeOf(ee[ae.Left.Value])
	if leftType != TypeUnknown && leftType != TypeOf(right) {
		log.Fatalf("cannot change variable %s from type %s to type %s at %s", ae.Left.Value, leftType, TypeOf(right), ae.Assignment.Pos)
	}

	ee[ae.Left.Value] = right
	return right
}

func (ae AssignmentExecute) Type(typeMap TypeMap) Type {
	return ae.Right.Type(typeMap)
}

type PlusAssignmentExecute struct {
	Assignment *Assignment
	Left       IdentifierValue
	Right      Executable
}

func (pae PlusAssignmentExecute) Execute(ee ExecutionEnvironment) ExecutionResult {
	left := ee[pae.Left.Value]
	right := pae.Right.Execute(ee)

	if TypeOf(left) != TypeOf(right) {
		log.Fatalf("type mismatch %s/%s for += at %s", TypeOf(left), TypeOf(right), pae.Assignment.Pos)
	}

	plusOp := PlusOpMap[TypeOf(left)]
	if plusOp == nil {
		log.Fatalf("invalid type %s for += at %s", TypeOf(left), pae.Assignment.Pos)
	}

	newVal := plusOp(left, right)
	ee[pae.Left.Value] = newVal
	return newVal
}

func (pae PlusAssignmentExecute) Type(typeMap TypeMap) Type {
	return pae.Right.Type(typeMap)
}

func (l List) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	log.Printf("compliling list")

	var errs CompileErrors
	var list ListExecute

	for _, item := range l.Items {
		list.Items = append(list.Items, errs.Collect(item.Compile(typeMap)))
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

func (le ListExecute) Type(typeMap TypeMap) Type {
	return TypeList
}

func (e Expression) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	if e.Assignment == nil {
		return nil, errs.Append(fmt.Errorf("expression is empty at %s", e.Pos))
	}

	return e.Assignment.Compile(typeMap)
}

func (p Program) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	block := ExecutableBlock{}

	for _, c := range p.Commands {
		if c != nil && c.Expression != nil {
			e := errs.Collect(c.Expression.Compile(typeMap))
			block.Commands = append(block.Commands, e)
		}
	}
	return block, errs
}

func (a *Addition) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(a.Multiplication.Compile(typeMap))
	if len(a.Operations) == 0 {
		return ex, errs
	}

	operands := []Executable{}
	for _, opExpr := range a.Operations {
		operand := errs.Collect(opExpr.Operand.Compile(typeMap))
		if ex.Type(typeMap) != operand.Type(typeMap) {
			errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(typeMap), operand.Type(typeMap), a.Pos))
		}
		operands = append(operands, operand)
	}

	plus := PlusOpMap[ex.Type(typeMap)]
	minus := MinusOpMap[ex.Type(typeMap)]
	if plus == nil || minus == nil {
		errs.Append(fmt.Errorf("invalid type %s for +, - at %s", ex.Type(typeMap), a.Pos))
		return nil, errs
	}

	for i, operand := range operands {
		switch a.Operations[i].Op {
		case "+":
			ex = BinaryOperation{plus, ex, operand, ex.Type(typeMap)}
		case "-":
			ex = BinaryOperation{minus, ex, operand, ex.Type(typeMap)}
		default:
			errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", a.Operations[i].Op, ex.Type(typeMap), a.Pos))
		}
	}

	return ex, errs
}

type BinaryOperation struct {
	Func        func(any, any) any
	Left, Right Executable
	TypeVal     Type
}

func (bo BinaryOperation) Execute(ee ExecutionEnvironment) ExecutionResult {
	return bo.Func(bo.Left.Execute(ee), bo.Right.Execute(ee))
}

func (bo BinaryOperation) Type(typeMap TypeMap) Type {
	return bo.TypeVal
}

type ComparisonOperation struct {
	Func        func(any, any) BoolValue
	Left, Right Executable
}

func (co ComparisonOperation) Execute(ee ExecutionEnvironment) ExecutionResult {
	return co.Func(co.Left.Execute(ee), co.Right.Execute(ee))
}

func (co ComparisonOperation) Type(typeMap TypeMap) Type {
	return TypeBool
}

type ShortCircuitAnd struct {
	Left, Right Executable
}

// Execute returns the result of the left expression if it is false, otherwise it returns the result of the right expression.
func (sca ShortCircuitAnd) Execute(ee ExecutionEnvironment) ExecutionResult {
	leftValue := sca.Left.Execute(ee)
	if !bool(leftValue.(BoolValue)) {
		return BoolValue(false)
	}
	return sca.Right.Execute(ee)
}

func (sca ShortCircuitAnd) Type(typeMap TypeMap) Type {
	return TypeBool
}

type ShortCircuitOr struct {
	Left, Right Executable
}

// Execute returns the result of the left expression if it is true, otherwise it returns the result of the right expression.
func (sco ShortCircuitOr) Execute(ee ExecutionEnvironment) ExecutionResult {
	leftValue := sco.Left.Execute(ee)
	if bool(leftValue.(BoolValue)) {
		return BoolValue(true)
	}
	return sco.Right.Execute(ee)
}

func (sco ShortCircuitOr) Type(typeMap TypeMap) Type {
	return TypeBool
}
