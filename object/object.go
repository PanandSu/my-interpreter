package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"my-interpreter/ast"
	"strings"
)

const (
	NULL  = "NULL"
	ERROR = "ERROR"

	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	STRING  = "STRING"

	BUILTIN_FUNCTION = "BUILTIN_FUNCTION"

	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"

	MAP   = "MAP"
	ARRAY = "ARRAY"
)

type ObjectType string

// 哈希键结构体
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// 可哈希
type Hashable interface {
	Hash() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

// 整数
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Hash() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

// 布尔
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Hash() HashKey {
	var val uint64
	if b.Value {
		val = 1
	} else {
		val = 0
	}
	return HashKey{Type: b.Type(), Value: val}
}

// 字符串
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING
}
func (s *String) Inspect() string {
	return s.Value
}
func (s *String) Hash() HashKey {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Null
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}
func (n *Null) Inspect() string {
	return "null"
}

// 数组
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}
func (a *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type Pair struct {
	Key   Object
	Value Object
}

// 映射
type Map struct {
	Mappings map[HashKey]*Pair
}

func (m *Map) Type() ObjectType {
	return MAP
}
func (m *Map) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, v := range m.Mappings {
		pairs = append(pairs, v.Key.Inspect()+": "+v.Value.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// 错误
type Error struct {
	Msg string
}

func (e *Error) Type() ObjectType {
	return ERROR
}
func (e *Error) Inspect() string {
	return e.Msg
}

// 函数
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" ")
	out.WriteString(f.Body.String())
	return out.String()
}

// 返回值
type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return RETURN
}
func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

// 内置函数
type BuiltinFunction func(params ...Object) Object

type Builtins struct {
	Fn BuiltinFunction
}

func (b *Builtins) Type() ObjectType {
	return BUILTIN_FUNCTION
}
func (b *Builtins) Inspect() string {
	return "builtin function"
}
