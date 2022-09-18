package optionalstring

import (
	"github.com/ngicks/type-param-common/iterator"
)

type treeNodeType int

const (
	nonOptional treeNodeType = iota
	optional
)

// treeNode is node of optional string tree.
// It is seperated by optional part. left node is always optional.
// if lower parts have no optional part there must not be child node
type treeNode struct {
	left  *treeNode
	right *treeNode
	value []TextNode
	typ   treeNodeType
}

func (n *treeNode) Clone() []TextNode {
	if n.value == nil {
		return nil
	}
	cloned := make([]TextNode, len(n.value))
	copy(cloned, n.value)
	return cloned
}

func (n *treeNode) AddValue(v string, typ valueType) {
	n.value = append(n.value, TextNode{value: v, typ: typ})
}

func (n *treeNode) SetAsOptional() {
	n.typ = optional
}

func (n *treeNode) IsOptional() bool {
	return n.typ == optional
}

func (n *treeNode) Left() *treeNode {
	if n.left == nil {
		n.left = &treeNode{}
	}
	return n.left
}
func (n *treeNode) HasLeft() bool {
	return n.left != nil
}

func (n *treeNode) Right() *treeNode {
	if n.right == nil {
		n.right = &treeNode{}
	}
	return n.right
}
func (n *treeNode) HasRight() bool {
	return n.right != nil
}

func (n *treeNode) Flatten() []RawString {
	return n.flatten()
}
func (n *treeNode) flatten() []RawString {
	// root node must not be optional

	// treeNodes is value of self -> left -> right order.
	var cur RawString
	var total []RawString
	if c := n.Clone(); len(c) > 0 {
		cur = RawString(c).Clone()
	} else {
		cur = NewRawString()
	}
	total = []RawString{cur.Clone()}

	if n.HasLeft() {
		l := n.Left()
		totalCloned := iterator.
			FromSlice(total).
			Collect()
		total = total[:0]
		for _, s := range l.flatten() {
			for _, str := range totalCloned {
				total = append(total, str.Append(s))
			}
		}
		if l.IsOptional() {
			total = append(total, cur)
		}
	}

	if n.HasRight() {
		// right cannot be optional.

		totalCloned := iterator.
			FromSlice(total).
			Collect()
		total = total[:0]

		for _, s := range n.Right().flatten() {
			for _, str := range totalCloned {
				total = append(total, str.Append(s))
			}
		}
	}

	return total
}
