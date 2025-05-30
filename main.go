package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/Prayag2003/b-plus-tree-database-indexing/config"
	"github.com/Prayag2003/b-plus-tree-database-indexing/data"
	"github.com/Prayag2003/b-plus-tree-database-indexing/tree"
)

func init() {
	gob.Register(&tree.InternalNode{})
	gob.Register(&tree.LeafNode{})
}

func main() {
	// Set up logger with timestamp and file info
	logger := log.New(os.Stdout, "BPlusTreeDB: ", log.Ldate|log.Ltime|log.Lshortfile)

	cfg := config.LoadConfig()

	bpt, err := data.LoadTree(cfg.StoragePath)
	if err != nil {
		logger.Printf("Failed to load tree from disk (%s): %v. Creating new tree and inserting initial data.", cfg.StoragePath, err)
		bpt = tree.NewBPlusTree(cfg.TreeOrder)
		for k, v := range data.Emails {
			bpt.Insert(k, v)
		}
		logger.Println("Inserted initial data into B+ Tree")
	} else {
		logger.Println("Loaded B+ Tree from disk")
	}

	val, found := bpt.Search(1003)
	if found {
		logger.Printf("Found key 1003 with value: %v", val)
	} else {
		logger.Println("Key 1003 not found")
	}

	logger.Println("B+ Tree Structure:")
	bpt.PrettyPrint()

	if err := data.SaveTree(cfg.StoragePath, bpt); err != nil {
		logger.Printf("Failed to save tree: %v", err)
	} else {
		logger.Println("Saved B+ Tree to disk successfully")
	}
}
