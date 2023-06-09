package util

//
// IntSet is a set of integers. Defined as map[int]struct{}
//
type IntSet map[int]struct{}

// Ints constructs IntSet from slice of integers.
func Ints(values []int) IntSet {
	set := make(IntSet)

	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}

func OfInts(ints ...int) IntSet {
	return Ints(ints)
}

// Len returns number of elements in the receiver.
func (set IntSet) Len() int {
	return len(set)
}

// AsSlice returns the receiver's elements as slice of ints
func (set IntSet) AsSlice() []int {
	i := 0
	result := make([]int, len(set), len(set))

	for v := range set {
		result[i] = v
		i++
	}

	return result
}

func (set IntSet) Includes(v int) bool {
	_, ok := set[v]

	return ok
}

func (set IntSet) IncludesAny(values []int) bool {
	for _, v := range values {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}

func (set IntSet) HasIntersect(other IntSet) bool {
	for v := range other {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}

func (set IntSet) Intersect(other IntSet) IntSet {
	result := make(IntSet)

	for v := range set {
		if _, ok := other[v]; ok {
			result[v] = struct{}{}
		}
	}

	return result
}

func (set IntSet) Union(other IntSet) IntSet {
	result := make(IntSet)

	for v := range set {
		result[v] = struct{}{}
	}

	for v := range other {
		result[v] = struct{}{}
	}

	return result
}

func (set IntSet) Add(v int) {
	set[v] = struct{}{}
}

func (set IntSet) Remove(v int) {
	delete(set, v)
}

func (set IntSet) AddAll(values []int) {
	for _, v := range values {
		set[v] = struct{}{}
	}
}

func (set IntSet) AddSet(src IntSet) {
	for v := range src {
		set[v] = struct{}{}
	}
}

func (set IntSet) RemoveAll(values []int) {
	for _, v := range values {
		delete(set, v)
	}
}
