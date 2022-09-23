package flextime_test

import (
	"testing"
	_ "time/tzdata"

	"github.com/ngicks/flextime"
	"github.com/stretchr/testify/assert"
)

type replaceTimeTokenTestCase struct {
	input    string
	expected string
}

func TestReplaceTimeToken(t *testing.T) {
	cases := []replaceTimeTokenTestCase{
		{
			input:    "yyyy-MM-ddTHH:mm:ss.SSSSSSSSSZ",
			expected: "2006-01-02T15:04:05.000000000Z07:00",
		},
		{
			input:    "YYYY-MM-DDTHH:mm:ss.999999999Z",
			expected: "2006-01-02T15:04:05.999999999Z07:00",
		},
		{
			input:    `YYYY-MM-DD\[T\]HH:mm:ss`,
			expected: `2006-01-02[T]15:04:05`,
		},
		{
			input:    `YYYY-MM-DD'T'HH:mm:ss`,
			expected: `2006-01-02T15:04:05`,
		},
		{
			input:    `xxxx-'Www'-e`,
			expected: `xxxx-Www-e`,
		},
	}

	for _, testCase := range cases {
		out, err := flextime.ReplaceTimeToken(testCase.input)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, out)
	}
}
