package collections

import (
	"testing"
)

func is0To9(data []int) bool {
	for i := 0; i <= 9; i++ {
		if data[i] != i {
			return false
		}
	}
	return true
}

func TestListIterator(t *testing.T) {
	testList := NewList[int](10)
	for i := 0; i <= 6; i++ {
		testList.AddAtTail(i)
	}
	for i := 9; i >= 7; i-- {
		testList.AddAtHead(i)
	}

	var out []int
	iter := testList.Iterator()
	for iter.HasNext() {
		val := iter.Next()
		out = append(out, val)
	}

	if !is0To9(out) {
		t.Errorf("Iterator returns %v instead of 0-9", out)
		t.Errorf("List has size of %d", testList.Size())
		t.Errorf("List has data of %v", testList.data)
	}
}

func TestArrayStack(t *testing.T) {
	testStack := NewArrayStack[int](10)
	for i := 4; i >= 0; i-- {
		testStack.Push(i)
	}
	var out []int
	for i := 0; i < 5; i++ {
		out = append(out, testStack.Pop())
	}

	for i := 9; i >= 5; i-- {
		testStack.Push(i)
	}
	for i := 0; i < 5; i++ {
		out = append(out, testStack.Pop())
	}
	if !is0To9(out) {
		t.Errorf("Stack returns %v instead of 0-9", out)
		t.Errorf("Stack has data of %v", testStack.data)
	}
}
