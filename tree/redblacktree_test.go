package tree

import (
	"errors"
	"math/rand"
	"testing"
)

const size = 100_000

var (
	keys   []int
	values []int
)

func init() {
	keys = make([]int, size, size)
	values = make([]int, size, size)
	for i := 0; i < size; i++ {
		keys[i] = rand.Int()
		values[i] = rand.Int()
	}
}

func TestNodeReuse(t *testing.T) {
	n := node[int, int]{left: &node[int, int]{}, right: &node[int, int]{}, color: RED}
	n.reuse()
	if n.left != nil {
		t.Errorf("Got %v expected nil", n.left)
	}
	if n.right != nil {
		t.Errorf("Got %v expected nil", n.right)
	}
	if n.color != BLACK {
		t.Errorf("Got %v expected %v", n.color, BLACK)
	}
}

func TestTreeIsEmpty(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
	tree.Put(0, 0)
	if tree.isEmpty() {
		t.Errorf("Got %v expected false", tree.isEmpty())
	}
}

func TestTreePut(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	for i, v := range keys {
		tree.Put(v, values[i])
		tree.Put(v, values[i])
	}
	for _, v := range keys {
		cont := tree.Contains(v)
		if !cont {
			t.Errorf("Got %v expected true", cont)
		}
	}
}

func TestTreeGet(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	for i, v := range keys {
		val, err := tree.Get(i)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Got value: %v and error: %v expected error: %v", val, err, ErrNotFound)
		}
		val, err = tree.Get(v)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Got value: %v and error: %v expected error: %v", val, err, ErrNotFound)
		}
	}
	for i, v := range keys {
		tree.Put(v, values[i])
		val, err := tree.Get(v)
		if errors.Is(err, ErrNotFound) || val != values[i] {
			t.Errorf("Got value: %v and error: %v expected value: %v and error: %v", val, err, values[i], ErrNotFound)
		}
	}
}

func TestTreeDelete(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	tree.Delete(100)
	for i, v := range keys {
		tree.Put(v, values[i])
	}
	for _, v := range keys {
		if tree.isEmpty() {
			t.Errorf("Got %v expected false", tree.isEmpty())
		}
		tree.Delete(v)
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
	for i, v := range keys {
		tree.Put(v, values[i])
		if tree.isEmpty() {
			t.Errorf("Got %v expected false", tree.isEmpty())
		}
		tree.Delete(v)
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
}

func TestTreeDeleteMin(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	tree.DeleteMin()
	for i, v := range keys {
		tree.Put(v, values[i])
	}
	for _ = range keys {
		if tree.isEmpty() {
			t.Errorf("Got %v expected false", tree.isEmpty())
		}
		tree.DeleteMin()
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
	for i, v := range keys {
		tree.Put(v, values[i])
		tree.DeleteMin()
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
}

func TestTreeDeleteMax(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	tree.DeleteMax()
	for i, v := range keys {
		tree.Put(v, values[i])
	}
	for _ = range keys {
		if tree.isEmpty() {
			t.Errorf("Got %v expected false", tree.isEmpty())
		}
		tree.DeleteMax()
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
	for i, v := range keys {
		tree.Put(v, values[i])
		tree.DeleteMax()
	}
	if !tree.isEmpty() {
		t.Errorf("Got %v expected true", tree.isEmpty())
	}
}

func TestTreeFloor(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	val, err := tree.Floor(100)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Got value: %v and error: %v expected error: %v", val, err, ErrNotFound)
	}
	for i := 0; i <= 100; i++ {
		tree.Put(i, i)
		val, err = tree.Floor(100)
		if errors.Is(err, ErrNotFound) || val != i {
			t.Errorf("Got value: %v and error: %v expected value: %v and error: %v", val, err, i, ErrNotFound)
		}
	}
	for i := 100; i > 0; i-- {
		tree.Delete(i)
		val, err = tree.Floor(100)
		if errors.Is(err, ErrNotFound) || val != i-1 {
			t.Errorf("Got value: %v and error: %v expected value: %v and error: %v", val, err, i, ErrNotFound)
		}
	}
}

func TestTreeCeil(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	val, err := tree.Ceil(0)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Got value: %v and error: %v expected error: %v", val, err, ErrNotFound)
	}
	for i := 100; i > 0; i-- {
		tree.Put(i, i)
		val, err = tree.Ceil(0)
		if errors.Is(err, ErrNotFound) || val != i {
			t.Errorf("Got value: %v and error: %v expected value: %v and error: %v", val, err, i, ErrNotFound)
		}
	}
	for i := 1; i < 100; i++ {
		tree.Delete(i)
		val, err = tree.Ceil(0)
		if errors.Is(err, ErrNotFound) || val != i+1 {
			t.Errorf("Got value: %v and error: %v expected value: %v and error: %v", val, err, i, ErrNotFound)
		}
	}
}

func TestTreeMin(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	key, val, err := tree.Min()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Got key: %v and value: %v and error: %v expected error: %v", key, val, err, ErrNotFound)
	}
	for i := 100; i > 0; i-- {
		tree.Put(i, i)
		key, val, err = tree.Min()
		if errors.Is(err, ErrNotFound) || val != i || key != i {
			t.Errorf("Got key: %v and value: %v and error: %v expected key: %v and value: %v and error: %v", key, val, err, i, i, ErrNotFound)
		}
	}
	tree.Put(0, 100)
	for i := 1; i < 100; i++ {
		tree.Delete(i)
		key, val, err = tree.Min()
		if errors.Is(err, ErrNotFound) || val != 100 || key != 0 {
			t.Errorf("Got key: %v and value: %v and error: %v expected key: %v and value: %v and error: %v", key, val, err, 0, 100, ErrNotFound)
		}
	}
}

func TestTreeMax(t *testing.T) {
	tree := NewRedBlackTree[int, int]()
	key, val, err := tree.Max()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Got key: %v and value: %v and error: %v expected error: %v", key, val, err, ErrNotFound)
	}
	for i := 0; i < 100; i++ {
		tree.Put(i, i)
		key, val, err = tree.Max()
		if errors.Is(err, ErrNotFound) || val != i || key != i {
			t.Errorf("Got key: %v and value: %v and error: %v expected key: %v and value: %v and error: %v", key, val, err, i, i, ErrNotFound)
		}
	}
	tree.Put(100, 0)
	for i := 0; i < 100; i++ {
		tree.Delete(i)
		key, val, err = tree.Max()
		if errors.Is(err, ErrNotFound) || val != 0 || key != 100 {
			t.Errorf("Got key: %v and value: %v and error: %v expected key: %v and value: %v and error: %v", key, val, err, 0, 100, ErrNotFound)
		}
	}
}

func BenchmarkPut(b *testing.B) {
	tr := NewRedBlackTree[int, int]()
	for i := 0; i < b.N; i++ {
		tr.Put(keys[i%size], values[i%size])
	}
}
