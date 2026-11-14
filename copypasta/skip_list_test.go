package copypasta

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSkipListBasic(t *testing.T) {
	sl := newSkipList()
	if sl.len() != 0 {
		t.Fatal("expected empty list")
	}

	// insert and search
	for _, kv := range [][2]int{{1, 10}, {5, 50}, {3, 30}, {7, 70}, {2, 20}} {
		sl.insert(kv[0], kv[1])
	}
	if sl.len() != 5 {
		t.Fatalf("want len=5, got %d", sl.len())
	}

	for _, kv := range [][2]int{{1, 10}, {5, 50}, {3, 30}, {7, 70}, {2, 20}} {
		v, ok := sl.search(kv[0])
		if !ok || v != kv[1] {
			t.Fatalf("search(%d): got (%d,%v), want (%d,true)", kv[0], v, ok, kv[1])
		}
	}

	// update existing key
	sl.insert(3, 300)
	if v, ok := sl.search(3); !ok || v != 300 {
		t.Fatalf("update: want 300, got %d", v)
	}
	if sl.len() != 5 {
		t.Fatalf("update must not change len; got %d", sl.len())
	}

	// delete
	if !sl.delete(3) {
		t.Fatal("delete(3) should return true")
	}
	if _, ok := sl.search(3); ok {
		t.Fatal("deleted key must not be found")
	}
	if sl.delete(99) {
		t.Fatal("delete of absent key should return false")
	}
	if sl.len() != 4 {
		t.Fatalf("want len=4 after delete, got %d", sl.len())
	}
}

func TestSkipListBounds(t *testing.T) {
	sl := newSkipList()
	for _, k := range []int{2, 5, 8, 11, 14} {
		sl.insert(k, k*10)
	}

	// lowerBound
	if n := sl.lowerBound(5); n == nil || n.key != 5 {
		t.Fatalf("lowerBound(5): want 5, got %v", n)
	}
	if n := sl.lowerBound(6); n == nil || n.key != 8 {
		t.Fatalf("lowerBound(6): want 8, got %v", n)
	}
	if n := sl.lowerBound(15); n != nil {
		t.Fatalf("lowerBound(15): want nil, got %v", n)
	}

	// upperBound
	if n := sl.upperBound(5); n == nil || n.key != 8 {
		t.Fatalf("upperBound(5): want 8, got %v", n)
	}
	if n := sl.upperBound(14); n != nil {
		t.Fatalf("upperBound(14): want nil, got %v", n)
	}

	// floor
	if n := sl.floor(5); n == nil || n.key != 5 {
		t.Fatalf("floor(5): want 5, got %v", n)
	}
	if n := sl.floor(6); n == nil || n.key != 5 {
		t.Fatalf("floor(6): want 5, got %v", n)
	}
	if n := sl.floor(1); n != nil {
		t.Fatalf("floor(1): want nil, got %v", n)
	}
}

func TestSkipListOrder(t *testing.T) {
	sl := newSkipList()
	keys := make([]int, 100)
	for i := range keys {
		keys[i] = rand.Intn(1000)
		sl.insert(keys[i], keys[i])
	}
	sort.Ints(keys)

	// forEach must visit in ascending order
	prev := -1
	sl.forEach(func(k, _ int) {
		if k < prev {
			t.Fatalf("out of order: %d after %d", k, prev)
		}
		prev = k
	})
}

func BenchmarkSkipListInsert(b *testing.B) {
	sl := newSkipList()
	for i := range b.N {
		sl.insert(rand.Intn(1e6), i)
	}
}

func BenchmarkSkipListSearch(b *testing.B) {
	sl := newSkipList()
	for i := range 1000 {
		sl.insert(i, i)
	}
	b.ResetTimer()
	for range b.N {
		sl.search(rand.Intn(1000))
	}
}
