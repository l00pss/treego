package btree

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBTreeBasicOperations(t *testing.T) {
	btree := NewBTree[int, string](3)

	btree.Insert(10, "ten")
	btree.Insert(20, "twenty")
	btree.Insert(5, "five")
	btree.Insert(6, "six")
	btree.Insert(12, "twelve")
	btree.Insert(30, "thirty")
	btree.Insert(7, "seven")
	btree.Insert(17, "seventeen")

	if val, found := btree.Search(10); !found || val != "ten" {
		t.Errorf("Expected to find '10' -> 'ten', got found=%v, val=%v", found, val)
	}

	if val, found := btree.Search(25); found {
		t.Errorf("Expected not to find '25', but found val=%v", val)
	}

	if size := btree.Size(); size != 8 {
		t.Errorf("Expected size 8, got %d", size)
	}

	items := btree.InOrderTraversal()
	expectedKeys := []int{5, 6, 7, 10, 12, 17, 20, 30}

	if len(items) != len(expectedKeys) {
		t.Errorf("Expected %d items, got %d", len(expectedKeys), len(items))
	}

	for i, item := range items {
		if item.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, item.Key)
		}
	}
}

func TestBTreeDeletion(t *testing.T) {
	btree := NewBTree[int, string](3)

	keys := []int{10, 20, 5, 6, 12, 30, 7, 17, 25, 40, 50}
	for _, key := range keys {
		btree.Insert(key, "value")
	}

	if !btree.Delete(6) {
		t.Error("Expected to delete key 6")
	}

	if btree.Delete(100) {
		t.Error("Expected not to delete non-existent key 100")
	}

	if _, found := btree.Search(6); found {
		t.Error("Key 6 should not exist after deletion")
	}

	if size := btree.Size(); size != 10 {
		t.Errorf("Expected size 10 after deletion, got %d", size)
	}
}

func TestBTreeWithStrings(t *testing.T) {
	btree := NewBTree[string, int](2)

	btree.Insert("apple", 1)
	btree.Insert("banana", 2)
	btree.Insert("cherry", 3)
	btree.Insert("date", 4)

	if val, found := btree.Search("banana"); !found || val != 2 {
		t.Errorf("Expected to find 'banana' -> 2, got found=%v, val=%v", found, val)
	}

	items := btree.InOrderTraversal()
	expectedKeys := []string{"apple", "banana", "cherry", "date"}

	for i, item := range items {
		if item.Key != expectedKeys[i] {
			t.Errorf("Expected key %s at position %d, got %s", expectedKeys[i], i, item.Key)
		}
	}
}

func TestBTreeEmpty(t *testing.T) {
	btree := NewBTree[int, string](3)

	if !btree.IsEmpty() {
		t.Error("New B-tree should be empty")
	}

	if size := btree.Size(); size != 0 {
		t.Errorf("Empty B-tree should have size 0, got %d", size)
	}

	if height := btree.Height(); height != 0 {
		t.Errorf("Empty B-tree should have height 0, got %d", height)
	}

	btree.Insert(1, "one")

	if btree.IsEmpty() {
		t.Error("B-tree with one element should not be empty")
	}
}

func TestDeleteSingleElement(t *testing.T) {
	btree := NewBTree[int, int](3)
	btree.Insert(1, 10)

	if !btree.Delete(1) {
		t.Error("Delete single element should return true")
	}

	if !btree.IsEmpty() {
		t.Error("Tree should be empty after deleting single element")
	}
}

func TestDeleteAllElements(t *testing.T) {
	btree := NewBTree[int, int](3)
	n := 100

	for i := 1; i <= n; i++ {
		btree.Insert(i, i*10)
	}

	for i := 1; i <= n; i++ {
		if !btree.Delete(i) {
			t.Errorf("Delete(%d) should return true", i)
		}
	}

	if !btree.IsEmpty() {
		t.Error("Tree should be empty after deleting all elements")
	}
}

func TestDeleteReverseOrder(t *testing.T) {
	btree := NewBTree[int, int](3)
	n := 50

	for i := 1; i <= n; i++ {
		btree.Insert(i, i)
	}

	for i := n; i >= 1; i-- {
		if !btree.Delete(i) {
			t.Errorf("Delete(%d) should return true", i)
		}
		if btree.Size() != i-1 {
			t.Errorf("After deleting %d, expected size=%d, got=%d", i, i-1, btree.Size())
		}
	}
}

func TestDeleteRandomOrder(t *testing.T) {
	btree := NewBTree[int, int](3)
	n := 100

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i
		btree.Insert(i, i*10)
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for i, k := range keys {
		if !btree.Delete(k) {
			t.Errorf("Delete(%d) should return true", k)
		}
		if btree.Size() != n-i-1 {
			t.Errorf("After deleting %d keys, expected size=%d, got=%d", i+1, n-i-1, btree.Size())
		}
	}
}

