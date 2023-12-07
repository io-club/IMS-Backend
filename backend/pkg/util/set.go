package util

import (
	"encoding/json"
	"fmt"
)

// 集合
type Set[T comparable] map[T]struct{}

func NewSet[T comparable](ts ...T) Set[T] {
	set := Set[T]{}
	set.Add(ts...)
	return set
}

func (s Set[T]) Add(vs ...T) {
	for _, v := range vs {
		s[v] = struct{}{}
	}
}

func (s Set[T]) Remove(vs ...T) {
	for _, v := range vs {
		delete(s, v)
	}
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) ToSlice() []T {
	slice := []T{}
	for v := range s {
		slice = append(slice, v)
	}
	return slice
}

func (s Set[T]) ToMap() map[T]struct{} {
	return s
}

func (s Set[T]) Clear() {
	for v := range s {
		delete(s, v)
	}
}

func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Set[T]) Clone() Set[T] {
	clone := NewSet[T]()
	for v := range s {
		clone.Add(v)
	}
	return clone
}

func (s Set[T]) Equal(other Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for v := range s {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	union := s.Clone()
	for v := range other {
		union.Add(v)
	}
	return union
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	intersect := NewSet[T]()
	for v := range s {
		if other.Contains(v) {
			intersect.Add(v)
		}
	}
	return intersect
}

func (s Set[T]) Difference(other Set[T]) Set[T] {

	difference := NewSet[T]()
	for v := range s {
		if !other.Contains(v) {
			difference.Add(v)
		}
	}
	return difference
}

func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	return s.Difference(other).Union(other.Difference(s))
}

func (s Set[T]) IsSubset(other Set[T]) bool {
	if s.Len() > other.Len() {
		return false
	}
	for v := range s {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

func (s Set[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s Set[T]) IsProperSubset(other Set[T]) bool {
	return s.IsSubset(other) && !s.Equal(other)
}

func (s Set[T]) IsProperSuperset(other Set[T]) bool {
	return s.IsSuperset(other) && !s.Equal(other)
}

func (s Set[T]) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

func (s Set[T]) UnmarshalJSON(data []byte) error {
	slice := []T{}
	err := json.Unmarshal(data, &slice)
	if err != nil {
		return err
	}
	for _, v := range slice {
		s.Add(v)
	}
	return nil
}

func (s Set[T]) MarshalYAML() (interface{}, error) {
	return s.ToSlice(), nil
}

func (s Set[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	slice := []T{}
	err := unmarshal(&slice)
	if err != nil {
		return err
	}
	for _, v := range slice {
		s.Add(v)
	}
	return nil
}
