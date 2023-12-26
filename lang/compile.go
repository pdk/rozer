package lang

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alecthomas/participle/v2/lexer"
)

var (
	TagNull     = TagValue{Value: "#null"}
	TagComplete = TagValue{Value: "#complete"}
	TagContinue = TagValue{Value: "#continue"}
	TagBreak    = TagValue{Value: "#break"}
)

type Parameterized interface {
	Executable
	ParameterNames() []string
	Apply(*ExecutionEnvironment) ExecutionResult
}

type ExecutionEnvironment struct {
	global map[string]any
	local  map[string]any
}

func NewExecutionEnvironment() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		global: map[string]any{},
		local:  map[string]any{},
	}
}

func (ee *ExecutionEnvironment) NewLocalEnvironment() *ExecutionEnvironment {
	return &ExecutionEnvironment{
		global: ee.global,
		local:  map[string]any{},
	}
}

func (ee *ExecutionEnvironment) Get(key string) any {
	v, ok := ee.global[key]
	if ok {
		return v
	}

	return ee.local[key]

	// v, ok = ee.local[key]
	// if !ok {
	// 	log.Printf("accessing variable %s before assignment", key)
	// }

	// return v
}

func (ee *ExecutionEnvironment) SetGlobal(key string, value any) {
	ee.global[key] = value
}

func (ee *ExecutionEnvironment) GlobalExists(key string) bool {
	_, ok := ee.global[key]
	return ok
}

func (ee *ExecutionEnvironment) Set(key string, value any) {
	if ee.GlobalExists(key) {
		log.Fatalf("cannot reassign global variable %s", key)
	}

	ee.local[key] = value
}

type ExecutionResult any

type Executable interface {
	Execute(*ExecutionEnvironment) ExecutionResult
	Type(TypeMap) Type
	ListRep() []any
}

type TypeMap map[string]Type

type Type uint

const (
	TypeUnknown Type = iota
	TypeTag
	TypeIdentifier
	TypeBool
	TypeFloat
	TypeInteger
	TypeString
	TypeList
	TypeKeyValue
	TypeFunction
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
	case TypeTag:
		return "tag"
	case TypeIdentifier:
		return "identifier"
	case TypeFunction:
		return "function"
	default:
		return fmt.Sprintf("unknown type %d", t)
	}
}

