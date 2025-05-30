package tree

import "fmt"

func (tree *BPlusTree) PrettyPrint() {
	fmt.Println("B+ Tree Structure (Pretty Print):")
	prettyPrint(tree.Root, "", true)
}

func prettyPrint(node Node, prefix string, isTail bool) {
	branch := "├─ "
	if isTail {
		branch = "└─ "
	}

	if node.IsLeaf() {
		leaf := node.(*LeafNode)
		fmt.Printf("%s%sLeaf: [", prefix, branch)
		for i := 0; i < len(leaf.Keys); i++ {
			fmt.Printf("%d:%v", leaf.Keys[i], leaf.Values[i])
			if i < len(leaf.Keys)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")
	} else {
		internal := node.(*InternalNode)
		fmt.Printf("%s%sInternal: %v\n", prefix, branch, internal.Keys)

		for i := 0; i < len(internal.Children); i++ {
			isLast := i == len(internal.Children)-1
			newPrefix := prefix
			if isTail {
				newPrefix += "   "
			} else {
				newPrefix += "│  "
			}
			prettyPrint(internal.Children[i], newPrefix, isLast)
		}
	}
}
