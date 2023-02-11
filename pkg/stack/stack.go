package stack

import (
	"errors"
)

var ErrEmptyStack error = errors.New("empty stack")

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() Stack[T] {
	items := make([]T, 0)
	return Stack[T]{items: items}
}

func (s *Stack[T]) Push(item *T) {
	s.items = append(s.items, *item)
}

func (s *Stack[T]) Pop() (*T, error) {
	if len(s.items) == 0 {
		return nil, ErrEmptyStack
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]

	return &item, nil
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}