func TestDeleteAndReinsert(t *testing.T) {
	btree := NewBTree[int, int](3)

	for i := 1; i <= 20; i++ {
		btree.Insert(i, i*10)
	}

	for i := 1; i <= 10; i++ {
		btree.Delete(i)
	}

	for i := 1; i <= 10; i++ {
		btree.Insert(i, i*100)
	}

	for i := 1; i <= 10; i++ {
		value, found := btree.Search(i)
		if !found || value != i*100 {
			t.Errorf("Reinserted key %d: expected %d, got %d", i, i*100, value)
		}
	}

	if btree.Size() != 20 {
		t.Errorf("Expected size=20, got=%d", btree.Size())
	}
}

// === Degree Boundary Tests ===

func TestMinimumDegree(t *testing.T) {
	btree := NewBTree[int, int](2)

	for i := 1; i <= 100; i++ {
		btree.Insert(i, i)
	}

	if btree.Size() != 100 {
		t.Errorf("Expected size=100, got=%d", btree.Size())
	}

	for i := 1; i <= 100; i++ {
		value, found := btree.Search(i)
		if !found || value != i {
			t.Errorf("Search(%d) failed", i)
		}
	}
}

func TestLargeDegree(t *testing.T) {
	btree := NewBTree[int, int](50)

	for i := 1; i <= 1000; i++ {
		btree.Insert(i, i)
	}

	if btree.Size() != 1000 {
		t.Errorf("Expected size=1000, got=%d", btree.Size())
	}

	for i := 1; i <= 1000; i++ {
		if _, found := btree.Search(i); !found {
			t.Errorf("Search(%d) should find key", i)
		}
	}
}

func TestDegreeOne(t *testing.T) {
	btree := NewBTree[int, int](1)

	for i := 1; i <= 10; i++ {
		btree.Insert(i, i)
	}

	if btree.Size() != 10 {
		t.Errorf("Expected size=10, got=%d", btree.Size())
	}
}

// === Tree Structure Validation ===

func (bt *BTree[K, V]) validate() error {
	if bt.root == nil {
		return nil
	}
	return bt.validateNode(bt.root, true)
}

func (bt *BTree[K, V]) validateNode(node *Node[K, V], isRoot bool) error {
	maxKeys := 2*bt.degree - 1
	minKeys := bt.degree - 1

	if len(node.keys) > maxKeys {
		return fmt.Errorf("node has too many keys: %d > %d", len(node.keys), maxKeys)
	}

	if !isRoot && len(node.keys) < minKeys {
		return fmt.Errorf("non-root node has too few keys: %d < %d", len(node.keys), minKeys)
	}

	if len(node.keys) != len(node.values) {
		return fmt.Errorf("keys and values count mismatch: %d keys, %d values", len(node.keys), len(node.values))
	}

	for i := 1; i < len(node.keys); i++ {
		if node.keys[i-1] >= node.keys[i] {
			return fmt.Errorf("keys not sorted at index %d", i)
		}
	}

	if !node.isLeaf {
		if len(node.children) != len(node.keys)+1 {
			return fmt.Errorf("children count mismatch: %d children, %d keys", len(node.children), len(node.keys))
		}

		for i, child := range node.children {
			if err := bt.validateNode(child, false); err != nil {
				return fmt.Errorf("child %d invalid: %v", i, err)
			}
		}
	}

	return nil
}

func TestTreeStructureAfterInserts(t *testing.T) {
	btree := NewBTree[int, int](3)

	for i := 1; i <= 100; i++ {
		btree.Insert(i, i)
		if err := btree.validate(); err != nil {
			t.Errorf("Invalid tree after inserting %d: %v", i, err)
		}
	}
}

func TestTreeStructureAfterDeletes(t *testing.T) {
	btree := NewBTree[int, int](3)

	for i := 1; i <= 50; i++ {
		btree.Insert(i, i)
	}

	for i := 1; i <= 50; i++ {
		btree.Delete(i)
		if err := btree.validate(); err != nil {
			t.Errorf("Invalid tree after deleting %d: %v", i, err)
		}
	}
}

func TestTreeStructureRandomOps(t *testing.T) {
	btree := NewBTree[int, int](3)
	keys := make(map[int]bool)

	for i := 0; i < 500; i++ {
		op := rand.Intn(3)
		key := rand.Intn(100)

		switch op {
		case 0, 1:
			if !keys[key] {
				btree.Insert(key, key*10)
				keys[key] = true
			}
		case 2:
			btree.Delete(key)
			delete(keys, key)
		}

		if err := btree.validate(); err != nil {
			t.Errorf("Invalid tree at iteration %d (op=%d, key=%d): %v", i, op, key, err)
		}
	}
}

func TestInOrderTraversalSorted(t *testing.T) {
	btree := NewBTree[int, int](3)

	for i := 0; i < 100; i++ {
		btree.Insert(i, i)
	}

	items := btree.InOrderTraversal()
	for i := 1; i < len(items); i++ {
		if items[i-1].Key >= items[i].Key {
			t.Error("InOrderTraversal not sorted")
			break
		}
	}
}

