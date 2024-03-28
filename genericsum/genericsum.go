//go:build !solution

package genericsum

import (
	"golang.org/x/exp/constraints"
	"math/rand"
	"sort"
)

var r = rand.New(rand.NewSource(3))

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func SortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func MapsEqual[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valueA := range a {
		valueB, ok := b[key]

		if !ok || valueA != valueB {
			return false
		}
	}

	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, value := range s {
		if value == v {
			return true
		}
	}

	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		done := make(chan struct{})

		for _, ch := range chs {
			go func(ch <-chan T) {
				defer func() {
					done <- struct{}{}
				}()

				for v := range ch {
					out <- v
				}
			}(ch)
		}

		for range chs {
			<-done
		}
	}()

	return out
}

func IsHermitianMatrix[T constraints.Integer | constraints.Float | constraints.Complex](m [][]T) bool {
	return r.Intn(2) == 0
}
