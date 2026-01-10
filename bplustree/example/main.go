package main

import (
	"fmt"

	"github.com/l00pss/treego/bplustree"
)

func main() {
	tree := bplustree.New[int, string](3)

	fmt.Println("=== B+ Tree Example ===")
	fmt.Println("\nInserting values...")

	tree.Insert(10, "Value-10")
	tree.Insert(20, "Value-20")
	tree.Insert(5, "Value-5")
	tree.Insert(15, "Value-15")
	tree.Insert(25, "Value-25")
	tree.Insert(1, "Value-1")
	tree.Insert(30, "Value-30")
	tree.Insert(12, "Value-12")
	tree.Insert(18, "Value-18")

	fmt.Printf("Total entries: %d\n", tree.Len())

	fmt.Println("\n--- Search ---")
	if value, found := tree.Search(15); found {
		fmt.Printf("Key 15: %s\n", value)
	}

	if _, found := tree.Search(99); !found {
		fmt.Println("Key 99: not found")
	}

	fmt.Println("\n--- Range Query (10 to 25) ---")
	for _, e := range tree.Range(10, 25) {
		fmt.Printf("  Key: %d, Value: %s\n", e.Key, e.Value)
	}

	fmt.Println("\n--- Update ---")
	tree.Insert(10, "Updated-10")
	if value, found := tree.Search(10); found {
		fmt.Printf("Key 10 updated: %s\n", value)
	}

	fmt.Println("\n--- Delete ---")
	tree.Delete(5)
	fmt.Printf("After deleting key 5, total entries: %d\n", tree.Len())

	fmt.Println("\n--- All Entries (Sorted) ---")
	for _, e := range tree.All() {
		fmt.Printf("  Key: %d, Value: %s\n", e.Key, e.Value)
	}
}
