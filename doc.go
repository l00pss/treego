// Package treego implements a generic B-tree data structure in Go.
//
// A B-tree is a self-balancing tree data structure that maintains sorted data
// and allows searches, sequential access, insertions, and deletions in logarithmic time.
// The B-tree generalizes the binary search tree, allowing for nodes with more than two children.
//
// This implementation provides:
//   - Generic types for both keys and values using Go generics
//   - Configurable minimum degree for performance tuning
//   - Complete set of operations: Insert, Search, Delete, Traversal
//   - High performance with O(log n) complexity for basic operations
//
// Example usage:
//
//	// Create a B-tree with integer keys and string values
//	bt := treego.NewBTree[int, string](3)
//
//	// Insert some data
//	bt.Insert(10, "ten")
//	bt.Insert(5, "five")
//	bt.Insert(20, "twenty")
//
//	// Search for a key
//	if value, found := bt.Search(10); found {
//	    fmt.Printf("Found: %s\n", value)
//	}
//
//	// Get all items in sorted order
//	items := bt.InOrderTraversal()
//	for _, item := range items {
//	    fmt.Printf("%d -> %s\n", item.Key, item.Value)
//	}
//
// The B-tree is particularly useful for:
//   - Database indexing
//   - File system implementations
//   - Large datasets that don't fit in memory
//   - Applications requiring sorted key-value storage
//
// Performance characteristics:
//   - Insert: O(log n)
//   - Search: O(log n)
//   - Delete: O(log n)
//   - Space: O(n)
//
// The minimum degree parameter affects performance:
//   - Lower degrees (2-4): Better for small datasets, less memory usage
//   - Higher degrees (16+): Better for large datasets, fewer disk accesses
package treego
