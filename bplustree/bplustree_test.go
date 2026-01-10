package bplustree

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func TestInsertAndSearch(t *testing.T) {
	tree := New[int, string](3)

	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")
	tree.Insert(5, "five")
	tree.Insert(15, "fifteen")
	tree.Insert(25, "twenty-five")
	tree.Insert(1, "one")
	tree.Insert(30, "thirty")

	tests := []struct {
		key      int
		expected string
		found    bool
	}{
		{10, "ten", true},
		{20, "twenty", true},
		{5, "five", true},
		{15, "fifteen", true},
		{25, "twenty-five", true},
		{1, "one", true},
		{30, "thirty", true},
		{100, "", false},
		{0, "", false},
	}

	for _, tc := range tests {
		value, found := tree.Search(tc.key)
		if found != tc.found {
			t.Errorf("Search(%d): expected found=%v, got=%v", tc.key, tc.found, found)
		}
		if found && value != tc.expected {
			t.Errorf("Search(%d): expected value=%s, got=%s", tc.key, tc.expected, value)
		}
	}
}

func TestUpdate(t *testing.T) {
	tree := New[int, string](3)

	tree.Insert(10, "original")
	tree.Insert(10, "updated")

	value, found := tree.Search(10)
	if !found || value != "updated" {
		t.Errorf("Update failed: expected 'updated', got '%s'", value)
	}
}

func TestSplit(t *testing.T) {
	tree := New[int, int](2)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, i*10)
	}

	for i := 1; i <= 10; i++ {
		value, found := tree.Search(i)
		if !found || value != i*10 {
			t.Errorf("Search(%d): expected %d, got %d, found=%v", i, i*10, value, found)
		}
	}

	if tree.Len() != 10 {
		t.Errorf("Len(): expected 10, got %d", tree.Len())
	}
}

func TestDelete(t *testing.T) {
	tree := New[int, string](3)

	tree.Insert(10, "ten")
	tree.Insert(20, "twenty")
	tree.Insert(5, "five")

	if !tree.Delete(10) {
		t.Error("Delete(10) should return true")
	}

	_, found := tree.Search(10)
	if found {
		t.Error("Search(10) should return false after deletion")
	}

	if tree.Delete(100) {
		t.Error("Delete(100) should return false for non-existent key")
	}

	if tree.Len() != 2 {
		t.Errorf("Len() after deletion: expected 2, got %d", tree.Len())
	}
}

func TestRange(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 20; i++ {
		tree.Insert(i, i*10)
	}

	result := tree.Range(5, 15)

	if len(result) != 11 {
		t.Errorf("Range(5, 15): expected 11 entries, got %d", len(result))
	}

	for i, e := range result {
		expectedKey := 5 + i
		if e.Key != expectedKey {
			t.Errorf("Range entry %d: expected key %d, got %d", i, expectedKey, e.Key)
		}
	}
}

func TestAll(t *testing.T) {
	tree := New[int, int](3)

	keys := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
	for _, k := range keys {
		tree.Insert(k, k*10)
	}

	all := tree.All()

	if len(all) != len(keys) {
		t.Errorf("All(): expected %d entries, got %d", len(keys), len(all))
	}

	for i := 1; i < len(all); i++ {
		if all[i-1].Key >= all[i].Key {
			t.Error("All() entries are not sorted")
			break
		}
	}
}

func TestEmptyTree(t *testing.T) {
	tree := New[int, string](3)

	_, found := tree.Search(1)
	if found {
		t.Error("Search on empty tree should return false")
	}

	if tree.Delete(1) {
		t.Error("Delete on empty tree should return false")
	}

	if tree.Len() != 0 {
		t.Errorf("Len() on empty tree: expected 0, got %d", tree.Len())
	}

	if tree.All() != nil {
		t.Error("All() on empty tree should return nil")
	}

	if tree.Range(1, 10) != nil {
		t.Error("Range() on empty tree should return nil")
	}
}

