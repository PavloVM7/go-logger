// Copyright Ⓒ 2025 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package set contains type and methods for working with set
package set

// Set is a set
type Set[T comparable] map[T]struct{}

// New creates a new set
func New[T comparable]() Set[T] {
	return make(map[T]struct{})
}

// NewCapacity creates a new set with capacity
func NewCapacity[T comparable](capacity int) Set[T] {
	return make(map[T]struct{}, capacity)
}

// NewWithValues creates a new set with values
func NewWithValues[T comparable](elems ...T) Set[T] {
	set := NewCapacity[T](len(elems))
	set.AddAll(elems...)
	return set
}

// ToSlice returns a slice of elements
func (s Set[T]) ToSlice() []T {
	keys := make([]T, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}

	return keys
}

// Clear clears the set
func (s Set[T]) Clear() {
	clear(s)
}

// Add adds an element
func (s Set[T]) Add(elem T) bool {
	if _, ok := s[elem]; ok {
		return false
	}
	s[elem] = struct{}{}
	return true
}

// AddAll adds multiple elements
func (s Set[T]) AddAll(elems ...T) int {
	count := 0
	for _, elem := range elems {
		if s.Add(elem) {
			count++
		}
	}

	return count
}

// Remove removes an element
func (s Set[T]) Remove(elem T) {
	delete(s, elem)
}

// Contains returns true if the element is in the set
func (s Set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}

// Len returns the length of the set
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty returns true if the set is empty
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Copy creates a copy of the set
func (s Set[T]) Copy() Set[T] {
	result := NewCapacity[T](len(s))
	for k, v := range s {
		result[k] = v
	}

	return result
}

// TrimToLen trims the set to the length
func (s Set[T]) TrimToLen() {
	cpy := s.Copy() //nolint:staticcheck
	r := &s
	*r = cpy
}

// Size returns the size of the set
func (s Set[T]) Size() int {
	return s.Len()
}

// ForEach iterates over the set and applies a function 'f' to each element of the set
func (s Set[T]) ForEach(f func(T)) {
	for k := range s {
		f(k)
	}
}
