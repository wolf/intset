/*
Package intset implements a small set of integers.

The integers allowed to be members of the set range from 0 up to but not
including the size of the set's "universe".

This set implementation is efficient in time, but uses plenty of space.

How it works:

An IntSet contains two slices of int.  One represents the set itself.  For an
IntSet, is, the first is.length slots of the set slice are the actual values
in the set.  The physical length of that slice is permanently set at the size
of the set's universe.  That is, for a set that can hold integers from 0 up to
but not including 52, len(is.set) == 52; but only the first is.length slots
contain anything useful.  The is.keys slice is what provides random access to
the members of the set.  The keys slice is indexed by the set element, and
contains the index of that element within the set slice.  So for instance, if
7 is in the set, then is.keys[7] contains an index into the set slice.  If
that index is less than is.length, then maybe 7 is in the set.  We just have
to look in the set slot indexed by is.keys[7].  If that value is 7, then 7 is
in the set.

The only real problem that can occur is trying to add a value to a set that is
outside that set's universe.  There are at least three ways to handle that
situation: (1) ignore it, don't add the value; (2) return an error; (3) grow
the set's universe to include the new value, then add it.  This choice also
influences another design decision.  One useful design property is the ability
to "chain" operations, e.g., intset.Difference(setA, setB).Choose().  To allow
this to happen, most routines need to return the receiver (and nothing else).
For an error like this, the benefits of ignoring it may exceed the benefits of
reporting it.  The size of a set's universe is conceptually part of the type
of the set.  Trying to add a value outside that universe is a programming
error, not a run-time error.
*/

package intset

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type IntSet struct {
	length int
	keys   []int
	set    []int
}

func New(universeSize int, values ...int) *IntSet {
	result := &IntSet{0, make([]int, universeSize), make([]int, universeSize)}
	return result.Add(values...)
}

// Copy duplicates a set.

// The new set contains exactly the same members as the receiver set.  Copy is
// O(n) where n is the number of elements actually in the receiver set.
func (is *IntSet) Copy() *IntSet {
	return New(len(is.keys), is.set[:is.length]...)
}

func (is *IntSet) String() string {
	result := "["
	for i, v := range is.set[:is.length] {
		result += fmt.Sprint(v)
		if i < is.length-1 {
			result += " "
		}
	}
	result += "]"
	return result
}

func (is *IntSet) Length() int {
	return is.length
}

func (is *IntSet) Empty() bool {
	return is.length == 0
}

// UniverseSize is 1 greater than the largest value that can be a member of the set.
func (is *IntSet) UniverseSize() int {
	return len(is.keys)
}

// Contains returns true if the given value is a member of the receiver set.

// Contains is O(1).
func (is *IntSet) Contains(value int) bool {
	if value < 0 || len(is.keys) <= value { // is value within the set universe?
		return false // if not, then value can't possibly be in the set
	}
	setSlot := is.keys[value]
	return setSlot < is.length && is.set[setSlot] == value
}

// Values returns a new slice containing (in no particular order) all the elements in the receiver set.

// Values is O(n), n the number of elements in the receiver set.
func (is *IntSet) Values() []int {
	result := make([]int, is.length)
	copy(result, is.set[:is.length])
	return result
}

// Add puts a new element into the receiver set, if that element is in the set's universe.

// Add is O(1) for a single value and O(n), n the number of values to add, for a list of values.
// Add returns the receiver set to allow method chaining.
func (is *IntSet) Add(values ...int) *IntSet {
	for _, v := range values {
		if 0 <= v && v < len(is.keys) && !is.Contains(v) {
			is.keys[v] = is.length
			is.set[is.length] = v
			is.length++
		}
	}
	return is
}

// Remove removes elements from the receiver set.

// Remove is O(1) for a single value and O(n), n the number of values to remove, for a list of values.
// Remove returns the receiver set to allow method chaining.
func (is *IntSet) Remove(values ...int) *IntSet {
	for _, v := range values {
		if is.length == 0 {
			break
		}
		if is.Contains(v) {
			valueToMove := is.set[is.length-1]
			if valueToMove != v {
				valueToReplace := is.keys[v]
				is.set[valueToReplace] = valueToMove
				is.keys[valueToMove] = valueToReplace
			}
			is.length--
		}
	}
	return is
}

// Union with a receiver updates the receiver set to also include all the elements in other.

// Union is O(n), n the number of elements in other.
// Union returns the receiver set to allow method chaining.
func (is *IntSet) Union(other *IntSet) *IntSet {
	return is.Add(other.set[:other.length]...)
}

// Union with two arguments produces a new set that contains exactly all the elements in both lhs and rhs.

// Union is O(n+m), n the number of elements in lhs, m the number of elements in rhs.
func Union(lhs, rhs *IntSet) *IntSet {
	return lhs.Copy().Union(rhs)
}

// Difference with a receiver updates the receiver set to remove any element that also appears in other.

// Difference is O(n), n the number of elements in other.
// Difference returns the receiver set to allow method chaining.
func (is *IntSet) Difference(other *IntSet) *IntSet {
	return is.Remove(other.set[:other.length]...)
}

// Difference with two arguments produces a new set that contains exactly all the elements in lhs that do not appear in rhs.

// Difference is O(n+m), n the number of elements in lhs, m the number of elements in rhs
func Difference(lhs, rhs *IntSet) *IntSet {
	return lhs.Copy().Difference(rhs)
}

// Choose returns an element at random from the receiver set.  The set itself is not modified.

// Choose is the order of complexity of rand.Int
func (is *IntSet) Choose() (choice int, err error) {
	if is.Empty() {
		return 0, nil // error("set is empty")
	}
	max := big.NewInt(int64(is.length))
	ip, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	choice = is.set[ip.Int64()]
	return
}