func TestLargeDataset(t *testing.T) {
	tree := New[int, int](4)
	n := 10000

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, k := range keys {
		tree.Insert(k, k*2)
	}

	if tree.Len() != n {
		t.Errorf("Len(): expected %d, got %d", n, tree.Len())
	}

	for _, k := range keys {
		value, found := tree.Search(k)
		if !found || value != k*2 {
			t.Errorf("Search(%d): expected %d, got %d, found=%v", k, k*2, value, found)
		}
	}

	all := tree.All()
	if !slices.IsSortedFunc(all, func(a, b Entry[int, int]) int {
		return a.Key - b.Key
	}) {
		t.Error("All() entries are not sorted")
	}

	for i := 0; i < n/2; i++ {
		tree.Delete(keys[i])
	}

	if tree.Len() != n/2 {
		t.Errorf("Len() after deletions: expected %d, got %d", n/2, tree.Len())
	}
}

func TestStringKeys(t *testing.T) {
	tree := New[string, int](3)

	tree.Insert("apple", 1)
	tree.Insert("banana", 2)
	tree.Insert("cherry", 3)
	tree.Insert("date", 4)

	value, found := tree.Search("banana")
	if !found || value != 2 {
		t.Errorf("Search('banana'): expected 2, got %d", value)
	}

	result := tree.Range("banana", "date")
	if len(result) != 3 {
		t.Errorf("Range: expected 3 entries, got %d", len(result))
	}
}

func TestLeafLinking(t *testing.T) {
	tree := New[int, int](2)

	for i := 1; i <= 20; i++ {
		tree.Insert(i, i)
	}

	all := tree.All()
	for i, e := range all {
		if e.Key != i+1 {
			t.Errorf("Leaf linking broken at index %d: expected key %d, got %d", i, i+1, e.Key)
		}
	}
}

// === Edge Cases ===

func TestDeleteSingleElement(t *testing.T) {
	tree := New[int, int](3)
	tree.Insert(1, 10)

	if !tree.Delete(1) {
		t.Error("Delete single element should return true")
	}

	if tree.Len() != 0 {
		t.Errorf("Tree should be empty, got len=%d", tree.Len())
	}

	if tree.root != nil {
		t.Error("Root should be nil after deleting single element")
	}
}

func TestDeleteAllElements(t *testing.T) {
	tree := New[int, int](3)
	n := 100

	for i := 1; i <= n; i++ {
		tree.Insert(i, i*10)
	}

	for i := 1; i <= n; i++ {
		if !tree.Delete(i) {
			t.Errorf("Delete(%d) should return true", i)
		}
	}

	if tree.Len() != 0 {
		t.Errorf("Tree should be empty after deleting all, got len=%d", tree.Len())
	}

	if tree.root != nil {
		t.Error("Root should be nil after deleting all elements")
	}
}

func TestDeleteReverseOrder(t *testing.T) {
	tree := New[int, int](3)
	n := 50

	for i := 1; i <= n; i++ {
		tree.Insert(i, i)
	}

	for i := n; i >= 1; i-- {
		if !tree.Delete(i) {
			t.Errorf("Delete(%d) should return true", i)
		}
		if tree.Len() != i-1 {
			t.Errorf("After deleting %d, expected len=%d, got=%d", i, i-1, tree.Len())
		}
	}
}

func TestDeleteRandomOrder(t *testing.T) {
	tree := New[int, int](3)
	n := 100

	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = i
		tree.Insert(i, i*10)
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for i, k := range keys {
		if !tree.Delete(k) {
			t.Errorf("Delete(%d) should return true", k)
		}
		if tree.Len() != n-i-1 {
			t.Errorf("After deleting %d keys, expected len=%d, got=%d", i+1, n-i-1, tree.Len())
		}
	}
}

