package tree

import (
	"errors"
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

// Constants
const (
	RED   = true
	BLACK = false
)

var (
	ErrNotFound = errors.New("key not found")
)

// NewRedBlackTree returns new RBT
func NewRedBlackTree[K constraints.Ordered, V comparable]() *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{
		pool: sync.Pool{
			New: func() interface{} {
				return &node[K, V]{}
			},
		},
	}
}

type node[K constraints.Ordered, V comparable] struct {
	key   K
	value V
	left  *node[K, V]
	right *node[K, V]
	color bool
}

func (n *node[K, V]) reuse() {
	n.left = nil
	n.right = nil
	n.color = BLACK
}

// RedBlackTree struct
type RedBlackTree[K constraints.Ordered, V comparable] struct {
	root *node[K, V]
	pool sync.Pool
}

func (t *RedBlackTree[K, V]) isEmpty() bool {
	return t.root == nil
}

// Put puts value under the key
func (t *RedBlackTree[K, V]) Put(key K, value V) {
	t.root = t.put(t.root, key, value)
	t.root.color = BLACK
}

func (t *RedBlackTree[K, V]) put(x *node[K, V], key K, value V) *node[K, V] {
	if x == nil {
		h := t.pool.Get().(*node[K, V])
		h.key = key
		h.value = value
		h.color = RED
		return h
	}
	if x.key > key {
		x.left = t.put(x.left, key, value)
	} else if x.key < key {
		x.right = t.put(x.right, key, value)
	} else {
		x.value = value
	}
	if t.isRed(x.right) && !t.isRed(x.left) {
		x = t.rotateLeft(x)
	}
	if t.isRed(x.left) && t.isRed(x.left.left) {
		x = t.rotateRight(x)
	}
	if t.isRed(x.left) && t.isRed(x.right) {
		t.flipColors(x)
	}
	return x
}

// Get returns value for the given key
func (t *RedBlackTree[K, V]) Get(key K) (data V, err error) {
	x := t.root
	for x != nil {
		if x.key > key {
			x = x.left
		} else if x.key < key {
			x = x.right
		} else {
			return x.value, nil
		}
	}
	return data, ErrNotFound
}

// Contains check that tree contains key
func (t *RedBlackTree[K, V]) Contains(key K) bool {
	_, err := t.Get(key)
	return !errors.Is(err, ErrNotFound)
}

// Delete removes by key
func (t *RedBlackTree[K, V]) Delete(key K) {
	if t.isEmpty() || !t.Contains(key) {
		return
	}
	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}
	t.root = t.deleteKey(t.root, key)
	if !t.isEmpty() {
		t.root.color = BLACK
	}
}

func (t *RedBlackTree[K, V]) deleteKey(x *node[K, V], key K) *node[K, V] {
	if x.key > key {
		if !t.isRed(x.left) && !t.isRed(x.left.left) {
			x = t.moveRedLeft(x)
		}
		x.left = t.deleteKey(x.left, key)
	} else {
		if t.isRed(x.left) {
			x = t.rotateRight(x)
		}
		if key == x.key && x.right == nil {
			x.reuse()
			t.pool.Put(x)
			return nil
		}
		if !t.isRed(x.right) && !t.isRed(x.right.left) {
			x = t.moveRedRight(x)
		}
		if key == x.key {
			h := t.min(x.right)
			x.key = h.key
			x.value = h.value
			x.right = t.deleteMin(x.right)
		} else {
			x.right = t.deleteKey(x.right, key)
		}
	}
	return t.balance(x)
}

// DeleteMin removes min key and associated value
func (t *RedBlackTree[K, V]) DeleteMin() {
	if t.isEmpty() {
		return
	}
	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}
	t.root = t.deleteMin(t.root)
	if !t.isEmpty() {
		t.root.color = BLACK
	}
}

func (t *RedBlackTree[K, V]) deleteMin(x *node[K, V]) *node[K, V] {
	if x.left == nil {
		x.reuse()
		t.pool.Put(x)
		return nil
	}
	if !t.isRed(x.left) && !t.isRed(x.left.left) {
		x = t.moveRedLeft(x)
	}
	x.left = t.deleteMin(x.left)
	return t.balance(x)
}

// DeleteMax removes max key and associated value
func (t *RedBlackTree[K, V]) DeleteMax() {
	if t.isEmpty() {
		return
	}
	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}
	t.root = t.deleteMax(t.root)
	if !t.isEmpty() {
		t.root.color = BLACK
	}
}

