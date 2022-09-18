package optionalstring

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nodeTestCase struct {
	input    *treeNode
	expected []rawString
}

func TestNode(t *testing.T) {
	cases := []nodeTestCase{
		{
			input: &treeNode{ // [A]
				typ: text,
				left: &treeNode{
					typ:   optional,
					value: []Value{{typ: Normal, value: "A"}},
				},
			},
			expected: []rawString{
				{},
				{{typ: Normal, value: "A"}},
			},
		},
		{
			input: &treeNode{ // [A]B
				typ: text,
				left: &treeNode{
					typ:   optional,
					value: []Value{{typ: Normal, value: "A"}},
				},
				right: &treeNode{
					typ:   text,
					value: []Value{{typ: Normal, value: "B"}},
				},
			},
			expected: []rawString{
				{{typ: Normal, value: "B"}},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "B"}},
			},
		},
		{
			input: &treeNode{ // A[B]C
				typ:   text,
				value: []Value{{typ: Normal, value: "A"}},
				left: &treeNode{
					typ:   optional,
					value: []Value{{typ: Normal, value: "B"}},
				},
				right: &treeNode{
					typ:   text,
					value: []Value{{typ: Normal, value: "C"}},
				},
			},
			expected: []rawString{
				{{typ: Normal, value: "A"}, {typ: Normal, value: "C"}},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "B"}, {typ: Normal, value: "C"}},
			},
		},
		{
			input: &treeNode{ // [A[B]C]
				typ: text,
				left: &treeNode{
					typ:   optional,
					value: []Value{{typ: Normal, value: "A"}},
					left: &treeNode{
						typ:   optional,
						value: []Value{{typ: Normal, value: "B"}},
					},
					right: &treeNode{
						typ:   text,
						value: []Value{{typ: Normal, value: "C"}},
					},
				},
			},
			expected: []rawString{
				{},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "C"}},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "B"}, {typ: Normal, value: "C"}},
			},
		},
	}

	for _, tc := range cases {
		flatten := tc.input.Flatten()

		sort.Slice(tc.expected, func(i, j int) bool {
			return tc.expected[i].String() < tc.expected[j].String()
		})
		sort.Slice(flatten, func(i, j int) bool {
			return flatten[i].String() < flatten[j].String()
		})

		assert.Equal(t, tc.expected, flatten)
	}

}