func TestDeleteAndReinsert(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 20; i++ {
		tree.Insert(i, i*10)
	}

	for i := 1; i <= 10; i++ {
		tree.Delete(i)
	}

	for i := 1; i <= 10; i++ {
		tree.Insert(i, i*100)
	}

	for i := 1; i <= 10; i++ {
		value, found := tree.Search(i)
		if !found || value != i*100 {
			t.Errorf("Reinserted key %d: expected %d, got %d", i, i*100, value)
		}
	}

	if tree.Len() != 20 {
		t.Errorf("Expected len=20, got=%d", tree.Len())
	}
}

// === Range Edge Cases ===

func TestRangeEmptyResult(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, i)
	}

	result := tree.Range(100, 200)
	if len(result) != 0 {
		t.Errorf("Range with no matches should return empty, got %d", len(result))
	}
}

func TestRangeSingleResult(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, i*10)
	}

	result := tree.Range(5, 5)
	if len(result) != 1 {
		t.Errorf("Range(5,5) should return 1 entry, got %d", len(result))
	}
	if result[0].Key != 5 {
		t.Errorf("Range(5,5) should return key 5, got %d", result[0].Key)
	}
}

func TestRangeAllElements(t *testing.T) {
	tree := New[int, int](3)
	n := 50

	for i := 1; i <= n; i++ {
		tree.Insert(i, i)
	}

	result := tree.Range(1, n)
	if len(result) != n {
		t.Errorf("Range(1,%d) should return %d entries, got %d", n, n, len(result))
	}
}

func TestRangeBoundary(t *testing.T) {
	tree := New[int, int](3)

	for i := 10; i <= 100; i += 10 {
		tree.Insert(i, i)
	}

	result := tree.Range(25, 75)
	expected := []int{30, 40, 50, 60, 70}

	if len(result) != len(expected) {
		t.Errorf("Expected %d entries, got %d", len(expected), len(result))
	}

	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("Entry %d: expected key %d, got %d", i, expected[i], e.Key)
		}
	}
}

// === Degree Boundary Tests ===

func TestMinimumDegree(t *testing.T) {
	tree := New[int, int](2)

	for i := 1; i <= 100; i++ {
		tree.Insert(i, i)
	}

	if tree.Len() != 100 {
		t.Errorf("Expected len=100, got=%d", tree.Len())
	}

	for i := 1; i <= 100; i++ {
		value, found := tree.Search(i)
		if !found || value != i {
			t.Errorf("Search(%d) failed", i)
		}
	}
}

func TestLargeDegree(t *testing.T) {
	tree := New[int, int](50)

	for i := 1; i <= 1000; i++ {
		tree.Insert(i, i)
	}

	if tree.Len() != 1000 {
		t.Errorf("Expected len=1000, got=%d", tree.Len())
	}

	for i := 1; i <= 1000; i++ {
		if _, found := tree.Search(i); !found {
			t.Errorf("Search(%d) should find key", i)
		}
	}
}

func TestDegreeOne(t *testing.T) {
	tree := New[int, int](1)

	for i := 1; i <= 10; i++ {
		tree.Insert(i, i)
	}

	if tree.Len() != 10 {
		t.Errorf("Expected len=10, got=%d", tree.Len())
	}
}

// === Tree Structure Validation ===

func (t *BPlusTree[K, V]) validate() error {
	if t.root == nil {
		return nil
	}
	return t.validateNode(t.root, nil, nil, 0)
}

