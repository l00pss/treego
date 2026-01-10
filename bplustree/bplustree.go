package bplustree

import "cmp"

type Entry[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

type node[K cmp.Ordered, V any] struct {
	isLeaf   bool
	keys     []K
	children []*node[K, V]
	entries  []Entry[K, V]
	next     *node[K, V]
	parent   *node[K, V]
}

type BPlusTree[K cmp.Ordered, V any] struct {
	root   *node[K, V]
	degree int
}

func New[K cmp.Ordered, V any](degree int) *BPlusTree[K, V] {
	if degree < 2 {
		degree = 2
	}
	return &BPlusTree[K, V]{degree: degree}
}

func (t *BPlusTree[K, V]) Search(key K) (V, bool) {
	if t.root == nil {
		var zero V
		return zero, false
	}
	leaf := t.findLeaf(key)
	for _, e := range leaf.entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	var zero V
	return zero, false
}

func (t *BPlusTree[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = &node[K, V]{isLeaf: true}
		t.root.entries = []Entry[K, V]{{Key: key, Value: value}}
		return
	}

	leaf := t.findLeaf(key)

	for i, e := range leaf.entries {
		if e.Key == key {
			leaf.entries[i].Value = value
			return
		}
	}

	t.insertIntoLeaf(leaf, key, value)

	if len(leaf.entries) > t.maxLeafEntries() {
		t.splitLeaf(leaf)
	}
}

func (t *BPlusTree[K, V]) Delete(key K) bool {
	if t.root == nil {
		return false
	}

	leaf := t.findLeaf(key)
	idx := -1
	for i, e := range leaf.entries {
		if e.Key == key {
			idx = i
			break
		}
	}

	if idx == -1 {
		return false
	}

	leaf.entries = append(leaf.entries[:idx], leaf.entries[idx+1:]...)

	if leaf == t.root {
		if len(leaf.entries) == 0 {
			t.root = nil
		}
		return true
	}

	minEntries := t.minLeafEntries()
	if len(leaf.entries) < minEntries {
		t.rebalanceLeaf(leaf)
	}

	return true
}

func (t *BPlusTree[K, V]) Range(start, end K) []Entry[K, V] {
	if t.root == nil {
		return nil
	}

	var result []Entry[K, V]
	leaf := t.findLeaf(start)

	for leaf != nil {
		for _, e := range leaf.entries {
			if e.Key >= start && e.Key <= end {
				result = append(result, e)
			} else if e.Key > end {
				return result
			}
		}
		leaf = leaf.next
	}
	return result
}

func (t *BPlusTree[K, V]) All() []Entry[K, V] {
	if t.root == nil {
		return nil
	}

	var result []Entry[K, V]
	leaf := t.firstLeaf()
	for leaf != nil {
		result = append(result, leaf.entries...)
		leaf = leaf.next
	}
	return result
}

func (t *BPlusTree[K, V]) Len() int {
	if t.root == nil {
		return 0
	}
	count := 0
	leaf := t.firstLeaf()
	for leaf != nil {
		count += len(leaf.entries)
		leaf = leaf.next
	}
	return count
}

func (t *BPlusTree[K, V]) findLeaf(key K) *node[K, V] {
	n := t.root
	for !n.isLeaf {
		i := 0
		for i < len(n.keys) && key >= n.keys[i] {
			i++
		}
		n = n.children[i]
	}
	return n
}

func (t *BPlusTree[K, V]) firstLeaf() *node[K, V] {
	if t.root == nil {
		return nil
	}
	n := t.root
	for !n.isLeaf {
		n = n.children[0]
	}
	return n
}

func (t *BPlusTree[K, V]) insertIntoLeaf(leaf *node[K, V], key K, value V) {
	entry := Entry[K, V]{Key: key, Value: value}
	i := 0
	for i < len(leaf.entries) && leaf.entries[i].Key < key {
		i++
	}
	leaf.entries = append(leaf.entries[:i], append([]Entry[K, V]{entry}, leaf.entries[i:]...)...)
}

func (t *BPlusTree[K, V]) splitLeaf(leaf *node[K, V]) {
	mid := len(leaf.entries) / 2

	newLeaf := &node[K, V]{
		isLeaf:  true,
		entries: make([]Entry[K, V], len(leaf.entries[mid:])),
		next:    leaf.next,
		parent:  leaf.parent,
	}
	copy(newLeaf.entries, leaf.entries[mid:])
	leaf.entries = leaf.entries[:mid]
	leaf.next = newLeaf

	t.insertIntoParent(leaf, newLeaf.entries[0].Key, newLeaf)
}

func (t *BPlusTree[K, V]) splitInternal(n *node[K, V]) {
	mid := len(n.keys) / 2
	promoteKey := n.keys[mid]

	newNode := &node[K, V]{
		isLeaf:   false,
		keys:     make([]K, len(n.keys[mid+1:])),
		children: make([]*node[K, V], len(n.children[mid+1:])),
		parent:   n.parent,
	}
	copy(newNode.keys, n.keys[mid+1:])
	copy(newNode.children, n.children[mid+1:])

	for _, child := range newNode.children {
		child.parent = newNode
	}

	n.keys = n.keys[:mid]
	n.children = n.children[:mid+1]

	t.insertIntoParent(n, promoteKey, newNode)
}

func (t *BPlusTree[K, V]) insertIntoParent(left *node[K, V], key K, right *node[K, V]) {
	if left.parent == nil {
		newRoot := &node[K, V]{
			isLeaf:   false,
			keys:     []K{key},
			children: []*node[K, V]{left, right},
		}
		t.root = newRoot
		left.parent = newRoot
		right.parent = newRoot
		return
	}

	parent := left.parent
	right.parent = parent

	i := 0
	for i < len(parent.children) && parent.children[i] != left {
		i++
	}

	parent.keys = append(parent.keys[:i], append([]K{key}, parent.keys[i:]...)...)
	parent.children = append(parent.children[:i+1], append([]*node[K, V]{right}, parent.children[i+1:]...)...)

	if len(parent.keys) > t.maxInternalKeys() {
		t.splitInternal(parent)
	}
}

func (t *BPlusTree[K, V]) rebalanceLeaf(leaf *node[K, V]) {
	parent := leaf.parent
	if parent == nil {
		return
	}

	idx := 0
	for idx < len(parent.children) && parent.children[idx] != leaf {
		idx++
	}

	if idx > 0 {
		leftSibling := parent.children[idx-1]
		if len(leftSibling.entries) > t.minLeafEntries() {
			borrowed := leftSibling.entries[len(leftSibling.entries)-1]
			leftSibling.entries = leftSibling.entries[:len(leftSibling.entries)-1]
			leaf.entries = append([]Entry[K, V]{borrowed}, leaf.entries...)
			parent.keys[idx-1] = leaf.entries[0].Key
			return
		}
	}

	if idx < len(parent.children)-1 {
		rightSibling := parent.children[idx+1]
		if len(rightSibling.entries) > t.minLeafEntries() {
			borrowed := rightSibling.entries[0]
			rightSibling.entries = rightSibling.entries[1:]
			leaf.entries = append(leaf.entries, borrowed)
			parent.keys[idx] = rightSibling.entries[0].Key
			return
		}
	}

	if idx > 0 {
		leftSibling := parent.children[idx-1]
		leftSibling.entries = append(leftSibling.entries, leaf.entries...)
		leftSibling.next = leaf.next
		t.deleteFromParent(parent, idx-1, leaf)
	} else if idx < len(parent.children)-1 {
		rightSibling := parent.children[idx+1]
		leaf.entries = append(leaf.entries, rightSibling.entries...)
		leaf.next = rightSibling.next
		t.deleteFromParent(parent, idx, rightSibling)
	}
}

func (t *BPlusTree[K, V]) deleteFromParent(parent *node[K, V], keyIdx int, child *node[K, V]) {
	parent.keys = append(parent.keys[:keyIdx], parent.keys[keyIdx+1:]...)

	childIdx := keyIdx + 1
	parent.children = append(parent.children[:childIdx], parent.children[childIdx+1:]...)

	if parent == t.root && len(parent.keys) == 0 {
		t.root = parent.children[0]
		t.root.parent = nil
		return
	}

	if parent.parent != nil && len(parent.keys) < t.minInternalKeys() {
		t.rebalanceInternal(parent)
	}
}

func (t *BPlusTree[K, V]) rebalanceInternal(n *node[K, V]) {
	parent := n.parent
	if parent == nil {
		return
	}

	idx := 0
	for idx < len(parent.children) && parent.children[idx] != n {
		idx++
	}

	if idx > 0 {
		leftSibling := parent.children[idx-1]
		if len(leftSibling.keys) > t.minInternalKeys() {
			borrowedKey := leftSibling.keys[len(leftSibling.keys)-1]
			borrowedChild := leftSibling.children[len(leftSibling.children)-1]

			leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
			leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]

			n.keys = append([]K{parent.keys[idx-1]}, n.keys...)
			n.children = append([]*node[K, V]{borrowedChild}, n.children...)
			borrowedChild.parent = n

			parent.keys[idx-1] = borrowedKey
			return
		}
	}

	if idx < len(parent.children)-1 {
		rightSibling := parent.children[idx+1]
		if len(rightSibling.keys) > t.minInternalKeys() {
			borrowedKey := rightSibling.keys[0]
			borrowedChild := rightSibling.children[0]

			rightSibling.keys = rightSibling.keys[1:]
			rightSibling.children = rightSibling.children[1:]

			n.keys = append(n.keys, parent.keys[idx])
			n.children = append(n.children, borrowedChild)
			borrowedChild.parent = n

			parent.keys[idx] = borrowedKey
			return
		}
	}

	if idx > 0 {
		leftSibling := parent.children[idx-1]
		leftSibling.keys = append(leftSibling.keys, parent.keys[idx-1])
		leftSibling.keys = append(leftSibling.keys, n.keys...)
		leftSibling.children = append(leftSibling.children, n.children...)
		for _, child := range n.children {
			child.parent = leftSibling
		}
		t.deleteFromParent(parent, idx-1, n)
	} else if idx < len(parent.children)-1 {
		rightSibling := parent.children[idx+1]
		n.keys = append(n.keys, parent.keys[idx])
		n.keys = append(n.keys, rightSibling.keys...)
		n.children = append(n.children, rightSibling.children...)
		for _, child := range rightSibling.children {
			child.parent = n
		}
		t.deleteFromParent(parent, idx, rightSibling)
	}
}

func (t *BPlusTree[K, V]) maxLeafEntries() int {
	return t.degree*2 - 1
}

func (t *BPlusTree[K, V]) minLeafEntries() int {
	return t.degree - 1
}

func (t *BPlusTree[K, V]) maxInternalKeys() int {
	return t.degree*2 - 1
}

func (t *BPlusTree[K, V]) minInternalKeys() int {
	return t.degree - 1
}
