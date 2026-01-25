package rtree

import (
	"math"
	"testing"
)

// TestNewRTree tests R-tree creation
func TestNewRTree(t *testing.T) {
	tree := NewRTree(2, 4)

	if tree.minEntries != 2 {
		t.Errorf("Expected minEntries to be 2, got %d", tree.minEntries)
	}

	if tree.maxEntries != 4 {
		t.Errorf("Expected maxEntries to be 4, got %d", tree.maxEntries)
	}

	if tree.size != 0 {
		t.Errorf("Expected size to be 0, got %d", tree.size)
	}

	if tree.root == nil {
		t.Error("Root should not be nil")
	}

	if !tree.root.isLeaf {
		t.Error("Root should be a leaf initially")
	}
}

// TestRectangleArea tests rectangle area calculation
func TestRectangleArea(t *testing.T) {
	tests := []struct {
		rect     Rectangle
		expected float64
	}{
		{NewRectangle(0, 0, 10, 10), 100.0},
		{NewRectangle(0, 0, 5, 5), 25.0},
		{NewRectangle(2, 3, 7, 8), 25.0},
		{NewRectangle(0, 0, 0, 0), 0.0},
	}

	for _, test := range tests {
		area := test.rect.Area()
		if area != test.expected {
			t.Errorf("Expected area %.2f, got %.2f for rect %+v",
				test.expected, area, test.rect)
		}
	}
}

// TestRectangleIntersects tests rectangle intersection
func TestRectangleIntersects(t *testing.T) {
	tests := []struct {
		r1       Rectangle
		r2       Rectangle
		expected bool
	}{
		{NewRectangle(0, 0, 10, 10), NewRectangle(5, 5, 15, 15), true},
		{NewRectangle(0, 0, 10, 10), NewRectangle(11, 11, 20, 20), false},
		{NewRectangle(0, 0, 10, 10), NewRectangle(2, 2, 8, 8), true},
		{NewRectangle(0, 0, 10, 10), NewRectangle(10, 10, 20, 20), true},
		{NewRectangle(0, 0, 5, 5), NewRectangle(6, 6, 10, 10), false},
	}

	for _, test := range tests {
		result := test.r1.Intersects(test.r2)
		if result != test.expected {
			t.Errorf("Expected %v.Intersects(%v) to be %v, got %v",
				test.r1, test.r2, test.expected, result)
		}
	}
}

// TestRectangleContains tests rectangle containment
func TestRectangleContains(t *testing.T) {
	tests := []struct {
		r1       Rectangle
		r2       Rectangle
		expected bool
	}{
		{NewRectangle(0, 0, 10, 10), NewRectangle(2, 2, 8, 8), true},
		{NewRectangle(0, 0, 10, 10), NewRectangle(5, 5, 15, 15), false},
		{NewRectangle(0, 0, 10, 10), NewRectangle(0, 0, 10, 10), true},
		{NewRectangle(2, 2, 8, 8), NewRectangle(0, 0, 10, 10), false},
	}

	for _, test := range tests {
		result := test.r1.Contains(test.r2)
		if result != test.expected {
			t.Errorf("Expected %v.Contains(%v) to be %v, got %v",
				test.r1, test.r2, test.expected, result)
		}
	}
}

// TestRectangleContainsPoint tests point containment
func TestRectangleContainsPoint(t *testing.T) {
	rect := NewRectangle(0, 0, 10, 10)

	tests := []struct {
		point    Point
		expected bool
	}{
		{Point{5, 5}, true},
		{Point{0, 0}, true},
		{Point{10, 10}, true},
		{Point{11, 11}, false},
		{Point{-1, 5}, false},
		{Point{5, 11}, false},
	}

	for _, test := range tests {
		result := rect.ContainsPoint(test.point)
		if result != test.expected {
			t.Errorf("Expected rect.ContainsPoint(%v) to be %v, got %v",
				test.point, test.expected, result)
		}
	}
}

// TestRectangleUnion tests rectangle union
func TestRectangleUnion(t *testing.T) {
	r1 := NewRectangle(0, 0, 5, 5)
	r2 := NewRectangle(3, 3, 8, 8)

	union := r1.Union(r2)
	expected := NewRectangle(0, 0, 8, 8)

	if union != expected {
		t.Errorf("Expected union to be %v, got %v", expected, union)
	}
}

// TestRectangleDistance tests distance calculation
func TestRectangleDistance(t *testing.T) {
	rect := NewRectangle(0, 0, 10, 10)

	tests := []struct {
		point    Point
		expected float64
	}{
		{Point{5, 5}, 0.0},             // Inside
		{Point{15, 5}, 5.0},            // Right
		{Point{5, 15}, 5.0},            // Top
		{Point{-5, 5}, 5.0},            // Left
		{Point{5, -5}, 5.0},            // Bottom
		{Point{15, 15}, math.Sqrt(50)}, // Diagonal
	}

	for _, test := range tests {
		result := rect.Distance(test.point)
		if math.Abs(result-test.expected) > 0.0001 {
			t.Errorf("Expected distance from rect to %v to be %.4f, got %.4f",
				test.point, test.expected, result)
		}
	}
}

