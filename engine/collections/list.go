package collections

// A List ADT implemented as a circular array with a set size
// Only required functions are implemented -> there's no remove
type List[T any] struct {
	data    []T
	start   int
	end     int
	head    int
	size    int
	maxSize int
}

func NewList[T any](maxSize int) List[T] {
	return List[T]{
		data:    make([]T, maxSize),
		start:   maxSize / 2,
		end:     maxSize/2 + 1,
		head:    maxSize/2 + 1,
		size:    0,
		maxSize: maxSize,
	}
}

// Unpredictable behavior when inserting past maxSize
func (list *List[T]) AddAtTail(data T) {
	list.data[list.end] = data
	list.end = (list.end + 1) % list.maxSize
	list.size++
}

func (list *List[T]) AddAtHead(data T) {
	list.data[list.start] = data
	list.start--
	if list.start < 0 {
		list.start = list.maxSize - 1
	}
	list.size++
	// head is initially placed assuming AddAtTail is first
	// this corrects for that
	if list.size == 1 {
		list.head = list.start + 1
	}
}

func (list *List[T]) Size() int {
	return list.size
}

func (list *List[T]) Get(idx int) T {
	return list.data[(list.head+idx)%list.maxSize]
}

type ListIterator[T any] struct {
	list *List[T]
	idx  int
}

func (list *List[T]) Iterator() ListIterator[T] {
	return ListIterator[T]{
		list: list,
		idx:  -1,
	}
}

func (iter *ListIterator[T]) HasNext() bool {
	return iter.idx < iter.list.size-1
}

func (iter *ListIterator[T]) Next() T {
	iter.idx++
	return iter.list.Get(iter.idx)
}
