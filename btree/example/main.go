package main

import (
	"fmt"

	"github.com/l00pss/treego/btree"
)

func main() {
	// Example 1: B-tree with integer keys and string values
	fmt.Println("=== B-Tree with int keys and string values ===")
	btreeIntString := btree.NewBTree[int, string](3)

	// Insert some data
	data := map[int]string{
		10: "ten",
		20: "twenty",
		5:  "five",
		15: "fifteen",
		25: "twenty-five",
		1:  "one",
		30: "thirty",
	}

	for key, value := range data {
		btreeIntString.Insert(key, value)
	}

	fmt.Printf("Tree size: %d\n", btreeIntString.Size())
	fmt.Printf("Tree height: %d\n", btreeIntString.Height())

	// Search operations
	if value, found := btreeIntString.Search(15); found {
		fmt.Printf("Found key 15 with value: %s\n", value)
	}

	if _, found := btreeIntString.Search(100); !found {
		fmt.Println("Key 100 not found")
	}

	// In-order traversal
	fmt.Println("In-order traversal:")
	items := btreeIntString.InOrderTraversal()
	for _, item := range items {
		fmt.Printf("  %d -> %s\n", item.Key, item.Value)
	}

	// Example 2: B-tree with string keys and integer values
	fmt.Println("\n=== B-Tree with string keys and int values ===")
	btreeStringInt := btree.NewBTree[string, int](2)

	fruits := map[string]int{
		"apple":      5,
		"banana":     3,
		"cherry":     8,
		"date":       2,
		"elderberry": 1,
	}

	for fruit, count := range fruits {
		btreeStringInt.Insert(fruit, count)
	}

	fmt.Printf("Fruits inventory (sorted):\n")
	fruitItems := btreeStringInt.InOrderTraversal()
	for _, item := range fruitItems {
		fmt.Printf("  %s: %d\n", item.Key, item.Value)
	}

	// Example 3: Complex operations
	fmt.Println("\n=== Complex Operations ===")
	btree := btree.NewBTree[int, int](4)

	// Insert many values
	for i := 1; i <= 50; i++ {
		btree.Insert(i, i*i) // storing square values
	}

	fmt.Printf("Large tree size: %d\n", btree.Size())
	fmt.Printf("Large tree height: %d\n", btree.Height())

	// Delete some values
	fmt.Println("Deleting values 10, 25, 40...")
	btree.Delete(10)
	btree.Delete(25)
	btree.Delete(40)

	fmt.Printf("Tree size after deletions: %d\n", btree.Size())

	// Check if deleted values are gone
	for _, key := range []int{10, 25, 40} {
		if _, found := btree.Search(key); !found {
			fmt.Printf("Confirmed: key %d was deleted\n", key)
		}
	}

	// Print tree structure
	fmt.Println("\nTree structure:")
	fmt.Print(btreeStringInt.String())
}
