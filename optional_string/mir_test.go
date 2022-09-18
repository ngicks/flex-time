package optionalstring

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nodeTestCase struct {
	input    *treeNode
	expected []RawString
}

func TestNode(t *testing.T) {
	cases := []nodeTestCase{
		{
			input: &treeNode{ // [A]
				typ: nonOptional,
				left: &treeNode{
					typ:   optional,
					value: []TextNode{{typ: Normal, value: "A"}},
				},
			},
			expected: []RawString{
				{},
				{{typ: Normal, value: "A"}},
			},
		},
		{
			input: &treeNode{ // [A]B
				typ: nonOptional,
				left: &treeNode{
					typ:   optional,
					value: []TextNode{{typ: Normal, value: "A"}},
				},
				right: &treeNode{
					typ:   nonOptional,
					value: []TextNode{{typ: Normal, value: "B"}},
				},
			},
			expected: []RawString{
				{{typ: Normal, value: "B"}},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "B"}},
			},
		},
		{
			input: &treeNode{ // A[B]C
				typ:   nonOptional,
				value: []TextNode{{typ: Normal, value: "A"}},
				left: &treeNode{
					typ:   optional,
					value: []TextNode{{typ: Normal, value: "B"}},
				},
				right: &treeNode{
					typ:   nonOptional,
					value: []TextNode{{typ: Normal, value: "C"}},
				},
			},
			expected: []RawString{
				{{typ: Normal, value: "A"}, {typ: Normal, value: "C"}},
				{{typ: Normal, value: "A"}, {typ: Normal, value: "B"}, {typ: Normal, value: "C"}},
			},
		},
		{
			input: &treeNode{ // [A[B]C]
				typ: nonOptional,
				left: &treeNode{
					typ:   optional,
					value: []TextNode{{typ: Normal, value: "A"}},
					left: &treeNode{
						typ:   optional,
						value: []TextNode{{typ: Normal, value: "B"}},
					},
					right: &treeNode{
						typ:   nonOptional,
						value: []TextNode{{typ: Normal, value: "C"}},
					},
				},
			},
			expected: []RawString{
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
