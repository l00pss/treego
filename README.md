# TreeGo

Generic tree data structures for Go. Built for real-world use.

![TreeGo Logo](logo.png)

<div align="center">
  <a href="https://www.buymeacoffee.com/l00pss" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
</div>

## What's Inside

| Package | Description | Best For |
|---------|-------------|----------|
| `btree` | Classic B-Tree | General purpose key-value storage |
| `bplustree` | B+ Tree | Range queries, databases, sequential access |

## Installation

```bash
go get github.com/l00pss/treego/btree
go get github.com/l00pss/treego/bplustree
```

## Quick Examples

### B-Tree

```go
import "github.com/l00pss/treego/btree"

bt := btree.NewBTree[int, string](3)

bt.Insert(10, "ten")
bt.Insert(5, "five")
bt.Insert(20, "twenty")

value, found := bt.Search(10)  // "ten", true

bt.Delete(5)

for _, item := range bt.InOrderTraversal() {
    fmt.Printf("%d: %s\n", item.Key, item.Value)
}
```

### B+ Tree

```go
import "github.com/l00pss/treego/bplustree"

tree := bplustree.New[int, string](3)

tree.Insert(10, "ten")
tree.Insert(5, "five")
tree.Insert(20, "twenty")

value, found := tree.Search(10)  // "ten", true

// Range query - B+ Tree's superpower
for _, e := range tree.Range(5, 15) {
    fmt.Printf("%d: %s\n", e.Key, e.Value)
}

tree.Delete(5)
```

## B-Tree vs B+ Tree

| Feature | B-Tree | B+ Tree |
|---------|--------|---------|
| Data location | All nodes | Leaf nodes only |
| Range queries | Slower | Fast (linked leaves) |
| Point lookups | Fast | Fast |
| Sequential scan | No | Yes |
| Use case | General storage | Databases, file systems |

## API

### B-Tree

```go
NewBTree[K, V](degree)      // Create tree
Insert(key, value)          // Add or update
Search(key) (V, bool)       // Find by key
Delete(key) bool            // Remove key
InOrderTraversal() []KV     // All items sorted
Size() int                  // Count of items
Height() int                // Tree height
IsEmpty() bool              // Check if empty
```

### B+ Tree

```go
New[K, V](degree)           // Create tree
Insert(key, value)          // Add or update
Search(key) (V, bool)       // Find by key
Delete(key) bool            // Remove key
Range(start, end) []Entry   // Range query
All() []Entry               // All items sorted
Len() int                   // Count of items
```

## Benchmarks

Tested on Apple M1 Pro.

### B-Tree

| Operation | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| Insert (degree=3) | 135.8 | 163 | 3 |
| Insert (degree=10) | 41.5 | 64 | 0 |
| Insert (degree=50) | 22.6 | 55 | 0 |
| Search | 32.9 | 0 | 0 |
| Search (random) | 83.3 | 0 | 0 |
| Delete (degree=3) | 56.0 | 0 | 0 |
| Delete (degree=10) | 33.1 | 0 | 0 |
| Delete (degree=50) | 34.5 | 0 | 0 |
| Mixed ops | 83.0 | 22 | 0 |

### B+ Tree

| Operation | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| Insert (degree=3) | 113.3 | 121 | 1 |
| Insert (degree=10) | 79.3 | 69 | 0 |
| Insert (degree=50) | 145.1 | 56 | 0 |
| Search | 31.6 | 0 | 0 |
| Search (random) | 84.7 | 0 | 0 |
| Delete (degree=3) | 40.0 | 0 | 0 |
| Delete (degree=10) | 26.6 | 0 | 0 |
| Delete (degree=50) | 27.6 | 0 | 0 |
| Range (100 items) | 739.0 | 4080 | 8 |
| Mixed ops | 32.3 | 0 | 0 |

## Running Tests

```bash
go test ./btree/... ./bplustree/... -v

# Benchmarks
go test ./btree/... ./bplustree/... -bench=. -benchmem
```

## License

MIT
