package tree

type Node interface {
	IsLeaf() bool
	GetKeys() []int
}

type InternalNode struct {
	Keys     []int
	Children []Node
}

type LeafNode struct {
	Keys   []int
	Values []any
	Next   *LeafNode
}

func (in *InternalNode) IsLeaf() bool {
	return false
}

func (in *InternalNode) GetKeys() []int {
	return in.Keys
}

func (ln *LeafNode) IsLeaf() bool {
	return true
}

func (ln *LeafNode) GetKeys() []int {
	return ln.Keys
}