// TestInsertSingleItem tests inserting a single item
func TestInsertSingleItem(t *testing.T) {
	tree := NewRTree(2, 4)

	item := &Item{
		Bounds: NewRectangle(0, 0, 10, 10),
		Data:   "Test Item",
	}

	tree.Insert(item)

	if tree.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", tree.Size())
	}

	if !tree.root.isLeaf {
		t.Error("Root should still be a leaf after one insertion")
	}

	if len(tree.root.items) != 1 {
		t.Errorf("Expected root to have 1 item, got %d", len(tree.root.items))
	}
}

// TestInsertMultipleItems tests inserting multiple items
func TestInsertMultipleItems(t *testing.T) {
	tree := NewRTree(2, 4)

	items := []struct {
		bounds Rectangle
		data   string
	}{
		{NewRectangle(0, 0, 10, 10), "Item 1"},
		{NewRectangle(20, 20, 30, 30), "Item 2"},
		{NewRectangle(5, 5, 15, 15), "Item 3"},
		{NewRectangle(25, 25, 35, 35), "Item 4"},
	}

	for _, item := range items {
		tree.Insert(&Item{Bounds: item.bounds, Data: item.data})
	}

	if tree.Size() != 4 {
		t.Errorf("Expected size to be 4, got %d", tree.Size())
	}
}

// TestInsertTriggersSplit tests that insertion triggers node splitting
func TestInsertTriggersSplit(t *testing.T) {
	tree := NewRTree(2, 4)

	// Insert 5 items to trigger a split (max is 4)
	for i := 0; i < 5; i++ {
		tree.Insert(&Item{
			Bounds: NewRectangle(float64(i*10), float64(i*10),
				float64(i*10+5), float64(i*10+5)),
			Data: i,
		})
	}

	if tree.Size() != 5 {
		t.Errorf("Expected size to be 5, got %d", tree.Size())
	}

	// After split, root should no longer be a leaf
	if tree.root.isLeaf {
		t.Error("Root should not be a leaf after split")
	}

	if len(tree.root.children) < 2 {
		t.Errorf("Root should have at least 2 children after split, got %d",
			len(tree.root.children))
	}
}

// TestSearch tests basic search functionality
func TestSearch(t *testing.T) {
	tree := NewRTree(2, 4)

	// Insert test items
	tree.Insert(&Item{Bounds: NewRectangle(0, 0, 10, 10), Data: "A"})
	tree.Insert(&Item{Bounds: NewRectangle(20, 20, 30, 30), Data: "B"})
	tree.Insert(&Item{Bounds: NewRectangle(5, 5, 15, 15), Data: "C"})
	tree.Insert(&Item{Bounds: NewRectangle(100, 100, 110, 110), Data: "D"})

	// Search for items intersecting with (0, 0, 20, 20)
	results := tree.Search(NewRectangle(0, 0, 20, 20))

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify correct items are returned
	found := make(map[string]bool)
	for _, item := range results {
		found[item.Data.(string)] = true
	}

	expected := []string{"A", "B", "C"}
	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("Expected to find item %s", exp)
		}
	}

	if found["D"] {
		t.Error("Item D should not be in results")
	}
}

// TestSearchEmpty tests search on empty tree
func TestSearchEmpty(t *testing.T) {
	tree := NewRTree(2, 4)

	results := tree.Search(NewRectangle(0, 0, 10, 10))

	if len(results) != 0 {
		t.Errorf("Expected 0 results from empty tree, got %d", len(results))
	}
}

