# TreeGo 

A high-performance, generic B-tree implementation in Go with full type safety and comprehensive functionality.

![TreeGo Logo](logo.png)

## Features

- **Generic Implementation**: Works with any ordered type (int, string, float, etc.)
- **Type Safe**: Full compile-time type checking with Go generics
- **High Performance**: Optimized B-tree operations with configurable degree
- **Complete API**: Insert, Search, Delete, Traversal, and utility methods
- **Well Tested**: Comprehensive test suite with benchmarks
- **Zero Dependencies**: Pure Go implementation

## Installation

```bash
go get github.com/l00pss/treego
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/l00pss/treego"
)

func main() {
    // Create a B-tree with int keys and string values, degree 3
    bt := treego.NewBTree[int, string](3)
    
    // Insert key-value pairs
    bt.Insert(10, "ten")
    bt.Insert(20, "twenty")
    bt.Insert(5, "five")
    
    // Search for a key
    if value, found := bt.Search(10); found {
        fmt.Printf("Found: %s\n", value) // Output: Found: ten
    }
    
    // Get all items in sorted order
    items := bt.InOrderTraversal()
    for _, item := range items {
        fmt.Printf("%d -> %s\n", item.Key, item.Value)
    }
    // Output:
    // 5 -> five
    // 10 -> ten  
    // 20 -> twenty
}
```

## API Reference

### Creating a B-tree

```go
// Create a B-tree with minimum degree t
// Higher degree = wider tree, better for datasets in memory
// Lower degree = taller tree, better for disk-based storage
bt := treego.NewBTree[KeyType, ValueType](degree)
```

### Core Operations

#### Insert
```go
bt.Insert(key, value)  // O(log n)
```

#### Search  
```go
value, found := bt.Search(key)  // O(log n)
```

#### Delete
```go
deleted := bt.Delete(key)  // O(log n) - returns true if key existed
```

### Traversal and Inspection

#### In-order Traversal
```go
items := bt.InOrderTraversal()  // Returns []KeyValue[K, V] in sorted order
```

#### Tree Properties
```go
size := bt.Size()        // Total number of keys
height := bt.Height()    // Height of the tree
empty := bt.IsEmpty()    // Check if tree is empty
```

#### Tree Visualization
```go
fmt.Print(bt.String())   // Print tree structure
```

## Supported Types

The B-tree works with any type that satisfies the `Ordered` constraint:

- **Integers**: `int`, `int8`, `int16`, `int32`, `int64`
- **Unsigned integers**: `uint`, `uint8`, `uint16`, `uint32`, `uint64`  
- **Floating point**: `float32`, `float64`
- **Strings**: `string`

## Examples

### Example 1: String Dictionary
```go
dict := treego.NewBTree[string, string](3)

dict.Insert("hello", "мərхaba")
dict.Insert("world", "dünya") 
dict.Insert("tree", "ağac")

if translation, found := dict.Search("hello"); found {
    fmt.Printf("hello -> %s\n", translation)
}
```

### Example 2: Student Grades
```go
grades := treego.NewBTree[int, float64](4)

grades.Insert(12345, 85.5)  // Student ID -> Grade
grades.Insert(12346, 92.0)
grades.Insert(12344, 78.5)

// Get all grades in order of student ID
allGrades := grades.InOrderTraversal()
for _, record := range allGrades {
    fmt.Printf("Student %d: %.1f\n", record.Key, record.Value)
}
```

### Example 3: Large Dataset Performance
```go
bt := treego.NewBTree[int, int](64)  // Higher degree for large datasets

// Insert 1 million records
for i := 0; i < 1000000; i++ {
    bt.Insert(i, i*i)
}

fmt.Printf("Tree with %d items has height %d\n", bt.Size(), bt.Height())
```

## Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| Insert    | O(log n)       | O(1)             |
| Search    | O(log n)       | O(1)             |
| Delete    | O(log n)       | O(1)             |
| Traversal | O(n)           | O(n)             |

## Choosing the Right Degree

- **Degree 2-4**: Good for small datasets or when minimizing memory usage
- **Degree 8-32**: Balanced choice for most use cases
- **Degree 64+**: Optimal for very large datasets, especially when stored on disk

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

