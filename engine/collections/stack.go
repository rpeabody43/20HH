package collections

// A Stack ADT with a predetermined max size
type ArrayStack[T any] struct {
	data []T
	ptr  int
}

func NewArrayStack[T any](maxSize int) ArrayStack[T] {
	return ArrayStack[T]{
		data: make([]T, maxSize),
		ptr:  -1,
	}
}

func (s *ArrayStack[T]) Push(data T) {
	s.ptr++
	// Failsafe for pushing off end of array
	if s.ptr == len(s.data) {
		s.data = append(s.data, data)
	} else {
		s.data[s.ptr] = data
	}
}

func (s *ArrayStack[T]) Pop() T {
	toReturn := s.data[s.ptr]
	s.ptr--
	return toReturn
}
