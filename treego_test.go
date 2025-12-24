package treego

import (
	"testing"
)

func TestBTreeBasicOperations(t *testing.T) {
	// Test with int keys and string values
	btree := NewBTree[int, string](3)

	// Test insertion
	btree.Insert(10, "ten")
	btree.Insert(20, "twenty")
	btree.Insert(5, "five")
	btree.Insert(6, "six")
	btree.Insert(12, "twelve")
	btree.Insert(30, "thirty")
	btree.Insert(7, "seven")
	btree.Insert(17, "seventeen")

	// Test search
	if val, found := btree.Search(10); !found || val != "ten" {
		t.Errorf("Expected to find '10' -> 'ten', got found=%v, val=%v", found, val)
	}

	if val, found := btree.Search(25); found {
		t.Errorf("Expected not to find '25', but found val=%v", val)
	}

	// Test size
	if size := btree.Size(); size != 8 {
		t.Errorf("Expected size 8, got %d", size)
	}

	// Test in-order traversal
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

	// Insert test data
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17, 25, 40, 50}
	for _, key := range keys {
		btree.Insert(key, "value")
	}

	// Test deletion
	if !btree.Delete(6) {
		t.Error("Expected to delete key 6")
	}

	if btree.Delete(100) {
		t.Error("Expected not to delete non-existent key 100")
	}

	// Verify deletion
	if _, found := btree.Search(6); found {
		t.Error("Key 6 should not exist after deletion")
	}

	// Test size after deletion
	if size := btree.Size(); size != 10 {
		t.Errorf("Expected size 10 after deletion, got %d", size)
	}
}

func TestBTreeWithStrings(t *testing.T) {
	btree := NewBTree[string, int](2)

	// Test with string keys
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

	// Insert one element
	btree.Insert(1, "one")

	if btree.IsEmpty() {
		t.Error("B-tree with one element should not be empty")
	}
}

// Benchmark tests
func BenchmarkBTreeInsert(b *testing.B) {
	btree := NewBTree[int, int](16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.Insert(i, i)
	}
}

func BenchmarkBTreeSearch(b *testing.B) {
	btree := NewBTree[int, int](16)

	// Pre-populate
	for i := 0; i < 10000; i++ {
		btree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.Search(i % 10000)
	}
}

func BenchmarkBTreeDelete(b *testing.B) {
	btree := NewBTree[int, int](16)

	// Pre-populate
	for i := 0; i < b.N*2; i++ {
		btree.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btree.Delete(i)
	}
}
