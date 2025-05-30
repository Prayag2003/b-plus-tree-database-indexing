package tree

import (
	"errors"
	"sort"
)

type BPlusTree struct {
	root  Node
	order int
}

func NewBPlusTree(order int) *BPlusTree {
	if order < 3 {
		panic("Order must be >= 3")
	}

	return &BPlusTree{
		root:  &LeafNode{},
		order: order,
	}
}

func (tree *BPlusTree) Insert(key int, value any) error {
	root := tree.root

	if root.IsLeaf() {
		leaf := root.(*LeafNode)
		insertIntoLeaf(leaf, key, value, tree.order)

		if len(leaf.Keys) > tree.order-1 {
			left, right, promotedKey := splitLeafNode(leaf, tree.order)

			tree.root = &InternalNode{
				Keys:     []int{promotedKey},
				Children: []Node{left, right},
			}
		}
		return nil
	}

	// Later we'll support internal node insertion here
	return errors.New("internal node insert not implemented yet")
}

/*
	Example: Inserting key=25, value="Z" into a leaf node

	Assume:
	- Tree order = 4 (max 3 keys before split)
	- Current leaf state:
		leaf.Keys   = [10, 20, 30]
		leaf.Values = ["A", "B", "C"]

	Steps:
	1. Find correct index to insert:
	- sort.SearchInts([10, 20, 30], 25) → returns index 2

	2. Check for existing key at index 2:
	- leaf.Keys[2] = 30 ≠ 25 → continue with insert

	3. Make space in keys slice:
	- append dummy value:         [10, 20, 30, 0]
	- shift right using copy:     [10, 20, 30, 30]
	- insert key at index 2:      [10, 20, 25, 30]

	4. Do the same for values slice:
	- append nil:                 ["A", "B", "C", nil]
	- shift right using copy:     ["A", "B", "C", "C"]
	- insert value at index 2:    ["A", "B", "Z", "C"]

	Final result:
		leaf.Keys   = [10, 20, 25, 30]
		leaf.Values = ["A", "B", "Z", "C"]

	Note:
	- This causes a leaf overflow (4 keys in a tree of order 4).
	- Next step: split the leaf and promote middle key to parent.
*/

func insertIntoLeaf(leaf *LeafNode, key int, value any, order int) {
	/*
		sort.SearchInts() uses binary search.
		It returns the first index where leaf.Keys[i] >= key
	*/
	i := sort.SearchInts(leaf.Keys, key)

	// If the key already exists at index i, overwrite the existing value.
	if i < len(leaf.Keys) && leaf.Keys[i] == key {
		leaf.Values[i] = value
		return
	}

	leaf.Keys = append(leaf.Keys, 0)

	// Shift all elements from i onward to the right by 1 position to make space
	copy(leaf.Keys[i+1:], leaf.Keys[i:])

	leaf.Keys[i] = key

	leaf.Values = append(leaf.Values, nil)
	copy(leaf.Values[i+1:], leaf.Values[i:])
	leaf.Values[i] = value

}

/*
	What Happens When a Leaf Node Overflows?
	When a leaf node reaches order number of keys and you try to insert another one:

	We must:
		1. Split the leaf into two leaf nodes.
		2. Move half the keys/values to the new right sibling.
		3. Update sibling pointers (for sequential leaf traversal).
		4. Promote the first key of the right sibling to the parent internal node.
		5. If there’s no parent (root was a leaf), create a new root internal node.
*/

func splitLeafNode(leaf *LeafNode, order int) (left *LeafNode, right *LeafNode, promotedKey int) {
	mid := order / 2

	right = &LeafNode{
		Keys:   append([]int(nil), leaf.Keys[mid:]...),
		Values: append([]any(nil), leaf.Values[mid:]...),
		Next:   leaf.Next,
	}

	left = &LeafNode{
		Keys:   append([]int(nil), leaf.Keys[:mid]...),
		Values: append([]any(nil), leaf.Values[:mid]...),
		Next:   right,
	}

	promotedKey = right.Keys[0]
	return left, right, promotedKey
}