func (t *BPlusTree[K, V]) validateNode(n *node[K, V], minKey, maxKey *K, depth int) error {
	if n.isLeaf {
		if len(n.entries) > t.maxLeafEntries() {
			return fmt.Errorf("leaf has too many entries: %d > %d", len(n.entries), t.maxLeafEntries())
		}

		if n != t.root && len(n.entries) < t.minLeafEntries() {
			return fmt.Errorf("non-root leaf has too few entries: %d < %d", len(n.entries), t.minLeafEntries())
		}

		for i := 1; i < len(n.entries); i++ {
			if n.entries[i-1].Key >= n.entries[i].Key {
				return fmt.Errorf("leaf entries not sorted at index %d", i)
			}
		}

		for _, e := range n.entries {
			if minKey != nil && e.Key < *minKey {
				return fmt.Errorf("leaf key %v < minKey %v", e.Key, *minKey)
			}
			if maxKey != nil && e.Key >= *maxKey {
				return fmt.Errorf("leaf key %v >= maxKey %v", e.Key, *maxKey)
			}
		}
	} else {
		if len(n.keys) > t.maxInternalKeys() {
			return fmt.Errorf("internal node has too many keys: %d > %d", len(n.keys), t.maxInternalKeys())
		}

		if n != t.root && len(n.keys) < t.minInternalKeys() {
			return fmt.Errorf("non-root internal has too few keys: %d < %d", len(n.keys), t.minInternalKeys())
		}

		if len(n.children) != len(n.keys)+1 {
			return fmt.Errorf("internal node children count mismatch: %d children, %d keys", len(n.children), len(n.keys))
		}

		for i := 1; i < len(n.keys); i++ {
			if n.keys[i-1] >= n.keys[i] {
				return fmt.Errorf("internal keys not sorted at index %d", i)
			}
		}

		for i, child := range n.children {
			if child.parent != n {
				return fmt.Errorf("child %d has wrong parent", i)
			}

			var childMin, childMax *K
			if i > 0 {
				childMin = &n.keys[i-1]
			}
			if i < len(n.keys) {
				childMax = &n.keys[i]
			}

			if err := t.validateNode(child, childMin, childMax, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *BPlusTree[K, V]) height() int {
	if t.root == nil {
		return 0
	}
	h := 1
	n := t.root
	for !n.isLeaf {
		h++
		n = n.children[0]
	}
	return h
}

func (t *BPlusTree[K, V]) countLeaves() int {
	if t.root == nil {
		return 0
	}
	count := 0
	leaf := t.firstLeaf()
	for leaf != nil {
		count++
		leaf = leaf.next
	}
	return count
}

func TestTreeStructureAfterInserts(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 100; i++ {
		tree.Insert(i, i)
		if err := tree.validate(); err != nil {
			t.Errorf("Invalid tree after inserting %d: %v", i, err)
		}
	}
}

func TestTreeStructureAfterDeletes(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 50; i++ {
		tree.Insert(i, i)
	}

	for i := 1; i <= 50; i++ {
		tree.Delete(i)
		if err := tree.validate(); err != nil {
			t.Errorf("Invalid tree after deleting %d: %v", i, err)
		}
	}
}

func TestTreeStructureRandomOps(t *testing.T) {
	tree := New[int, int](3)

	for i := 0; i < 500; i++ {
		op := rand.Intn(3)
		key := rand.Intn(100)

		switch op {
		case 0, 1:
			tree.Insert(key, key*10)
		case 2:
			tree.Delete(key)
		}

		if err := tree.validate(); err != nil {
			t.Errorf("Invalid tree at iteration %d (op=%d, key=%d): %v", i, op, key, err)
		}
	}
}

func TestLeafChainIntegrity(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 100; i++ {
		tree.Insert(i, i)
	}

	leafCount := tree.countLeaves()
	if leafCount == 0 {
		t.Error("Tree should have leaves")
	}

	leaf := tree.firstLeaf()
	visited := 0
	var prevKey int

	for leaf != nil {
		for _, e := range leaf.entries {
			if visited > 0 && e.Key <= prevKey {
				t.Errorf("Leaf chain not sorted: %d after %d", e.Key, prevKey)
			}
			prevKey = e.Key
			visited++
		}
		leaf = leaf.next
	}

	if visited != tree.Len() {
		t.Errorf("Leaf chain missing entries: visited=%d, len=%d", visited, tree.Len())
	}
}

func TestTreeHeight(t *testing.T) {
	tree := New[int, int](3)

	for i := 1; i <= 1000; i++ {
		tree.Insert(i, i)
	}

	height := tree.height()

	maxHeight := 10
	if height > maxHeight {
		t.Errorf("Tree height %d seems too large for 1000 entries with degree 3", height)
	}

	if height < 2 {
		t.Errorf("Tree height %d seems too small for 1000 entries", height)
	}
}

// === Stress Tests ===

func TestStressInsertDelete(t *testing.T) {
	tree := New[int, int](4)

	for round := 0; round < 10; round++ {
		for i := 0; i < 1000; i++ {
			tree.Insert(i, i)
		}

		if tree.Len() != 1000 {
			t.Errorf("Round %d insert: expected len=1000, got=%d", round, tree.Len())
		}

		for i := 0; i < 1000; i++ {
			tree.Delete(i)
		}

		if tree.Len() != 0 {
			t.Errorf("Round %d delete: expected len=0, got=%d", round, tree.Len())
		}
	}
}

func TestStressMixedOps(t *testing.T) {
	tree := New[int, int](3)
	expected := make(map[int]int)

	for i := 0; i < 5000; i++ {
		op := rand.Intn(10)
		key := rand.Intn(500)

		if op < 6 {
			value := rand.Intn(10000)
			tree.Insert(key, value)
			expected[key] = value
		} else {
			tree.Delete(key)
			delete(expected, key)
		}
	}

	if tree.Len() != len(expected) {
		t.Errorf("Length mismatch: tree=%d, expected=%d", tree.Len(), len(expected))
	}

	for k, v := range expected {
		got, found := tree.Search(k)
		if !found {
			t.Errorf("Key %d not found", k)
		} else if got != v {
			t.Errorf("Key %d: expected %d, got %d", k, v, got)
		}
	}
}

// === Benchmarks ===

func BenchmarkInsertSequential(b *testing.B) {
	for _, degree := range []int{3, 10, 50} {
		b.Run(fmt.Sprintf("degree=%d", degree), func(b *testing.B) {
			tree := New[int, int](degree)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Insert(i, i)
			}
		})
	}
}

func BenchmarkInsertRandom(b *testing.B) {
	keys := make([]int, b.N)
	for i := range keys {
		keys[i] = rand.Int()
	}

	b.ResetTimer()
	tree := New[int, int](10)
	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i], i)
	}
}

