package lo

// MIT License
//
// Copyright (c) 2022 Samuel Berthe

func Ternary[T any](condition bool, ifOutput, elseOutput T) T {
	if condition {
		return ifOutput
	}

	return elseOutput
}

func TernaryF[T any](condition bool, ifFunc, elseFunc func() T) T {
	if condition {
		return ifFunc()
	}

	return elseFunc()
}

type IfElse[T any] struct {
	result T
	done   bool
}

func If[T any](condition bool, result T) *IfElse[T] {
	if condition {
		return &IfElse[T]{result, true}
	}

	var t T
	return &IfElse[T]{t, false}
}

func IfF[T any](condition bool, resultF func() T) *IfElse[T] {
	if condition {
		return &IfElse[T]{resultF(), true}
	}

	var t T
	return &IfElse[T]{t, false}
}

func (i *IfElse[T]) ElseIf(condition bool, result T) *IfElse[T] {
	if !i.done && condition {
		i.result = result
		i.done = true
	}

	return i
}

func (i *IfElse[T]) ElseIfF(condition bool, resultF func() T) *IfElse[T] {
	if !i.done && condition {
		i.result = resultF()
		i.done = true
	}

	return i
}

func (i *IfElse[T]) Else(result T) T {
	if i.done {
		return i.result
	}

	return result
}

func (i *IfElse[T]) ElseF(resultF func() T) T {
	if i.done {
		return i.result
	}

	return resultF()
}

type SwitchCase[T comparable, R any] struct {
	predicate T
	result    R
	done      bool
}

func Switch[T comparable, R any](predicate T) *SwitchCase[T, R] {
	var result R

	return &SwitchCase[T, R]{
		predicate,
		result,
		false,
	}
}

func (s *SwitchCase[T, R]) Case(val T, result R) *SwitchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = result
		s.done = true
	}

	return s
}

func (s *SwitchCase[T, R]) CaseF(val T, cb func() R) *SwitchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = cb()
		s.done = true
	}

	return s
}

func (s *SwitchCase[T, R]) Default(result R) R {
	if !s.done {
		s.result = result
	}

	return s.result
}

func (s *SwitchCase[T, R]) DefaultF(cb func() R) R {
	if !s.done {
		s.result = cb()
	}

	return s.result
}
