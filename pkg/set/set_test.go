package set

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := New[int]()
	require.NotNil(t, s)
	assert.Equal(t, 0, s.Len())
	assert.True(t, s.IsEmpty())
}

func TestNewCapacity(t *testing.T) {
	s := NewCapacity[int](10)
	require.NotNil(t, s)
	assert.Equal(t, 0, s.Len())

	s.AddAll(1, 2, 3)
	assert.Equal(t, 3, s.Len())
}

func TestNewWithValues(t *testing.T) {
	s := NewWithValues(1, 2, 3, 4, 5)
	require.NotNil(t, s)
	assert.Equal(t, 5, s.Len())
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(5))
	assert.False(t, s.Contains(6))
}

func TestNewWithValues_Duplicates(t *testing.T) {
	s := NewWithValues(1, 2, 2, 3, 3, 3)
	assert.Equal(t, 3, s.Len())
}

func TestNewWithValues_Empty(t *testing.T) {
	s := NewWithValues[int]()
	assert.Equal(t, 0, s.Len())
	assert.True(t, s.IsEmpty())
}

func TestSet_Add(t *testing.T) {
	s := New[int]()

	// Add new element
	added := s.Add(1)
	assert.True(t, added)
	assert.True(t, s.Contains(1))
	assert.Equal(t, 1, s.Len())

	// Add duplicate
	added = s.Add(1)
	assert.False(t, added)
	assert.Equal(t, 1, s.Len())

	// Add another element
	added = s.Add(2)
	assert.True(t, added)
	assert.Equal(t, 2, s.Len())
}

func TestSet_AddAll(t *testing.T) {
	s := New[int]()

	count := s.AddAll(1, 2, 3)
	assert.Equal(t, 3, count)
	assert.Equal(t, 3, s.Len())

	// Add with duplicates
	count = s.AddAll(1, 2, 4)
	assert.Equal(t, 1, count)
	assert.Equal(t, 4, s.Len())

	// Add empty
	count = s.AddAll()
	assert.Equal(t, 0, count)
	assert.Equal(t, 4, s.Len())
}

func TestSet_Remove(t *testing.T) {
	s := NewWithValues(1, 2, 3)

	s.Remove(2)
	assert.Equal(t, 2, s.Len())
	assert.False(t, s.Contains(2))

	// Remove non-existent element (should not panic)
	s.Remove(4)
	assert.Equal(t, 2, s.Len())
}

func TestSet_Contains(t *testing.T) {
	s := NewWithValues("a", "b", "c")

	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.True(t, s.Contains("c"))
	assert.False(t, s.Contains("d"))
}

func TestSet_Len(t *testing.T) {
	s := New[int]()
	assert.Equal(t, 0, s.Len())

	s.Add(1)
	assert.Equal(t, 1, s.Len())

	s.AddAll(2, 3, 4)
	assert.Equal(t, 4, s.Len())

	s.Remove(1)
	assert.Equal(t, 3, s.Len())
}

func TestSet_IsEmpty(t *testing.T) {
	s := New[int]()
	assert.True(t, s.IsEmpty())

	s.Add(1)
	assert.False(t, s.IsEmpty())

	s.Remove(1)
	assert.True(t, s.IsEmpty())
}

func TestSet_Clear(t *testing.T) {
	s := NewCapacity[int](5)
	s.AddAll(1, 2, 3, 4, 5)
	assert.Equal(t, 5, s.Len())

	s.Clear()
	assert.Equal(t, 0, s.Len())
	assert.True(t, s.IsEmpty())

	// Clear empty set (should not panic)
	s.Clear()
	assert.Equal(t, 0, s.Len())
}

func TestSet_ToSlice(t *testing.T) {
	s := NewWithValues(3, 1, 4, 1, 5, 9, 2, 6)
	slice := s.ToSlice()

	assert.Equal(t, 7, len(slice))

	// Check all elements are present
	for _, elem := range slice {
		assert.True(t, s.Contains(elem))
	}
}

func TestSet_ToSlice_Empty(t *testing.T) {
	s := New[int]()
	slice := s.ToSlice()

	assert.Equal(t, 0, len(slice))
}