func TestTreeHeight(t *testing.T) {
	btree := NewBTree[int, int](3)

	for i := 1; i <= 1000; i++ {
		btree.Insert(i, i)
	}

	height := btree.Height()

	if height > 10 {
		t.Errorf("Tree height %d seems too large for 1000 entries with degree 3", height)
	}

	if height < 2 {
		t.Errorf("Tree height %d seems too small for 1000 entries", height)
	}
}

// === Stress Tests ===

func TestStressInsertDelete(t *testing.T) {
	btree := NewBTree[int, int](4)

	for round := 0; round < 10; round++ {
		for i := 0; i < 1000; i++ {
			btree.Insert(i, i)
		}

		if btree.Size() != 1000 {
			t.Errorf("Round %d insert: expected size=1000, got=%d", round, btree.Size())
		}

		for i := 0; i < 1000; i++ {
			btree.Delete(i)
		}

		if !btree.IsEmpty() {
			t.Errorf("Round %d delete: expected empty tree, got size=%d", round, btree.Size())
		}
	}
}

func TestStressMixedOps(t *testing.T) {
	btree := NewBTree[int, int](3)
	expected := make(map[int]int)

	for i := 0; i < 5000; i++ {
		op := rand.Intn(10)
		key := rand.Intn(500)

		if op < 6 {
			value := rand.Intn(10000)
			if _, exists := expected[key]; !exists {
				btree.Insert(key, value)
				expected[key] = value
			}
		} else {
			btree.Delete(key)
			delete(expected, key)
		}
	}

	if btree.Size() != len(expected) {
		t.Errorf("Size mismatch: tree=%d, expected=%d", btree.Size(), len(expected))
	}

	for k, v := range expected {
		got, found := btree.Search(k)
		if !found {
			t.Errorf("Key %d not found", k)
		} else if got != v {
			t.Errorf("Key %d: expected %d, got %d", k, v, got)
		}
	}
}

func TestLargeDataset(t *testing.T) {
	btree := NewBTree[int, int](4)
	n := 10000

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys {
		btree.Insert(k, k*2)
	}

	if btree.Size() != n {
		t.Errorf("Size(): expected %d, got %d", n, btree.Size())
	}

	for _, k := range keys {
		value, found := btree.Search(k)
		if !found || value != k*2 {
			t.Errorf("Search(%d): expected %d, got %d, found=%v", k, k*2, value, found)
		}
	}

	for i := 0; i < n/2; i++ {
		btree.Delete(keys[i])
	}

	if btree.Size() != n/2 {
		t.Errorf("Size() after deletions: expected %d, got %d", n/2, btree.Size())
	}
}

// === Benchmarks ===

func BenchmarkBTreeInsertSequential(b *testing.B) {
	for _, degree := range []int{3, 10, 50} {
		b.Run(fmt.Sprintf("degree=%d", degree), func(b *testing.B) {
			btree := NewBTree[int, int](degree)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				btree.Insert(i, i)
			}
		})
	}
}

func BenchmarkBTreeInsertRandom(b *testing.B) {
	keys := make([]int, b.N)
	for i := range keys {
		keys[i] = rand.Int()
	}

	b.ResetTimer()
	btree := NewBTree[int, int](10)
	for i := 0; i < b.N; i++ {
		btree.Insert(keys[i], i)
	}
}

func BenchmarkBTreeSearch(b *testing.B) {
	btree := NewBTree[int, int](10)
	n := 100000

	for i := 0; i < n; i++ {
		btree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.Search(i % n)
	}
}

func BenchmarkBTreeSearchRandom(b *testing.B) {
	btree := NewBTree[int, int](10)
	n := 100000

	for i := 0; i < n; i++ {
		btree.Insert(i, i)
	}

	keys := make([]int, b.N)
	for i := range keys {
		keys[i] = rand.Intn(n)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.Search(keys[i])
	}
}

func BenchmarkBTreeDelete(b *testing.B) {
	for _, degree := range []int{3, 10, 50} {
		b.Run(fmt.Sprintf("degree=%d", degree), func(b *testing.B) {
			btree := NewBTree[int, int](degree)
			for i := 0; i < b.N; i++ {
				btree.Insert(i, i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				btree.Delete(i)
			}
		})
	}
}

func BenchmarkBTreeTraversal(b *testing.B) {
	btree := NewBTree[int, int](10)
	n := 10000

	for i := 0; i < n; i++ {
		btree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.InOrderTraversal()
	}
}

func BenchmarkBTreeMixedOps(b *testing.B) {
	btree := NewBTree[int, int](10)

	for i := 0; i < 10000; i++ {
		btree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := i % 10
		key := i % 10000

		switch {
		case op < 5:
			btree.Search(key)
		case op < 8:
			btree.Insert(key+10000, i)
		default:
			btree.Delete(key + 10000)
		}
	}
}