func TypeOf(x any) Type {
	switch x.(type) {
	case nil:
		return TypeUnknown
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
	case TagValue:
		return TypeTag
	case IdentifierValue:
		return TypeIdentifier
	case FunctionExecute:
		return TypeFunction
	default:
		log.Printf("unknown type %T", x)
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
		TypeTag:     func(a, b any) BoolValue { return BoolValue(a.(TagValue).Value == b.(TagValue).Value) },
	}
	NotEqualOpMap = [TypeCount]func(any, any) BoolValue{
		TypeBool:    func(a, b any) BoolValue { return BoolValue(a.(BoolValue) != b.(BoolValue)) },
		TypeFloat:   func(a, b any) BoolValue { return BoolValue(a.(FloatValue) != b.(FloatValue)) },
		TypeInteger: func(a, b any) BoolValue { return BoolValue(a.(IntegerValue) != b.(IntegerValue)) },
		TypeString:  func(a, b any) BoolValue { return BoolValue(a.(StringValue) != b.(StringValue)) },
		TypeTag:     func(a, b any) BoolValue { return BoolValue(a.(TagValue).Value != b.(TagValue).Value) },
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

func init() {
	PlusOpMap[TypeUnknown] = unknownHandler[any](PlusOpMap, "add")
	MinusOpMap[TypeUnknown] = unknownHandler[any](MinusOpMap, "subtract")
	MultOpMap[TypeUnknown] = unknownHandler[any](MultOpMap, "multiply")
	DivOpMap[TypeUnknown] = unknownHandler[any](DivOpMap, "divide")
	ModuloOpMap[TypeUnknown] = unknownHandler[any](ModuloOpMap, "modulo")
	EqualOpMap[TypeUnknown] = unknownHandler[BoolValue](EqualOpMap, "compare")
	NotEqualOpMap[TypeUnknown] = unknownHandler[BoolValue](NotEqualOpMap, "compare")
	LessThanOpMap[TypeUnknown] = unknownHandler[BoolValue](LessThanOpMap, "compare")
	LessThanOrEqualOpMap[TypeUnknown] = unknownHandler[BoolValue](LessThanOrEqualOpMap, "compare")
	GreaterThanOpMap[TypeUnknown] = unknownHandler[BoolValue](GreaterThanOpMap, "compare")
	GreaterThanOrEqualOpMap[TypeUnknown] = unknownHandler[BoolValue](GreaterThanOrEqualOpMap, "compare")
}

func unknownHandler[T any](m [TypeCount]func(any, any) T, desc string) func(any, any) T {
	return func(a, b any) T {
		aType := TypeOf(a)
		bType := TypeOf(b)
		if aType != bType || aType == TypeUnknown || bType == TypeUnknown {
			log.Fatalf("cannot %s types %s, %s at %s", desc, aType, bType, currentPos)
		}
		op := m[aType]
		if op == nil {
			log.Fatalf("cannot %s type %s at %s", desc, aType, currentPos)
		}
		return op(a, b)
	}
}

type CompileErrors struct {
	Errs *[]error
}

var (
	NoErrors = CompileErrors{}
)

func (ce *CompileErrors) Append(err ...error) CompileErrors {
	if ce.Errs == nil {
		ce.Errs = &[]error{}
	}

	*ce.Errs = append(*ce.Errs, err...)

	return *ce
}

func (ce *CompileErrors) Collect(e Executable, errs CompileErrors) Executable {
	if errs.Errs != nil {
		ce.Append(*errs.Errs...)
	}
	return e
}

func NewError(err error) CompileErrors {
	return CompileErrors{&[]error{err}}
}

func (ce CompileErrors) Len() int {
	if ce.Errs == nil {
		return 0
	}
	return len(*ce.Errs)
}

type ExecutableBlock struct {
	Commands []Executable
}

func (b ExecutableBlock) Execute(ee *ExecutionEnvironment) ExecutionResult {
	var lastResult ExecutionResult
	for _, c := range b.Commands {
		lastResult = c.Execute(ee)
		log.Printf("execution returned %v", lastResult)
	}
	return lastResult
}

func (b ExecutableBlock) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

func (b ExecutableBlock) ListRep() []any {
	commands := []any{}
	for _, c := range b.Commands {
		commands = append(commands, c.ListRep())
	}
	return []any{"block", commands}
}

type BoolValue bool

func (b BoolValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return b
}

func (b BoolValue) Type(typeMap TypeMap) Type {
	return TypeBool
}

func (b BoolValue) ListRep() []any {
	return []any{"bool", fmt.Sprintf("%t", b)}
}

type FloatValue float64

func (f FloatValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return f
}

func (f FloatValue) Type(typeMap TypeMap) Type {
	return TypeFloat
}

func (f FloatValue) ListRep() []any {
	return []any{"float", fmt.Sprintf("%f", f)}
}

type IntegerValue int64

func (i IntegerValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return i
}

func (i IntegerValue) Type(typeMap TypeMap) Type {
	return TypeInteger
}

func (i IntegerValue) ListRep() []any {
	return []any{"integer", fmt.Sprintf("%d", i)}
}

type StringValue string

func (s StringValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return s
}

func (s StringValue) Type(typeMap TypeMap) Type {
	return TypeString
}

func (s StringValue) ListRep() []any {
	return []any{"string", s}
}

type TagValue struct {
	Base  *Base
	Value string
}

func (t TagValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return t
}

func (t TagValue) Type(typeMap TypeMap) Type {
	return TypeTag
}

func (t TagValue) ListRep() []any {
	return []any{"tag", t.Value}
}

type IdentifierValue struct {
	Base  *Base
	Value string
}

func (i IdentifierValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return ee.Get(i.Value)
}

func (i IdentifierValue) Type(typeMap TypeMap) Type {
	return typeMap[i.Value]
}

func (i IdentifierValue) ListRep() []any {
	return []any{"ident", i.Value}
}

type ListValue []Executable

func (l ListValue) Execute(ee *ExecutionEnvironment) ExecutionResult {
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
	case b.Tag != nil:
		return TagValue{b, *b.Tag}, NoErrors
	case b.Ident != nil:
		return IdentifierValue{b, *b.Ident}, NoErrors
	case b.StringValue != nil:
		return StringValue(*b.StringValue), NoErrors
	case b.Subexpression != nil:
		return b.Subexpression.Compile(typeMap)
	case b.List != nil:
		return b.List.Compile(typeMap)
	// case b.StatementBlock != nil:
	// 	return b.StatementBlock.Compile(typeMap)
	case b.UnnamedFunction != nil:
		return b.UnnamedFunction.Compile(typeMap)
	case b.Invocation != nil:
		return b.Invocation.Compile(typeMap)
	}

	return nil, NewError(fmt.Errorf("cannot compile base %#v", b))
}

type UnaryExecuteNot struct {
	Unary   *Unary
	Operand Executable
}

func (uen UnaryExecuteNot) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return BoolValue(!uen.Operand.Execute(ee).(BoolValue))
}

