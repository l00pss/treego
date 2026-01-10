package btree

import (
	"fmt"
	"strings"
)

// Ordered constraint for types that can be compared
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// BTree represents a generic B-tree
type BTree[K Ordered, V any] struct {
	root   *Node[K, V]
	degree int // minimum degree (t)
}

// Node represents a node in the B-tree
type Node[K Ordered, V any] struct {
	keys     []K
	values   []V
	children []*Node[K, V]
	isLeaf   bool
}

// KeyValue represents a key-value pair
type KeyValue[K Ordered, V any] struct {
	Key   K
	Value V
}

// NewBTree creates a new B-tree with the specified minimum degree
func NewBTree[K Ordered, V any](degree int) *BTree[K, V] {
	if degree < 2 {
		degree = 2 // minimum degree should be at least 2
	}
	return &BTree[K, V]{
		root:   newNode[K, V](true),
		degree: degree,
	}
}

// newNode creates a new node
func newNode[K Ordered, V any](isLeaf bool) *Node[K, V] {
	return &Node[K, V]{
		keys:     make([]K, 0),
		values:   make([]V, 0),
		children: make([]*Node[K, V], 0),
		isLeaf:   isLeaf,
	}
}

// Insert inserts a key-value pair into the B-tree
func (bt *BTree[K, V]) Insert(key K, value V) {
	root := bt.root
	if bt.isFull(root) {
		// Root is full, need to split
		newRoot := newNode[K, V](false)
		newRoot.children = append(newRoot.children, root)
		bt.splitChild(newRoot, 0)
		bt.root = newRoot
	}
	bt.insertNonFull(bt.root, key, value)
}

// Search searches for a key in the B-tree
func (bt *BTree[K, V]) Search(key K) (V, bool) {
	return bt.searchNode(bt.root, key)
}

// Delete removes a key from the B-tree
func (bt *BTree[K, V]) Delete(key K) bool {
	deleted := bt.deleteFromNode(bt.root, key)
	if len(bt.root.keys) == 0 && !bt.root.isLeaf {
		bt.root = bt.root.children[0]
	}
	return deleted
}

// InOrderTraversal performs in-order traversal of the B-tree
func (bt *BTree[K, V]) InOrderTraversal() []KeyValue[K, V] {
	var result []KeyValue[K, V]
	bt.inOrderTraverseNode(bt.root, &result)
	return result
}

// Height returns the height of the B-tree
func (bt *BTree[K, V]) Height() int {
	return bt.getHeight(bt.root)
}

// Size returns the total number of keys in the B-tree
func (bt *BTree[K, V]) Size() int {
	return bt.getSize(bt.root)
}

// IsEmpty checks if the B-tree is empty
func (bt *BTree[K, V]) IsEmpty() bool {
	return len(bt.root.keys) == 0
}

// isFull checks if a node is full
func (bt *BTree[K, V]) isFull(node *Node[K, V]) bool {
	return len(node.keys) == 2*bt.degree-1
}

// insertNonFull inserts into a non-full node
func (bt *BTree[K, V]) insertNonFull(node *Node[K, V], key K, value V) {
	i := len(node.keys) - 1

	if node.isLeaf {
		// Insert into leaf node
		node.keys = append(node.keys, key)
		node.values = append(node.values, value)

		// Shift elements to maintain sorted order
		for i >= 0 && node.keys[i] > key {
			node.keys[i+1] = node.keys[i]
			node.values[i+1] = node.values[i]
			i--
		}
		node.keys[i+1] = key
		node.values[i+1] = value
	} else {
		// Find child to recurse on
		for i >= 0 && node.keys[i] > key {
			i--
		}
		i++

		if bt.isFull(node.children[i]) {
			bt.splitChild(node, i)
			if node.keys[i] < key {
				i++
			}
		}
		bt.insertNonFull(node.children[i], key, value)
	}
}

