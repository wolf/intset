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

func New(size int, values ...int) *IntSet {
	result := &IntSet{0, make([]int, size), make([]int, size)}
	return result.Add(values...)
}

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

func (is *IntSet) Empty() bool {
	return is.length == 0
}

func (is *IntSet) Length() int {
	return is.length
}

func (is *IntSet) Capacity() int {
	return len(is.keys)
}

func (is *IntSet) Contains(value int) bool {
	return value < len(is.keys) && is.keys[value] < is.length && is.set[is.keys[value]] == value
}

func (is *IntSet) Values() []int {
	result := make([]int, is.length)
	copy(result, is.set[:is.length])
	return result
}

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

func (is *IntSet) Remove(values ...int) *IntSet {
	for _, v := range values {
		if is.Contains(v) {
			kMoving := is.set[is.length-1]
			kToBeReplaced := is.keys[v]
			is.set[kToBeReplaced] = kMoving
			is.keys[kMoving] = kToBeReplaced
			is.length--
		}
	}
	return is
}

func (is *IntSet) Union(rhs *IntSet) *IntSet {
	return is.Add(rhs.set[:rhs.length]...)
}

func Union(lhs, rhs *IntSet) *IntSet {
	return lhs.Copy().Union(rhs)
}

func (is *IntSet) Difference(rhs *IntSet) *IntSet {
	return is.Remove(rhs.set[:rhs.length]...)
}

func Difference(lhs, rhs *IntSet) *IntSet {
	return lhs.Copy().Difference(rhs)
}

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