func (uen UnaryExecuteNot) Type(typeMap TypeMap) Type {
	return TypeBool
}

func (uen UnaryExecuteNot) ListRep() []any {
	return []any{"!", uen.Operand.ListRep()}
}

type UnaryExecuteSubtractFloat struct {
	Unary   *Unary
	Operand Executable
}

func (uemf UnaryExecuteSubtractFloat) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return FloatValue(-uemf.Operand.Execute(ee).(FloatValue))
}

func (uemf UnaryExecuteSubtractFloat) Type(typeMap TypeMap) Type {
	return TypeFloat
}

func (uemf UnaryExecuteSubtractFloat) ListRep() []any {
	return []any{"-", uemf.Operand.ListRep()}
}

type UnaryExecuteMinusInteger struct {
	Unary   *Unary
	Operand Executable
}

func (uemi UnaryExecuteMinusInteger) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return IntegerValue(-uemi.Operand.Execute(ee).(IntegerValue))
}

func (uemi UnaryExecuteMinusInteger) Type(typeMap TypeMap) Type {
	return TypeInteger
}

func (uemi UnaryExecuteMinusInteger) ListRep() []any {
	return []any{"-", uemi.Operand.ListRep()}
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

type InvocationExecute struct {
	*Invocation
	ProducerExecutable  Executable
	ExecutableArguments []Executable
}

func (i InvocationExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {

	// look up the function
	// log.Printf("begin function invocation with %#v", i.ProducerExecutable)
	producerResult := i.ProducerExecutable.Execute(ee)
	// log.Printf("the function is a %T", producerResult)

	switch producerResult.(type) {
	case Parameterized:
		// looks good. fall thru to execute the function
	default:
		log.Fatalf("invalid function invocation (expecting Parameterized, got %T) at %s: %s", producerResult, i.Pos, i.String())
	}

	functionExecute := producerResult.(Parameterized)
	params := functionExecute.ParameterNames()

	if len(params) != len(i.Arguments) {
		log.Fatalf("function invocation expecting %d params, but got %d at %s: %s", len(params), len(i.Arguments), i.Pos, i.String())
	}

	// compute the values of the arguments
	log.Printf("computing the invocation arguments values")
	values := make([]ExecutionResult, len(i.Arguments))
	for i, e := range i.ExecutableArguments {
		values[i] = e.Execute(ee)
		log.Printf("argument %d: %v", i, values[i])
	}

	// assign the values to the parameters in a new local environment
	local := ee.NewLocalEnvironment()
	for i, param := range params {
		local.Set(param, values[i])
	}

	// execute the function with the new environment
	result := functionExecute.Apply(local)
	log.Printf("invocation result: %v", result)

	return result
}

func (i InvocationExecute) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

func (i InvocationExecute) ListRep() []any {

	args := []any{}
	for _, arg := range i.ExecutableArguments {
		args = append(args, arg.ListRep())
	}

	return []any{"invocation", *i.Name, args}
}

func (i *Invocation) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	funcProducer := IdentifierValue{Value: *i.Name}

	execArgs := []Executable{}
	for _, arg := range i.Arguments {
		nextArg := errs.Collect(arg.Compile(typeMap))
		execArgs = append(execArgs, nextArg)
	}

	return InvocationExecute{i, funcProducer, execArgs}, errs
}

type FunctionExecute struct {
	*NamedFunction
	*UnnamedFunction
	ExecutableBlock
}

func (fe FunctionExecute) Apply(ee *ExecutionEnvironment) ExecutionResult {
	return fe.ExecutableBlock.Execute(ee)
}

func (fe FunctionExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return fe
}

func (fe FunctionExecute) Type(typeMap TypeMap) Type {
	return TypeFunction
}

func (fe FunctionExecute) ListRep() []any {
	params := []string{"params"}
	params = append(params, fe.ParameterNames()...)

	switch {
	case fe.NamedFunction != nil:
		return []any{"named function", fe.NamedFunction.Name, params, fe.ExecutableBlock.ListRep()}
	case fe.UnnamedFunction != nil:
		return []any{"unnamed function", params, fe.ExecutableBlock.ListRep()}
	default:
		return []any{"invalid function"}
	}
}

func (fe FunctionExecute) ParameterNames() []string {
	switch {
	case fe.NamedFunction != nil:
		return fe.NamedFunction.Params
	case fe.UnnamedFunction != nil:
		return fe.UnnamedFunction.Params
	default:
		return []string{}
	}
}

func (nf *NamedFunction) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	if nf.Name == nil {
		return nil, errs.Append(fmt.Errorf("invalid named function at %s", nf.Pos))
	}

	ex := errs.Collect(nf.Body.Compile(typeMap))

	return FunctionExecute{
		NamedFunction:   nf,
		ExecutableBlock: ex.(ExecutableBlock),
	}, errs
}

