package optionalstring

import (
	"fmt"

	"github.com/ngicks/type-param-common/iterator"
	parsec "github.com/prataprc/goparsec"
)

type treeNodeType int

const (
	text treeNodeType = iota
	optional
)

type treeNode struct {
	left  *treeNode
	right *treeNode
	value []Value
	typ   treeNodeType
}

func (n *treeNode) Clone() []Value {
	if n.value == nil {
		return nil
	}
	cloned := make([]Value, len(n.value))
	copy(cloned, n.value)
	return cloned
}

func (n *treeNode) AddValue(v string, typ valueType) {
	n.value = append(n.value, Value{value: v, typ: typ})
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

func (n *treeNode) Flatten() []rawString {
	return n.flatten()
}
func (n *treeNode) flatten() []rawString {
	// root node must not be optional

	// treeNodes is value of self -> left -> right order.
	var cur rawString
	var total []rawString
	if c := n.Clone(); len(c) > 0 {
		cur = rawString(c).Clone()
	} else {
		cur = newRawString()
	}
	total = []rawString{cur.Clone()}

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

func decode(node parsec.Queryable) *treeNode {
	root := &treeNode{}
	recursiveDecode(node.GetChildren(), root)
	return root
}

func recursiveDecode(nodes []parsec.Queryable, ctx *treeNode) {
	var onceFound bool

	for i := 0; i < len(nodes); i++ {
		if onceFound {
			recursiveDecode(nodes[i:], ctx.Right())
			return
		}

		switch nodes[i].GetName() {
		case OPTIONALSTRING:
			// skipping first node.
			recursiveDecode(nodes[i].GetChildren(), ctx)
		case OPTIONAL:
			var optNext *treeNode
			if !onceFound {
				onceFound = true
				optNext = ctx.Left()
			} else {
				panic(
					fmt.Sprintf(
						"incorrect implementation: %s, %s",
						nodes[i].GetName(),
						nodes[i].GetValue(),
					),
				)
			}
			optNext.SetAsOptional()
			recursiveDecode(nodes[i].GetChildren(), optNext)
		case CHARS:
			for _, v := range nodes[i].GetChildren() {
				switch v.GetName() {
				case NORMALCHARS:
					ctx.AddValue(v.GetValue(), Normal)
				case ESCAPEDCHAR:
					ctx.AddValue(v.GetValue(), SingleQuoteEscaped)
				default:
					panic(fmt.Sprintf("incorrect implementation: %s, %s", v.GetName(), v.GetValue()))
				}
			}
		case ESCAPED:
			ctx.AddValue(nodes[i].GetValue(), SingleQuoteEscaped)
		case ITEMS:
			recursiveDecode(nodes[i].GetChildren(), ctx)
		}
	}
}
