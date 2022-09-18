package flextime

import "testing"

type simpleCase[T any] struct {
	input    T
	expected T
}

func TestGetUntilClosingSingleSquote(t *testing.T) {
	cases := []simpleCase[string]{
		{
			input:    `aaaa'`,
			expected: `aaaa`,
		},
		{
			input:    `aaaa\'aaa'`,
			expected: `aaaa\'aaa`,
		},
		{
			input:    `aaaa\'aaa'`,
			expected: `aaaa\'aaa`,
		},
		{
			input:    `aa\\'`,
			expected: `aa\\`,
		},
	}

	for _, testCase := range cases {
		result := getUntilClosingSingleSquote(testCase.input)

		if testCase.expected != result {
			t.Errorf("not equal. expected = %s, actual = %s", testCase.expected, result)
		}
	}
}
