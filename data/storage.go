package data

import (
	"encoding/gob"
	"os"

	"github.com/Prayag2003/b-plus-tree-database-indexing/tree"
)

func SaveTree(filePath string, bpt *tree.BPlusTree) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	return enc.Encode(bpt)
}

func LoadTree(filePath string) (*tree.BPlusTree, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bpt *tree.BPlusTree
	dec := gob.NewDecoder(file)
	err = dec.Decode(&bpt)
	return bpt, err
}
