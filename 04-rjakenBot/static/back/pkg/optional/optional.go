package optional

import "fmt"

type Optional[T any] interface {
	HasValue() bool
	Value() T
	String() string
}

func Equal[T comparable](first Optional[T], second Optional[T]) bool {
	if !first.HasValue() && !second.HasValue() {
		return true
	}
	if first.HasValue() && second.HasValue() {
		return first.Value() == second.Value()
	}

	return false
}

///////////////////////////////////////////

type nilImpl[T any] struct{}

func (n nilImpl[T]) HasValue() bool {
	return false
}

func (n nilImpl[T]) Value() T {
	panic("optional value is empty")
}

func (n nilImpl[T]) String() string {
	return "Nil"
}

///////////////////////////////////////////

type someImpl[T any] struct {
	val T
}

func (s someImpl[T]) HasValue() bool {
	return true
}

func (s someImpl[T]) Value() T {
	return s.val
}

func (s someImpl[T]) String() string {
	return fmt.Sprint(s.val)
}

///////////////////////////////////////////

func Some[T any](val T) Optional[T] {
	return someImpl[T]{
		val: val,
	}
}

func Nil[T any]() Optional[T] {
	return nilImpl[T]{}
}
