package tree

import (
	"sort"
)

type BPlusTree struct {
	Root  Node
	Order int
}

func NewBPlusTree(order int) *BPlusTree {
	if order < 3 {
		panic("Order must be >= 3")
	}

	return &BPlusTree{
		Root:  &LeafNode{},
		Order: order,
	}
}

func (tree *BPlusTree) Insert(key int, value any) error {
	newRoot, split := insert(tree.Root, key, value, tree.Order)
	tree.Root = newRoot
	_ = split
	return nil
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

func insert(node Node, key int, value any, order int) (Node, bool) {
	if node.IsLeaf() {
		leaf := node.(*LeafNode)
		insertIntoLeaf(leaf, key, value)

		if len(leaf.Keys) > order-1 {
			left, right, promotedKey := splitLeafNode(leaf, order)
			newInternal := &InternalNode{
				Keys:     []int{promotedKey},
				Children: []Node{left, right},
			}
			return newInternal, true
		}
		return leaf, false
	}

	internal := node.(*InternalNode)
	i := sort.SearchInts(internal.Keys, key)
	if i < len(internal.Keys) && key >= internal.Keys[i] {
		i++
	}

	child := internal.Children[i]
	newChild, childSplit := insert(child, key, value, order)

	if !childSplit {
		internal.Children[i] = newChild
		return internal, false
	}

	// Promote from newChild
	newInternal := newChild.(*InternalNode)
	promotedKey := newInternal.Keys[0]
	left := newInternal.Children[0]
	right := newInternal.Children[1]

	internal.Keys = append(internal.Keys, 0)
	copy(internal.Keys[i+1:], internal.Keys[i:])
	internal.Keys[i] = promotedKey

	internal.Children[i] = left
	internal.Children = append(internal.Children, nil)
	copy(internal.Children[i+2:], internal.Children[i+1:])
	internal.Children[i+1] = right

	if len(internal.Keys) > order-1 {
		return splitInternalNode(internal)
	}
	return internal, false
}

func insertIntoLeaf(leaf *LeafNode, key int, value any) {
	i := sort.SearchInts(leaf.Keys, key)

	if i < len(leaf.Keys) && leaf.Keys[i] == key {
		leaf.Values[i] = value
		return
	}

	leaf.Keys = append(leaf.Keys, 0)
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

// splitLeafNode splits a leaf into two and returns the left and right nodes and promoted key.
func splitLeafNode(leaf *LeafNode, order int) (*LeafNode, *LeafNode, int) {
	mid := order / 2

	left := &LeafNode{
		Keys:   append([]int{}, leaf.Keys[:mid]...),
		Values: append([]any{}, leaf.Values[:mid]...),
	}
	right := &LeafNode{
		Keys:   append([]int{}, leaf.Keys[mid:]...),
		Values: append([]any{}, leaf.Values[mid:]...),
		Next:   leaf.Next,
	}
	left.Next = right

	return left, right, right.Keys[0]
}

func splitInternalNode(node *InternalNode) (*InternalNode, bool) {
	mid := len(node.Keys) / 2
	promotedKey := node.Keys[mid]

	left := &InternalNode{
		Keys:     append([]int{}, node.Keys[:mid]...),
		Children: append([]Node{}, node.Children[:mid+1]...),
	}

	right := &InternalNode{
		Keys:     append([]int{}, node.Keys[mid+1:]...),
		Children: append([]Node{}, node.Children[mid+1:]...),
	}

	newRoot := &InternalNode{
		Keys:     []int{promotedKey},
		Children: []Node{left, right},
	}
	return newRoot, true
}

func (tree *BPlusTree) Search(key int) (any, bool) {
	node := tree.Root
	for !node.IsLeaf() {
		internal := node.(*InternalNode)
		i := sort.SearchInts(internal.Keys, key)
		if i < len(internal.Keys) && key >= internal.Keys[i] {
			i++
		}
		node = internal.Children[i]
	}

	leaf := node.(*LeafNode)
	i := sort.SearchInts(leaf.Keys, key)
	if i < len(leaf.Keys) && leaf.Keys[i] == key {
		return leaf.Values[i], true
	}
	return nil, false
}

func (tree *BPlusTree) RangeSearch(start, end int) []any {
	var results []any
	node := tree.Root

	for !node.IsLeaf() {
		internal := node.(*InternalNode)
		i := sort.SearchInts(internal.Keys, start)
		if i >= len(internal.Children) {
			i = len(internal.Children) - 1
		}
		node = internal.Children[i]
	}

	for node != nil {
		leaf := node.(*LeafNode)
		for i, key := range leaf.Keys {
			if key >= start && key <= end {
				results = append(results, leaf.Values[i])
			} else if key > end {
				return results
			}
		}
		node = leaf.Next
	}

	return results
}

func (tree *BPlusTree) Delete(key int) {
	tree.Root = deleteKey(tree.Root, key)
}

func deleteKey(node Node, key int) Node {
	if node.IsLeaf() {
		leaf := node.(*LeafNode)
		i := sort.SearchInts(leaf.Keys, key)
		if i < len(leaf.Keys) && leaf.Keys[i] == key {
			leaf.Keys = append(leaf.Keys[:i], leaf.Keys[i+1:]...)
			leaf.Values = append(leaf.Values[:i], leaf.Values[i+1:]...)
		}
		return leaf
	}

	internal := node.(*InternalNode)
	i := sort.SearchInts(internal.Keys, key)
	if i >= len(internal.Children) {
		i = len(internal.Children) - 1
	}
	child := internal.Children[i]
	internal.Children[i] = deleteKey(child, key)

	return internal
}