func (uf *UnnamedFunction) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(uf.Body.Compile(typeMap))

	return FunctionExecute{
		UnnamedFunction: uf,
		ExecutableBlock: ex.(ExecutableBlock),
	}, errs
}

func (rb *RequiredBlock) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	block := ExecutableBlock{}

	for _, e := range rb.Statements {
		if e != nil {
			ex := errs.Collect(e.Compile(typeMap))
			block.Commands = append(block.Commands, ex)
		}
	}

	return block, errs
}

func hasTypeUnknown(items []Executable, typeMap TypeMap) bool {
	for _, item := range items {
		if item.Type(typeMap) == TypeUnknown {
			return true
		}
	}
	return false
}

func (m *Multiplication) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(m.Unary.Compile(typeMap))
	if len(m.Operations) == 0 {
		return ex, errs
	}

	multOp := MultOpMap[ex.Type(typeMap)]
	divOp := DivOpMap[ex.Type(typeMap)]
	moduloOp := ModuloOpMap[ex.Type(typeMap)]

	operands := []Executable{}
	for _, opExpr := range m.Operations {
		operand := errs.Collect(opExpr.Operand.Compile(typeMap))
		operands = append(operands, operand)
	}

	atLeastOneUnknown := ex.Type(typeMap) == TypeUnknown || hasTypeUnknown(operands, typeMap)

	if atLeastOneUnknown {
		multOp = MultOpMap[TypeUnknown]
		divOp = DivOpMap[TypeUnknown]
		moduloOp = ModuloOpMap[TypeUnknown]
	} else {
		for _, operand := range operands {
			if ex.Type(typeMap) != operand.Type(typeMap) {
				errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(typeMap), operand.Type(typeMap), m.Pos))
			}
		}
	}

	for i, operand := range operands {
		switch m.Operations[i].Op {
		case "*":
			if multOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{m.Pos, "*", multOp, ex, operand, ex.Type(typeMap)}
		case "/":
			if divOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{m.Pos, "/", divOp, ex, operand, ex.Type(typeMap)}
		case "%":
			if moduloOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
				continue
			}
			ex = BinaryOperation{m.Pos, "%", moduloOp, ex, operand, ex.Type(typeMap)}
		default:
			errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", m.Operations[i].Op, ex.Type(typeMap), m.Pos))
		}
	}

	return ex, errs
}

