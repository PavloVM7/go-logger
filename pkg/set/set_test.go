package set

import (
	"sort"
	"testing"
)

//revive:disable

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
func TestSet_ForEach_remove(t *testing.T) {
	s := NewWithValues[int](1, 2, 3, 4, 5)
	assertEqual(t, 5, s.Len())
	s.ForEach(func(i int) {
		s.Remove(i)
	})
	assertEqual(t, 0, s.Len())
}

func TestSet_ForEach(t *testing.T) {
	s := NewWithValues[int](1, 2, 3, 4, 5)
	sum := 0
	s.ForEach(func(i int) {
		sum += i
	})
	assertEqual(t, 15, sum)
}

func TestSet_Clear(t *testing.T) {
	s := NewCapacity[int](5)
	s.AddAll(1, 2, 3, 4, 5)
	assertEqual(t, 5, s.Len())
	s.Clear()
	assertEqual(t, 0, s.Len())
	if !s.IsEmpty() {
		t.Error("Expected set to be empty after Clear()")
	}
}

func TestSet_Size(t *testing.T) {
	s := NewCapacity[int](5)
	s.AddAll(1, 2, 3, 4, 5)
	assertEqual(t, 5, s.Len())
	assertEqual(t, 5, s.Size())
	s.Clear()
	assertEqual(t, 0, s.Len())
	assertEqual(t, 0, s.Size())
}

func TestSet_TrimToLen(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Remove(1)
	s.Remove(2)
	s.Remove(3)
	s.TrimToLen()
	sl := s.ToSlice()
	assertEqual(t, 2, s.Len())
	sort.Ints(sl)
	expected := []int{4, 5}
	if len(sl) != len(expected) {
		t.Errorf("Expected slice length %d, got %d", len(expected), len(sl))
	}
	for i, v := range expected {
		if i >= len(sl) || sl[i] != v {
			t.Errorf("At index %d: expected %v, got %v", i, v, sl[i])
		}
	}
}
