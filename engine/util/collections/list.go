package collections

// A List ADT implemented as a circular array with a set size
// Only required functions are implemented -> there's no remove
type List[T any] struct {
	data    []T // Underlying array
	head    int // First cell with data
	end     int // Cell after data
	size    int // Size of current list
	maxSize int // How big the list can get before overwriting itself
}

func NewList[T any](maxSize int) List[T] {
	return List[T]{
		data:    make([]T, maxSize),
		head:    maxSize/2 + 1,
		end:     maxSize/2 + 1,
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
	list.size++

	list.head--
	if list.head < 0 {
		list.head = list.maxSize - 1
	}
	list.data[list.head] = data
}

func (list *List[T]) Size() int {
	return list.size
}

func (list *List[T]) Get(idx int) T {
	return list.data[(list.head+idx)%list.maxSize]
}

func (list *List[T]) Set(idx int, val T) {
	list.data[(list.head+idx)%list.maxSize] = val
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