// TestSearchNoIntersection tests search with no intersecting items
func TestSearchNoIntersection(t *testing.T) {
	tree := NewRTree(2, 4)

	tree.Insert(&Item{Bounds: NewRectangle(0, 0, 10, 10), Data: "A"})
	tree.Insert(&Item{Bounds: NewRectangle(20, 20, 30, 30), Data: "B"})

	results := tree.Search(NewRectangle(100, 100, 110, 110))

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestSearchPoint tests point search functionality
func TestSearchPoint(t *testing.T) {
	tree := NewRTree(2, 4)

	tree.Insert(&Item{Bounds: NewRectangle(0, 0, 10, 10), Data: "A"})
	tree.Insert(&Item{Bounds: NewRectangle(20, 20, 30, 30), Data: "B"})
	tree.Insert(&Item{Bounds: NewRectangle(5, 5, 15, 15), Data: "C"})

	// Point inside A and C
	results := tree.SearchPoint(Point{7, 7})

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Point inside B only
	results = tree.SearchPoint(Point{25, 25})

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Data.(string) != "B" {
		t.Errorf("Expected to find item B, got %v", results[0].Data)
	}

	// Point outside all
	results = tree.SearchPoint(Point{100, 100})

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestNearestNeighbor tests k-nearest neighbor search
func TestNearestNeighbor(t *testing.T) {
	tree := NewRTree(2, 4)

	// Insert items at different locations
	tree.Insert(&Item{Bounds: NewPoint(0, 0), Data: "A"})
	tree.Insert(&Item{Bounds: NewPoint(10, 0), Data: "B"})
	tree.Insert(&Item{Bounds: NewPoint(5, 5), Data: "C"})
	tree.Insert(&Item{Bounds: NewPoint(20, 20), Data: "D"})
	tree.Insert(&Item{Bounds: NewPoint(30, 30), Data: "E"})

	// Find 3 nearest neighbors to (0, 0)
	results := tree.NearestNeighbor(Point{0, 0}, 3)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First result should be A (distance 0)
	if results[0].Data.(string) != "A" {
		t.Errorf("Expected first result to be A, got %v", results[0].Data)
	}
}

// TestNearestNeighborSingle tests single nearest neighbor
func TestNearestNeighborSingle(t *testing.T) {
	tree := NewRTree(2, 4)

	tree.Insert(&Item{Bounds: NewPoint(0, 0), Data: "A"})
	tree.Insert(&Item{Bounds: NewPoint(10, 10), Data: "B"})
	tree.Insert(&Item{Bounds: NewPoint(5, 5), Data: "C"})

	results := tree.NearestNeighbor(Point{6, 6}, 1)

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Data.(string) != "C" {
		t.Errorf("Expected nearest neighbor to be C, got %v", results[0].Data)
	}
}

// TestHeight tests tree height calculation
func TestHeight(t *testing.T) {
	tree := NewRTree(2, 4)

	// Empty tree should have height 1
	if tree.Height() != 1 {
		t.Errorf("Expected height of empty tree to be 1, got %d", tree.Height())
	}

	// Add items
	for i := 0; i < 3; i++ {
		tree.Insert(&Item{
			Bounds: NewRectangle(float64(i), float64(i), float64(i+1), float64(i+1)),
			Data:   i,
		})
	}

	height := tree.Height()
	if height < 1 {
		t.Errorf("Height should be at least 1, got %d", height)
	}
}

// TestLargeDataset tests with a larger dataset
func TestLargeDataset(t *testing.T) {
	tree := NewRTree(4, 16)

	// Insert 100 items
	n := 100
	for i := 0; i < n; i++ {
		x := float64(i % 10 * 10)
		y := float64(i / 10 * 10)
		tree.Insert(&Item{
			Bounds: NewRectangle(x, y, x+5, y+5),
			Data:   i,
		})
	}

	if tree.Size() != n {
		t.Errorf("Expected size to be %d, got %d", n, tree.Size())
	}

	// Search for items in a specific region
	results := tree.Search(NewRectangle(0, 0, 25, 25))

	if len(results) == 0 {
		t.Error("Expected to find items in search region")
	}

	// Test nearest neighbor search
	nearest := tree.NearestNeighbor(Point{50, 50}, 5)

	if len(nearest) != 5 {
		t.Errorf("Expected 5 nearest neighbors, got %d", len(nearest))
	}
}

// TestOverlappingRectangles tests handling of overlapping rectangles
func TestOverlappingRectangles(t *testing.T) {
	tree := NewRTree(2, 4)

	// Insert overlapping rectangles
	tree.Insert(&Item{Bounds: NewRectangle(0, 0, 20, 20), Data: "Large"})
	tree.Insert(&Item{Bounds: NewRectangle(5, 5, 10, 10), Data: "Small1"})
	tree.Insert(&Item{Bounds: NewRectangle(15, 15, 25, 25), Data: "Overlap"})

	// Search in overlapping region
	results := tree.Search(NewRectangle(7, 7, 17, 17))

	// Should find all three rectangles
	if len(results) != 3 {
		t.Errorf("Expected 3 overlapping results, got %d", len(results))
	}
}

// BenchmarkInsert benchmarks insertion performance
func BenchmarkInsert(b *testing.B) {
	tree := NewRTree(4, 16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := float64(i % 100)
		y := float64(i / 100)
		tree.Insert(&Item{
			Bounds: NewRectangle(x, y, x+1, y+1),
			Data:   i,
		})
	}
}

// BenchmarkSearch benchmarks search performance
func BenchmarkSearch(b *testing.B) {
	tree := NewRTree(4, 16)

	// Populate tree
	for i := 0; i < 1000; i++ {
		x := float64(i % 100)
		y := float64(i / 100)
		tree.Insert(&Item{
			Bounds: NewRectangle(x, y, x+1, y+1),
			Data:   i,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := float64(i % 50)
		y := float64(i % 50)
		tree.Search(NewRectangle(x, y, x+10, y+10))
	}
}

// BenchmarkNearestNeighbor benchmarks k-NN performance
func BenchmarkNearestNeighbor(b *testing.B) {
	tree := NewRTree(4, 16)

	// Populate tree
	for i := 0; i < 1000; i++ {
		x := float64(i % 100)
		y := float64(i / 100)
		tree.Insert(&Item{
			Bounds: NewPoint(x, y),
			Data:   i,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := float64(i % 50)
		y := float64(i % 50)
		tree.NearestNeighbor(Point{x, y}, 10)
	}
}
