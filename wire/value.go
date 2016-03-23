// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package wire

import (
	"fmt"
	"strings"
)

// Value holds the over-the-wire representation of a Thrift value.
//
// The Type of the value determines which field in the Value is valid.
type Value struct {
	typ Type

	tbool   bool
	tdouble float64
	ti8     int8
	ti16    int16
	ti32    int32
	ti64    int64
	tbinary []byte
	tstruct Struct
	tmap    Map
	tset    Set
	tlist   List
}

// Type retrieves the type of value inside a Value.
func (v *Value) Type() Type {
	return v.typ
}

// Get retrieves whatever value the given Value contains.
func (v *Value) Get() interface{} {
	switch v.typ {
	case TBool:
		return v.tbool
	case TI8:
		return v.ti8
	case TDouble:
		return v.tdouble
	case TI16:
		return v.ti16
	case TI32:
		return v.ti32
	case TI64:
		return v.ti64
	case TBinary:
		return v.tbinary
	case TStruct:
		return v.tstruct
	case TMap:
		return v.tmap
	case TSet:
		return v.tset
	case TList:
		return v.tlist
	default:
		panic(fmt.Sprintf("Unknown value type %v", v.typ))
	}
}

// NewValueBool constructs a new Value that contains a boolean.
func NewValueBool(v bool) Value {
	return Value{
		typ:   TBool,
		tbool: v,
	}
}

// GetBool gets the Bool value from a Value.
func (v *Value) GetBool() bool {
	return v.tbool
}

// NewValueI8 constructs a new Value that contains a byte
func NewValueI8(v int8) Value {
	return Value{
		typ: TI8,
		ti8: v,
	}
}

// GetI8 gets the I8 value from a Value.
func (v *Value) GetI8() int8 {
	return v.ti8
}

// NewValueDouble constructs a new Value that contains a double.
func NewValueDouble(v float64) Value {
	return Value{
		typ:     TDouble,
		tdouble: v,
	}
}

// GetDouble gets the Double value from a Value.
func (v *Value) GetDouble() float64 {
	return v.tdouble
}

// NewValueI16 constructs a new Value that contains a 16-bit integer.
func NewValueI16(v int16) Value {
	return Value{
		typ:  TI16,
		ti16: v,
	}
}

// GetI16 gets the I16 value from a Value.
func (v *Value) GetI16() int16 {
	return v.ti16
}

// NewValueI32 constructs a new Value that contains a 32-bit integer.
func NewValueI32(v int32) Value {
	return Value{
		typ:  TI32,
		ti32: v,
	}
}

// GetI32 gets the I32 value from a Value.
func (v *Value) GetI32() int32 {
	return v.ti32
}

// NewValueI64 constructs a new Value that contains a 64-bit integer.
func NewValueI64(v int64) Value {
	return Value{
		typ:  TI64,
		ti64: v,
	}
}

// GetI64 gets the I64 value from a Value.
func (v *Value) GetI64() int64 {
	return v.ti64
}

// NewValueBinary constructs a new Value that contains a binary string.
func NewValueBinary(v []byte) Value {
	return Value{
		typ:     TBinary,
		tbinary: v,
	}
}

// NewValueString constructs a new Value that contains a string.
func NewValueString(v string) Value {
	return Value{
		typ:     TBinary,
		tbinary: unsafeStringToBytes(v),
	}
}

// GetBinary gets the Binary value from a Value.
func (v *Value) GetBinary() []byte {
	return v.tbinary
}

// GetString gets a string value from a Value.
func (v *Value) GetString() string {
	return unsafeBytesToString(v.tbinary)
}

// NewValueStruct constructs a new Value that contains a struct.
func NewValueStruct(v Struct) Value {
	return Value{
		typ:     TStruct,
		tstruct: v,
	}
}

// GetStruct gets the Struct value from a Value.
func (v *Value) GetStruct() Struct {
	return v.tstruct
}

// NewValueMap constructs a new Value that contains a map.
func NewValueMap(v Map) Value {
	return Value{
		typ:  TMap,
		tmap: v,
	}
}

// GetMap gets the Map value from a Value.
func (v *Value) GetMap() Map {
	return v.tmap
}

// NewValueSet constructs a new Value that contains a set.
func NewValueSet(v Set) Value {
	return Value{
		typ:  TSet,
		tset: v,
	}
}

// GetSet gets the Set value from a Value.
func (v *Value) GetSet() Set {
	return v.tset
}

// NewValueList constructs a new Value that contains a list.
func NewValueList(v List) Value {
	return Value{
		typ:   TList,
		tlist: v,
	}
}

// GetList gets the List value from a Value.
func (v *Value) GetList() List {
	return v.tlist
}

func (v Value) String() string {
	switch v.typ {
	case TBool:
		return fmt.Sprintf("TBool(%v)", v.tbool)
	case TI8:
		return fmt.Sprintf("TI8(%v)", v.ti8)
	case TDouble:
		return fmt.Sprintf("TDouble(%v)", v.tdouble)
	case TI16:
		return fmt.Sprintf("TI16(%v)", v.ti16)
	case TI32:
		return fmt.Sprintf("TI32(%v)", v.ti32)
	case TI64:
		return fmt.Sprintf("TI64(%v)", v.ti64)
	case TBinary:
		return fmt.Sprintf("TBinary(%v)", v.tbinary)
	case TStruct:
		return fmt.Sprintf("TStruct(%v)", v.tstruct)
	case TMap:
		return fmt.Sprintf("TMap(%v)", v.tmap)
	case TSet:
		return fmt.Sprintf("TSet(%v)", v.tset)
	case TList:
		return fmt.Sprintf("TList(%v)", v.tlist)
	default:
		panic(fmt.Sprintf("Unknown value type %v", v.typ))
	}
}

// Struct provides a wire-level representation of a struct.
//
// At this level, structs don't have names or named fields.
type Struct struct {
	Fields []Field
}

func (s Struct) fieldMap() map[int16]Value {
	m := make(map[int16]Value, len(s.Fields))
	for _, f := range s.Fields {
		m[f.ID] = f.Value
	}
	return m
}

func (s Struct) String() string {
	fields := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		fields[i] = field.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
}

// Field is a single field inside a Struct.
type Field struct {
	ID    int16
	Value Value
}

func (f Field) String() string {
	return fmt.Sprintf("%v: %v", f.ID, f.Value)
}

// Set is a set of values.
type Set struct {
	ValueType Type
	Size      int
	Items     ValueList
}

func (s Set) String() string {
	items := make([]string, 0, s.Size)
	s.Items.ForEach(func(item Value) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf("[set]%v{%s}", s.ValueType, strings.Join(items, ", "))
}

// List is a list of values.
type List struct {
	ValueType Type
	Size      int
	Items     ValueList
}

func (l List) String() string {
	items := make([]string, 0, l.Size)
	l.Items.ForEach(func(item Value) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf("[]%v{%s}", l.ValueType, strings.Join(items, ", "))
}

// Map is a collection of key-value pairs.
type Map struct {
	KeyType   Type
	ValueType Type
	Size      int
	Items     MapItemList
}

func (m Map) String() string {
	items := make([]string, 0, m.Size)
	m.Items.ForEach(func(item MapItem) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf(
		"map[%v]%v{%s}", m.KeyType, m.ValueType, strings.Join(items, ", "),
	)
}

// MapItem is a single item in a Map.
type MapItem struct {
	Key   Value
	Value Value
}

func (mi MapItem) String() string {
	return fmt.Sprintf("%v: %v", mi.Key, mi.Value)
}
