package utils

import "reflect"

// ============================================================== interfaces ===

// RoProp const property interface.
type RoProp[T any] interface {
	Get() T
	IsNone() bool
}

// RwProp variable property interface
type RwProp[T any] interface {
	RoProp[T]
	Set(value T)
	IsZeroable() bool
	ConstRef() RoProp[T]
}

// PutProp putter variable property interface. The putter allows
// the client to define assignment using a client defined function.
// The putter will still set the property's Field value.
type PutProp[T any] interface {
	RwProp[T]
	Put(value T)
}

// ===================================================================== New ===

// NewRoProp create const property
func NewRoProp[T any](value T) RoProp[T] {
	return &constProp[T]{field: value}
}

// NewRwProp create variable property
func NewRwProp[T any](value T) RwProp[T] {
	return &VarProp[T]{Field: value}
}

// NewRwProp create variable and zeroable property
func NewRwPropZ[T any](value T) RwProp[T] {
	return &VarProp[T]{Field: value, zeroable: true}
}

// NewPutProp create putter variable property
func NewPutProp[T any](value T, putter func(value T)) PutProp[T] {
	return &putVarProp[T]{
		VarProp[T]{Field: value},
		putter,
	}
}

// NewPutProp create putter variable and zeroable property
func NewPutPropZ[T any](value T, putter func(value T)) PutProp[T] {
	return &putVarProp[T]{
		VarProp[T]{Field: value, zeroable: true},
		putter,
	}
}

// =============================================================== factories ===

// RoPropFactory const property factory
type RoPropFactory[T any] struct {
	Zeroable bool
}

// Construct const property constructor
func (f RoPropFactory[T]) Construct(value T) RoProp[T] {
	return &constProp[T]{field: value}
}

// RwPropFactory variable property factory
type RwPropFactory[T any] struct {
	Zeroable bool
}

// Construct variable property constructor
func (f RwPropFactory[T]) Construct(value T) RwProp[T] {
	return &VarProp[T]{Field: value, zeroable: f.Zeroable}
}

// PutPropFactory putter variable property factory
type PutPropFactory[T any] struct {
	Zeroable bool
}

// Construct putter variable property constructor
func (f PutPropFactory[T]) Construct(value T, putter func(value T)) PutProp[T] {
	return &putVarProp[T]{
		VarProp[T]{Field: value, zeroable: f.Zeroable},
		putter,
	}
}

// =============================================================== constProp ===

type constProp[T any] struct {
	// constProp is not exported to prevent the client from
	// creating a const property without a value, which is
	// of no practical use as it can't be set later on.
	// The client should use either RoPropFactory or NewRoProp.
	//
	field T
}

// Get property value getter
func (p *constProp[T]) Get() T {
	return p.field
}

// IsNone determines whether the property has a value set
func (p *constProp[T]) IsNone() bool {
	return isPropNil(p.field, false)
}

// ================================================================= VarProp ===

// VarProp a read/write property
type VarProp[T any] struct {
	zeroable bool
	Field    T
}

// Get property value getter
func (p *VarProp[T]) Get() T {
	return p.Field
}

// Set property value setter
func (p *VarProp[T]) Set(value T) {
	p.Field = value
}

// IsNone determines whether the property has a value set
func (p *VarProp[T]) IsNone() bool {
	return isPropNil(p.Field, p.zeroable)
}

// IsZeroable indicates whether a zero value is a valid value for this property
func (p *VarProp[T]) IsZeroable() bool {
	return p.zeroable
}

// ConstRef returns a read only reference to this property
func (p *VarProp[T]) ConstRef() RoProp[T] {
	return p
}

// ============================================================== putVarProp ===

type putVarProp[T any] struct {
	VarProp[T]
	putter func(value T)
}

// Put set the value of the property and invoke the putter function
func (p *putVarProp[T]) Put(value T) {
	p.Set(value)
	p.putter(value)
}

// ==================================================================== misc ===

func isPropNil[T any](value T, zeroable bool) bool {
	refV := reflect.ValueOf(value)
	refK := refV.Kind()

	if refK == 0 {
		// value has not been set yet
		//
		return true
	}

	switch refK {
	case
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.Slice,
		reflect.UnsafePointer:
		return IsNil(value)
	}
	return !zeroable && refV.IsZero()
}
