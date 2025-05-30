package main

import (
	"fmt"

	"github.com/Prayag2003/b-plus-tree-database-indexing/tree"
)

func main() {

	tree := tree.NewBPlusTree(4)

	tree.Insert(10, "A")
	tree.Insert(20, "B")
	tree.Insert(30, "C")
	tree.Insert(25, "Z")
	tree.Insert(5, "E")
	tree.Insert(15, "F")
	tree.Insert(35, "G")
	tree.Insert(40, "H")
	tree.Insert(50, "I")

	fmt.Println("Inserted keys into B+Tree.")

	val, found := tree.Search(35)
	if found {
		fmt.Printf("Found key 35 with value: %v\n", val)
	} else {
		fmt.Println("Key 35 not found.")
	}
}