func (t *RedBlackTree[K, V]) deleteMax(x *node[K, V]) *node[K, V] {
	if t.isRed(x.left) {
		x = t.rotateRight(x)
	}
	if x.right == nil {
		x.reuse()
		t.pool.Put(x)
		return nil
	}
	if !t.isRed(x.right) && !t.isRed(x.right.left) {
		x = t.moveRedRight(x)
	}
	x.right = t.deleteMax(x.right)
	return t.balance(x)
}

// Floor returns value which key is nearest less or equal to key
func (t *RedBlackTree[K, V]) Floor(key K) (data V, err error) {
	node := t.floor(t.root, key)
	if node != nil {
		return node.value, nil
	}
	return data, ErrNotFound
}

func (t *RedBlackTree[K, V]) floor(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}
	if key == x.key {
		return x
	}
	if x.key > key {
		return t.floor(x.left, key)
	}
	temp := t.floor(x.right, key)
	if temp == nil {
		return x
	}
	return temp
}

// Ceil returns value which key is nearest greater or equal to key
func (t *RedBlackTree[K, V]) Ceil(key K) (data V, err error) {
	node := t.ceil(t.root, key)
	if node != nil {
		return node.value, nil
	}
	return data, ErrNotFound
}

func (t *RedBlackTree[K, V]) ceil(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}
	if key == x.key {
		return x
	}
	if key > x.key {
		return t.ceil(x.right, key)
	}
	temp := t.ceil(x.left, key)
	if temp == nil {
		return x
	}
	return temp
}

func (t *RedBlackTree[K, V]) isRed(x *node[K, V]) bool {
	if x == nil {
		return false
	}
	return x.color == RED
}

func (t *RedBlackTree[K, V]) rotateLeft(x *node[K, V]) *node[K, V] {
	h := x.right
	x.right = h.left
	h.left = x
	h.color = x.color
	x.color = RED
	return h
}

func (t *RedBlackTree[K, V]) rotateRight(x *node[K, V]) *node[K, V] {
	h := x.left
	x.left = h.right
	h.right = x
	h.color = x.color
	x.color = RED
	return h
}

func (t *RedBlackTree[K, V]) flipColors(x *node[K, V]) {
	x.color = !x.color
	x.left.color = !x.left.color
	x.right.color = !x.right.color
}

func (t *RedBlackTree[K, V]) moveRedLeft(x *node[K, V]) *node[K, V] {
	t.flipColors(x)
	if t.isRed(x.right.left) {
		x.right = t.rotateRight(x.right)
		x = t.rotateLeft(x)
		t.flipColors(x)
	}
	return x
}

func (t *RedBlackTree[K, V]) moveRedRight(x *node[K, V]) *node[K, V] {
	t.flipColors(x)
	if t.isRed(x.left.left) {
		x = t.rotateRight(x)
		t.flipColors(x)
	}
	return x
}

func (t *RedBlackTree[K, V]) balance(x *node[K, V]) *node[K, V] {
	if t.isRed(x.right) {
		x = t.rotateLeft(x)
	}
	if t.isRed(x.left) && t.isRed(x.left.left) {
		x = t.rotateRight(x)
	}
	if t.isRed(x.left) && t.isRed(x.right) {
		t.flipColors(x)
	}
	return x
}

func (t *RedBlackTree[K, V]) min(x *node[K, V]) *node[K, V] {
	if x == nil {
		return nil
	}
	for {
		if x.left == nil {
			return x
		}
		x = x.left
	}
}

func (t *RedBlackTree[K, V]) max(x *node[K, V]) *node[K, V] {
	if x == nil {
		return nil
	}
	for {
		if x.right == nil {
			return x
		}
		x = x.right
	}
}

// Min returns minimum key, its value
func (t *RedBlackTree[K, V]) Min() (key K, value V, err error) {
	node := t.min(t.root)
	if node != nil {
		return node.key, node.value, nil
	}
	return key, value, ErrNotFound
}

// Max returns minimum key, its value
func (t *RedBlackTree[K, V]) Max() (key K, value V, err error) {
	node := t.max(t.root)
	if node != nil {
		return node.key, node.value, nil
	}
	return key, value, ErrNotFound
}

func (t *RedBlackTree[K, V]) String() string {
	str := "RedBlackTree\n"
	if !t.isEmpty() {
		output(t.root, "", true, &str)
	}
	return str
}

func (n *node[K, V]) string() string {
	return fmt.Sprintf("%v", n.key)
}

func output[K constraints.Ordered, V comparable](node *node[K, V], prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.string() + "\n"
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.left, newPrefix, true, str)
	}
}