// splitChild splits a full child node
func (bt *BTree[K, V]) splitChild(parent *Node[K, V], index int) {
	fullChild := parent.children[index]
	newChild := newNode[K, V](fullChild.isLeaf)

	mid := bt.degree - 1

	// Store middle key and value before truncating
	midKey := fullChild.keys[mid]
	midValue := fullChild.values[mid]

	// Move half of keys and values to new child
	newChild.keys = append(newChild.keys, fullChild.keys[bt.degree:]...)
	newChild.values = append(newChild.values, fullChild.values[bt.degree:]...)
	fullChild.keys = fullChild.keys[:mid]
	fullChild.values = fullChild.values[:mid]

	// Move children if not leaf
	if !fullChild.isLeaf {
		newChild.children = append(newChild.children, fullChild.children[bt.degree:]...)
		fullChild.children = fullChild.children[:bt.degree]
	}

	// Insert middle key into parent
	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = newChild

	parent.keys = append(parent.keys, midKey)
	parent.values = append(parent.values, midValue)

	// Shift keys and values in parent
	for i := len(parent.keys) - 1; i > index; i-- {
		parent.keys[i] = parent.keys[i-1]
		parent.values[i] = parent.values[i-1]
	}
	parent.keys[index] = midKey
	parent.values[index] = midValue
}

// searchNode searches for a key in a node
func (bt *BTree[K, V]) searchNode(node *Node[K, V], key K) (V, bool) {
	var zero V
	i := 0

	// Find the first key greater than or equal to key
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	// If found
	if i < len(node.keys) && key == node.keys[i] {
		return node.values[i], true
	}

	// If leaf node, key doesn't exist
	if node.isLeaf {
		return zero, false
	}

	// Recurse on appropriate child
	return bt.searchNode(node.children[i], key)
}

// deleteFromNode deletes a key from a node
func (bt *BTree[K, V]) deleteFromNode(node *Node[K, V], key K) bool {
	i := 0

	// Find the index of the key or the child that should contain the key
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	if i < len(node.keys) && key == node.keys[i] {
		// Key found in this node
		if node.isLeaf {
			// Delete from leaf
			copy(node.keys[i:], node.keys[i+1:])
			copy(node.values[i:], node.values[i+1:])
			node.keys = node.keys[:len(node.keys)-1]
			node.values = node.values[:len(node.values)-1]
			return true
		} else {
			// Delete from internal node
			return bt.deleteFromInternalNode(node, i)
		}
	} else if !node.isLeaf {
		// Key not found in this node, recurse on child
		if len(node.children[i].keys) >= bt.degree {
			return bt.deleteFromNode(node.children[i], key)
		} else {
			// Child has minimum keys, need to handle underflow
			bt.handleChildUnderflow(node, i)
			return bt.deleteFromNode(node, key)
		}
	}

	return false // Key not found
}

// deleteFromInternalNode deletes a key from an internal node
func (bt *BTree[K, V]) deleteFromInternalNode(node *Node[K, V], index int) bool {
	key := node.keys[index]

	// Case 1: Left child has at least t keys
	if len(node.children[index].keys) >= bt.degree {
		pred := bt.getPredecessor(node, index)
		node.keys[index] = pred.Key
		node.values[index] = pred.Value
		return bt.deleteFromNode(node.children[index], pred.Key)
	}

	// Case 2: Right child has at least t keys
	if len(node.children[index+1].keys) >= bt.degree {
		succ := bt.getSuccessor(node, index)
		node.keys[index] = succ.Key
		node.values[index] = succ.Value
		return bt.deleteFromNode(node.children[index+1], succ.Key)
	}

	// Case 3: Both children have t-1 keys, merge
	bt.mergeChildren(node, index)
	return bt.deleteFromNode(node.children[index], key)
}

// getPredecessor gets the predecessor of a key
func (bt *BTree[K, V]) getPredecessor(node *Node[K, V], index int) KeyValue[K, V] {
	curr := node.children[index]
	for !curr.isLeaf {
		curr = curr.children[len(curr.children)-1]
	}
	lastIndex := len(curr.keys) - 1
	return KeyValue[K, V]{Key: curr.keys[lastIndex], Value: curr.values[lastIndex]}
}

// getSuccessor gets the successor of a key
func (bt *BTree[K, V]) getSuccessor(node *Node[K, V], index int) KeyValue[K, V] {
	curr := node.children[index+1]
	for !curr.isLeaf {
		curr = curr.children[0]
	}
	return KeyValue[K, V]{Key: curr.keys[0], Value: curr.values[0]}
}

// handleChildUnderflow handles child underflow
func (bt *BTree[K, V]) handleChildUnderflow(node *Node[K, V], index int) {
	// Try borrowing from left sibling
	if index > 0 && len(node.children[index-1].keys) >= bt.degree {
		bt.borrowFromLeftSibling(node, index)
		return
	}

	// Try borrowing from right sibling
	if index < len(node.children)-1 && len(node.children[index+1].keys) >= bt.degree {
		bt.borrowFromRightSibling(node, index)
		return
	}

	// Merge with sibling
	if index > 0 {
		bt.mergeChildren(node, index-1)
	} else {
		bt.mergeChildren(node, index)
	}
}

