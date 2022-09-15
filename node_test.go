package flextime

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
					value: []value{{typ: normal, value: "A"}},
				},
			},
			expected: []rawString{
				{},
				{{typ: normal, value: "A"}},
			},
		},
		{
			input: &treeNode{ // [A]B
				typ: text,
				left: &treeNode{
					typ:   optional,
					value: []value{{typ: normal, value: "A"}},
				},
				right: &treeNode{
					typ:   text,
					value: []value{{typ: normal, value: "B"}},
				},
			},
			expected: []rawString{
				{{typ: normal, value: "B"}},
				{{typ: normal, value: "A"}, {typ: normal, value: "B"}},
			},
		},
		{
			input: &treeNode{ // A[B]C
				typ:   text,
				value: []value{{typ: normal, value: "A"}},
				left: &treeNode{
					typ:   optional,
					value: []value{{typ: normal, value: "B"}},
				},
				right: &treeNode{
					typ:   text,
					value: []value{{typ: normal, value: "C"}},
				},
			},
			expected: []rawString{
				{{typ: normal, value: "A"}, {typ: normal, value: "C"}},
				{{typ: normal, value: "A"}, {typ: normal, value: "B"}, {typ: normal, value: "C"}},
			},
		},
		{
			input: &treeNode{ // [A[B]C]
				typ: text,
				left: &treeNode{
					typ:   optional,
					value: []value{{typ: normal, value: "A"}},
					left: &treeNode{
						typ:   optional,
						value: []value{{typ: normal, value: "B"}},
					},
					right: &treeNode{
						typ:   text,
						value: []value{{typ: normal, value: "C"}},
					},
				},
			},
			expected: []rawString{
				{},
				{{typ: normal, value: "A"}, {typ: normal, value: "C"}},
				{{typ: normal, value: "A"}, {typ: normal, value: "B"}, {typ: normal, value: "C"}},
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