func BenchmarkSearch(b *testing.B) {
	tree := New[int, int](10)
	n := 100000

	for i := 0; i < n; i++ {
		tree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(i % n)
	}
}

func BenchmarkSearchRandom(b *testing.B) {
	tree := New[int, int](10)
	n := 100000

	for i := 0; i < n; i++ {
		tree.Insert(i, i)
	}

	keys := make([]int, b.N)
	for i := range keys {
		keys[i] = rand.Intn(n)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(keys[i])
	}
}

func BenchmarkDelete(b *testing.B) {
	for _, degree := range []int{3, 10, 50} {
		b.Run(fmt.Sprintf("degree=%d", degree), func(b *testing.B) {
			tree := New[int, int](degree)
			for i := 0; i < b.N; i++ {
				tree.Insert(i, i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Delete(i)
			}
		})
	}
}

func BenchmarkRange(b *testing.B) {
	tree := New[int, int](10)
	n := 100000

	for i := 0; i < n; i++ {
		tree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := i % (n - 100)
		tree.Range(start, start+100)
	}
}

func BenchmarkMixedOps(b *testing.B) {
	tree := New[int, int](10)

	for i := 0; i < 10000; i++ {
		tree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := i % 10
		key := i % 10000

		switch {
		case op < 5:
			tree.Search(key)
		case op < 8:
			tree.Insert(key+10000, i)
		default:
			tree.Delete(key + 10000)
		}
	}
}
