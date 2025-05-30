package main

import (
	"fmt"

	"github.com/Prayag2003/b-plus-tree-database-indexing/tree"
)

func main() {
	bpt := tree.NewBPlusTree(4)

	bpt.Insert(10, "A")
	bpt.Insert(20, "B")
	bpt.Insert(30, "C")
	bpt.Insert(25, "Z")

	fmt.Println("Inserted keys into B+Tree.")
}