func (s *Series) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(s.FromValue.Compile(typeMap))
	if s.ToValue == nil {
		return ex, errs
	}

	rex := errs.Collect(s.ToValue.Compile(typeMap))

	return SeriesExecute{s, ex, rex}, errs
}

type SeriesExecute struct {
	Series *Series

	From Executable
	To   Executable
}

type InnerFunctionExecute struct {
	Pos lexer.Position

	Function func() ExecutionResult
}

func (ife InnerFunctionExecute) Apply(ee *ExecutionEnvironment) ExecutionResult {
	return ife.Function()
}

func (ife InnerFunctionExecute) ParameterNames() []string {
	return []string{}
}

func (ife InnerFunctionExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return ife.Function()
}

func (ife InnerFunctionExecute) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

func (ife InnerFunctionExecute) ListRep() []any {
	return []any{"inner function"}
}

func (se SeriesExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	from := se.From.Execute(ee)
	to := se.To.Execute(ee)

	switch {
	case TypeOf(from) == TypeInteger && TypeOf(to) == TypeInteger:
		// log.Printf("constructing integer series")
		fromInt := from.(IntegerValue)
		toInt := to.(IntegerValue)
		incr := 1
		if fromInt > toInt {
			incr = -1
		}
		complete := false
		return InnerFunctionExecute{
			Pos: se.Series.Pos,
			Function: func() ExecutionResult {
				if complete {
					return TagComplete
				}
				retVal := fromInt
				complete = retVal == toInt
				fromInt += IntegerValue(incr)
				return retVal
			},
		}
	}

	return from
}

func (se SeriesExecute) Type(typeMap TypeMap) Type {
	return TypeList
}

func (se SeriesExecute) ListRep() []any {
	return []any{"series", se.From.ListRep(), se.To.ListRep()}
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

func (kve KeyValueExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	return KeyValueResult{kve.Key.Execute(ee), kve.Value.Execute(ee)}
}

func (kve KeyValueExecute) Type(typeMap TypeMap) Type {
	return TypeKeyValue
}

func (kve KeyValueExecute) ListRep() []any {
	return []any{"keyvalue", kve.Key.ListRep(), kve.Value.ListRep()}
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

	ex := errs.Collect(c.Series.Compile(typeMap))
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
			ex = ComparisonOperation{c.Pos, "==", equalOp, ex, operand}
		case "!=":
			notEqualOp := NotEqualOpMap[ex.Type(typeMap)]
			if notEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
				continue
			}
			ex = ComparisonOperation{c.Pos, "!=", notEqualOp, ex, operand}
		case "<":
			lessThanOp := LessThanOpMap[ex.Type(typeMap)]
			if lessThanOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
				continue
			}
			ex = ComparisonOperation{c.Pos, "<", lessThanOp, ex, operand}
		case "<=":
			lessThanOrEqualOp := LessThanOrEqualOpMap[ex.Type(typeMap)]
			if lessThanOrEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{c.Pos, "<=", lessThanOrEqualOp, ex, operand}
		case ">":
			greaterThanOp := GreaterThanOpMap[ex.Type(typeMap)]
			if greaterThanOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{c.Pos, ">", greaterThanOp, ex, operand}
		case ">=":
			greaterThanOrEqualOp := GreaterThanOrEqualOpMap[ex.Type(typeMap)]
			if greaterThanOrEqualOp == nil {
				errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", c.Operations[i].Op, ex.Type(typeMap), c.Pos))
			}
			ex = ComparisonOperation{c.Pos, ">=", greaterThanOrEqualOp, ex, operand}
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

	pe := PipeExecute{
		Pipe:       p,
		DoComplete: []bool{false},
		Commands:   []Executable{ex},
	}
	for _, op := range p.Operations {
		pe.Commands = append(pe.Commands, errs.Collect(op.Operand.Compile(typeMap)))
		pe.DoComplete = append(pe.DoComplete, op.Op == ">>>")
	}

	return pe, errs
}