func TestSet_Copy(t *testing.T) {
	s := NewWithValues(1, 2, 3)
	cpy := s.Copy()

	assert.Equal(t, s.Len(), cpy.Len())
	assert.True(t, cpy.Contains(1))
	assert.True(t, cpy.Contains(2))
	assert.True(t, cpy.Contains(3))

	// Modify original
	s.Add(4)
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, 3, cpy.Len())
	assert.False(t, cpy.Contains(4))
}

func TestSet_Copy_Empty(t *testing.T) {
	s := New[int]()
	cpy := s.Copy()

	assert.Equal(t, 0, cpy.Len())
	assert.True(t, cpy.IsEmpty())
}

func TestSet_TrimToLen(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)

	// Remove some elements to create capacity overhead
	s.Remove(1)
	s.Remove(2)
	s.Remove(3)

	// Before trim, verify set state
	assert.Equal(t, 2, s.Len())
	assert.True(t, s.Contains(4))
	assert.True(t, s.Contains(5))

	// Trim should not change contents
	s.TrimToLen()

	sl := s.ToSlice()
	assert.Equal(t, 2, s.Len())
	slices.Sort(sl)
	expected := []int{4, 5}
	assert.Equal(t, expected, sl)
}

func TestSet_TrimToLen_Empty(t *testing.T) {
	s := New[int]()
	s.TrimToLen()
	assert.Equal(t, 0, s.Len())
}

func TestSet_TrimToLen_NoRemovals(t *testing.T) {
	s := NewWithValues(1, 2, 3)
	originalLen := s.Len()
	s.TrimToLen()
	assert.Equal(t, originalLen, s.Len())
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
	assert.True(t, s.Contains(3))
}

func TestSet_TrimToLen_Nil(t *testing.T) {
	var s Set[int]
	s.TrimToLen()
	// Should not panic
	t.Cleanup(func() {
		// Cleanup to use t parameter
	})
}

func TestSet_Size(t *testing.T) {
	s := NewCapacity[int](5)
	s.AddAll(1, 2, 3, 4, 5)
	assert.Equal(t, 5, s.Len())
	assert.Equal(t, 5, s.Size())

	s.Clear()
	assert.Equal(t, 0, s.Len())
	assert.Equal(t, 0, s.Size())
}

func TestSet_ForEach(t *testing.T) {
	s := NewWithValues(1, 2, 3, 4, 5)
	sum := 0
	s.ForEach(func(i int) {
		sum += i
	})
	assert.Equal(t, 15, sum)
}

func TestSet_ForEach_Empty(t *testing.T) {
	s := New[int]()
	callCount := 0
	s.ForEach(func(_ int) {
		callCount++
	})
	assert.Equal(t, 0, callCount)
}

func TestSet_ForEach_Remove(t *testing.T) {
	s := NewWithValues(1, 2, 3, 4, 5)
	assert.Equal(t, 5, s.Len())

	s.ForEach(func(i int) {
		s.Remove(i)
	})

	assert.Equal(t, 0, s.Len())
}

func TestSet_StringType(t *testing.T) {
	s := NewWithValues("apple", "banana", "cherry")

	assert.Equal(t, 3, s.Len())
	assert.True(t, s.Contains("apple"))
	assert.True(t, s.Contains("banana"))
	assert.True(t, s.Contains("cherry"))
	assert.False(t, s.Contains("orange"))

	s.Remove("banana")
	assert.Equal(t, 2, s.Len())
	assert.False(t, s.Contains("banana"))
}

func TestSet_StructType(t *testing.T) {
	type Point struct {
		X, Y int
	}

	s := New[Point]()
	p1 := Point{1, 2}
	p2 := Point{3, 4}

	s.Add(p1)
	s.Add(p2)

	assert.Equal(t, 2, s.Len())
	assert.True(t, s.Contains(p1))
	assert.True(t, s.Contains(p2))

	s.Add(p1) // Duplicate
	assert.Equal(t, 2, s.Len())
}

func TestSet_PointerType(t *testing.T) {
	s := New[*int]()

	a := 1
	b := 2
	c := 3

	s.Add(&a)
	s.Add(&b)
	s.Add(&c)

	assert.Equal(t, 3, s.Len())
	assert.True(t, s.Contains(&a))
	assert.True(t, s.Contains(&b))
	assert.True(t, s.Contains(&c))
}