// borrowFromLeftSibling borrows a key from left sibling
func (bt *BTree[K, V]) borrowFromLeftSibling(parent *Node[K, V], index int) {
	child := parent.children[index]
	sibling := parent.children[index-1]

	// Move parent key to child
	child.keys = append([]K{parent.keys[index-1]}, child.keys...)
	child.values = append([]V{parent.values[index-1]}, child.values...)

	// Move sibling's last key to parent
	lastIndex := len(sibling.keys) - 1
	parent.keys[index-1] = sibling.keys[lastIndex]
	parent.values[index-1] = sibling.values[lastIndex]
	sibling.keys = sibling.keys[:lastIndex]
	sibling.values = sibling.values[:lastIndex]

	// Move child if not leaf
	if !child.isLeaf {
		child.children = append([]*Node[K, V]{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}
}

// borrowFromRightSibling borrows a key from right sibling
func (bt *BTree[K, V]) borrowFromRightSibling(parent *Node[K, V], index int) {
	child := parent.children[index]
	sibling := parent.children[index+1]

	// Move parent key to child
	child.keys = append(child.keys, parent.keys[index])
	child.values = append(child.values, parent.values[index])

	// Move sibling's first key to parent
	parent.keys[index] = sibling.keys[0]
	parent.values[index] = sibling.values[0]
	sibling.keys = sibling.keys[1:]
	sibling.values = sibling.values[1:]

	// Move child if not leaf
	if !child.isLeaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}
}

// mergeChildren merges two children
func (bt *BTree[K, V]) mergeChildren(parent *Node[K, V], index int) {
	child := parent.children[index]
	sibling := parent.children[index+1]

	// Move parent key to child
	child.keys = append(child.keys, parent.keys[index])
	child.values = append(child.values, parent.values[index])

	// Move all keys and values from sibling to child
	child.keys = append(child.keys, sibling.keys...)
	child.values = append(child.values, sibling.values...)

	// Move children if not leaf
	if !child.isLeaf {
		child.children = append(child.children, sibling.children...)
	}

	// Remove key from parent
	copy(parent.keys[index:], parent.keys[index+1:])
	copy(parent.values[index:], parent.values[index+1:])
	parent.keys = parent.keys[:len(parent.keys)-1]
	parent.values = parent.values[:len(parent.values)-1]

	// Remove child pointer from parent
	copy(parent.children[index+1:], parent.children[index+2:])
	parent.children = parent.children[:len(parent.children)-1]
}

// inOrderTraverseNode performs in-order traversal of a node
func (bt *BTree[K, V]) inOrderTraverseNode(node *Node[K, V], result *[]KeyValue[K, V]) {
	i := 0
	for i < len(node.keys) {
		if !node.isLeaf {
			bt.inOrderTraverseNode(node.children[i], result)
		}
		*result = append(*result, KeyValue[K, V]{Key: node.keys[i], Value: node.values[i]})
		i++
	}

	if !node.isLeaf {
		bt.inOrderTraverseNode(node.children[i], result)
	}
}

// getHeight calculates the height of a node
func (bt *BTree[K, V]) getHeight(node *Node[K, V]) int {
	if node.isLeaf {
		return 0
	}
	return 1 + bt.getHeight(node.children[0])
}

// getSize calculates the total number of keys in a subtree
func (bt *BTree[K, V]) getSize(node *Node[K, V]) int {
	size := len(node.keys)
	if !node.isLeaf {
		for _, child := range node.children {
			size += bt.getSize(child)
		}
	}
	return size
}

// String returns a string representation of the B-tree
func (bt *BTree[K, V]) String() string {
	return bt.nodeString(bt.root, 0)
}

// nodeString returns a string representation of a node
func (bt *BTree[K, V]) nodeString(node *Node[K, V], level int) string {
	indent := strings.Repeat("  ", level)
	result := fmt.Sprintf("%sNode(leaf=%v): %v\n", indent, node.isLeaf, node.keys)

	if !node.isLeaf {
		for _, child := range node.children {
			result += bt.nodeString(child, level+1)
		}
	}

	return result
}