type PipeExecute struct {
	Pipe *Pipe

	DoComplete []bool
	Commands   []Executable
}

func (pe PipeExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {

	results := make([]ExecutionResult, len(pe.Commands))
	for i, c := range pe.Commands {
		results[i] = c.Execute(ee)
	}

	functions := make([]Parameterized, len(results))
	for i, result := range results {
		switch fn := result.(type) {
		case Parameterized:
			if i > 0 && len(fn.ParameterNames()) != 1 {
				log.Fatalf("invalid pipeline (every target must accept 1 argument) at %s: %s", pe.Pipe.Pos, pe.Pipe.String())
			}
			functions[i] = fn
		default:
			log.Fatalf("invalid pipeline (expecting function, got %T) at %s: %s", fn, pe.Pipe.Pos, pe.Pipe.String())
		}
	}

	var lastResult ExecutionResult
	for {
		log.Printf("executing pipeline, lastResult=%v", lastResult)
		for i, fn := range functions {
			fnEnv := ee.NewLocalEnvironment()
			if i > 0 {
				firstParam := fn.ParameterNames()[0]
				fnEnv.Set(firstParam, lastResult)
			}
			lastResult = fn.Apply(fnEnv)
			if TagContinue == lastResult ||
				TagComplete == lastResult ||
				TagBreak == lastResult ||
				TagNull == lastResult {
				break
			}
		}
		if TagComplete == lastResult ||
			TagBreak == lastResult ||
			TagNull == lastResult {
			break
		}
	}

	return lastResult
}

func (pe PipeExecute) Type(typeMap TypeMap) Type {
	return TypeUnknown
}

func (pe PipeExecute) ListRep() []any {
	commands := []any{}
	for _, c := range pe.Commands {
		commands = append(commands, c.ListRep())
	}
	return []any{">>", commands}
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
	opType := operand.Type(typeMap)
	if !ok {
		typeMap[ex.(IdentifierValue).Value] = opType
	} else if curType != TypeUnknown && opType != TypeUnknown && curType != opType {
		errs.Append(fmt.Errorf("cannot change type of variable %s from %s to %s at %s",
			ex.(IdentifierValue).Value, curType, opType, a.Pos))
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

func (ae AssignmentExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	right := ae.Right.Execute(ee)

	leftType := TypeOf(ee.Get(ae.Left.Value))
	if leftType != TypeUnknown && leftType != TypeOf(right) {
		log.Fatalf("cannot change type of variable %s from %s to %s at %s", ae.Left.Value, leftType, TypeOf(right), ae.Assignment.Pos)
	}

	ee.Set(ae.Left.Value, right)
	return right
}

func (ae AssignmentExecute) Type(typeMap TypeMap) Type {
	return ae.Right.Type(typeMap)
}

func (ae AssignmentExecute) ListRep() []any {
	return []any{ae.Assignment.Operation.Op, ae.Left.ListRep(), ae.Right.ListRep()}
}

type PlusAssignmentExecute struct {
	Assignment *Assignment
	Left       IdentifierValue
	Right      Executable
}

func (pae PlusAssignmentExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	left := ee.Get(pae.Left.Value)
	right := pae.Right.Execute(ee)

	if TypeOf(left) != TypeOf(right) {
		log.Fatalf("type mismatch %s/%s for += at %s", TypeOf(left), TypeOf(right), pae.Assignment.Pos)
	}

	plusOp := PlusOpMap[TypeOf(left)]
	if plusOp == nil {
		log.Fatalf("invalid type %s for += at %s", TypeOf(left), pae.Assignment.Pos)
	}

	newVal := plusOp(left, right)
	ee.Set(pae.Left.Value, newVal)
	return newVal
}

func (pae PlusAssignmentExecute) Type(typeMap TypeMap) Type {
	return pae.Right.Type(typeMap)
}

func (pae PlusAssignmentExecute) ListRep() []any {
	return []any{pae.Assignment.Operation.Op, pae.Left.ListRep(), pae.Right.ListRep()}
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

func (le ListExecute) Execute(ee *ExecutionEnvironment) ExecutionResult {
	result := ListResult{}

	for _, Item := range le.Items {
		result.Items = append(result.Items, Item.Execute(ee))
	}

	return result
}

func (le ListExecute) Type(typeMap TypeMap) Type {
	return TypeList
}

func (le ListExecute) ListRep() []any {
	items := []any{}
	for _, item := range le.Items {
		items = append(items, item.ListRep())
	}
	return []any{"list", items}
}

func (e Expression) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	if e.Assignment == nil {
		return nil, errs.Append(fmt.Errorf("expression is empty at %s", e.Pos))
	}

	return e.Assignment.Compile(typeMap)
}

func (p Program) Compile(typeMap TypeMap) (ProgramExecute, CompileErrors) {
	var errs CompileErrors

	block := ExecutableBlock{}
	functions := []FunctionExecute{}

	for _, c := range p.Commands {
		if c != nil {
			switch {
			case c.NamedFunction != nil:
				e := errs.Collect(c.NamedFunction.Compile(typeMap))
				functions = append(functions, e.(FunctionExecute))
			case c.Expression != nil:
				e := errs.Collect(c.Expression.Compile(typeMap))
				block.Commands = append(block.Commands, e)
			}
		}
	}

	return ProgramExecute{
		Program:         p,
		NamedFunctions:  functions,
		ExecutableBlock: block,
	}, errs
}

type ProgramExecute struct {
	Program        Program
	NamedFunctions []FunctionExecute
	ExecutableBlock
}

func (pe ProgramExecute) ExecuteProgram() ExecutionResult {
	execEnv := NewExecutionEnvironment()

	for _, fe := range pe.NamedFunctions {
		if execEnv.GlobalExists(*fe.NamedFunction.Name) {
			var pos lexer.Position
			if fe.NamedFunction.Name != nil {
				pos = fe.NamedFunction.Pos
			} else {
				pos = fe.UnnamedFunction.Pos
			}
			log.Fatalf("duplicate function %s at %s", *fe.NamedFunction.Name, pos)
		}
		execEnv.SetGlobal(*fe.NamedFunction.Name, fe)
	}

	return pe.ExecutableBlock.Execute(execEnv)
}

func (pe ProgramExecute) DumpProgram() {
	dump := pe.ExecutableBlock.ListRep()

	b, err := json.MarshalIndent(dump, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func (pe ProgramExecute) DumpFunctions() {
	functions := []any{}

	for _, fe := range pe.NamedFunctions {
		functions = append(functions, fe.ListRep())
	}

	b, err := json.MarshalIndent(functions, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func (a *Addition) Compile(typeMap TypeMap) (Executable, CompileErrors) {
	var errs CompileErrors

	ex := errs.Collect(a.Multiplication.Compile(typeMap))
	if len(a.Operations) == 0 {
		return ex, errs
	}

	plus := PlusOpMap[ex.Type(typeMap)]
	minus := MinusOpMap[ex.Type(typeMap)]

	operands := []Executable{}
	for _, opExpr := range a.Operations {
		operand := errs.Collect(opExpr.Operand.Compile(typeMap))
		operands = append(operands, operand)
	}

	atLeastOneUnknown := ex.Type(typeMap) == TypeUnknown || hasTypeUnknown(operands, typeMap)

	if atLeastOneUnknown {
		plus = PlusOpMap[TypeUnknown]
		minus = MinusOpMap[TypeUnknown]
	} else {
		for _, operand := range operands {
			if ex.Type(typeMap) != operand.Type(typeMap) {
				errs.Append(fmt.Errorf("type mismatch %s for %s at %s", ex.Type(typeMap), operand.Type(typeMap), a.Pos))
			}
		}
	}

	if plus == nil || minus == nil {
		errs.Append(fmt.Errorf("invalid type %s for +, - at %s", ex.Type(typeMap), a.Pos))
		return nil, errs
	}

	for i, operand := range operands {
		switch a.Operations[i].Op {
		case "+":
			ex = BinaryOperation{a.Pos, "+", plus, ex, operand, ex.Type(typeMap)}
		case "-":
			ex = BinaryOperation{a.Pos, "-", minus, ex, operand, ex.Type(typeMap)}
		default:
			errs.Append(fmt.Errorf("invalid operator %s for type %s at %s", a.Operations[i].Op, ex.Type(typeMap), a.Pos))
		}
	}

	return ex, errs
}

// very evil. very very bad.
var currentPos lexer.Position

type BinaryOperation struct {
	Pos lexer.Position

	Name        string
	Func        func(any, any) any
	Left, Right Executable
	TypeVal     Type
}

func (bo BinaryOperation) Execute(ee *ExecutionEnvironment) ExecutionResult {
	currentPos = bo.Pos
	return bo.Func(bo.Left.Execute(ee), bo.Right.Execute(ee))
}

func (bo BinaryOperation) Type(typeMap TypeMap) Type {
	return bo.TypeVal
}

func (bo BinaryOperation) ListRep() []any {
	return []any{bo.Name, bo.Left.ListRep(), bo.Right.ListRep()}
}

type ComparisonOperation struct {
	Pos lexer.Position

	Name        string
	Func        func(any, any) BoolValue
	Left, Right Executable
}

func (co ComparisonOperation) Execute(ee *ExecutionEnvironment) ExecutionResult {
	currentPos = co.Pos
	return co.Func(co.Left.Execute(ee), co.Right.Execute(ee))
}

func (co ComparisonOperation) Type(typeMap TypeMap) Type {
	return TypeBool
}

func (co ComparisonOperation) ListRep() []any {
	return []any{co.Name, co.Left.ListRep(), co.Right.ListRep()}
}

type ShortCircuitAnd struct {
	Left, Right Executable
}

// Execute returns the result of the left expression if it is false, otherwise it returns the result of the right expression.
func (sca ShortCircuitAnd) Execute(ee *ExecutionEnvironment) ExecutionResult {
	leftValue := sca.Left.Execute(ee)
	if !bool(leftValue.(BoolValue)) {
		return BoolValue(false)
	}
	return sca.Right.Execute(ee)
}

func (sca ShortCircuitAnd) Type(typeMap TypeMap) Type {
	return TypeBool
}

func (sca ShortCircuitAnd) ListRep() []any {
	return []any{"&&", sca.Left.ListRep(), sca.Right.ListRep()}
}

type ShortCircuitOr struct {
	Left, Right Executable
}

// Execute returns the result of the left expression if it is true, otherwise it returns the result of the right expression.
func (sco ShortCircuitOr) Execute(ee *ExecutionEnvironment) ExecutionResult {
	leftValue := sco.Left.Execute(ee)
	if bool(leftValue.(BoolValue)) {
		return BoolValue(true)
	}
	return sco.Right.Execute(ee)
}

func (sco ShortCircuitOr) Type(typeMap TypeMap) Type {
	return TypeBool
}

func (sco ShortCircuitOr) ListRep() []any {
	return []any{"||", sco.Left.ListRep(), sco.Right.ListRep()}
}
