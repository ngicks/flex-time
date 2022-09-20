package flextime_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngicks/flextime"
	"github.com/stretchr/testify/assert"
)

type combinedTestCase struct {
	input    any
	expected time.Time
}

func TestCombined(t *testing.T) {
	cases := []combinedTestCase{
		{
			input:    1666282966123,
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			// 1666282966123
			input:    []byte{49, 54, 54, 54, 50, 56, 50, 57, 54, 54, 49, 50, 51},
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			//"2022-10-20T16:22:46.123Z"
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 58, 52, 54, 46, 49, 50, 51, 90, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 58, 52, 54, 46, 49, 50, 51, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 58, 52, 54, 90, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 58, 52, 54, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 90, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 58, 50, 50, 34},
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 90, 34},
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 84, 49, 54, 34},
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, time.UTC),
		},
		{
			input:    []byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 34},
			expected: time.Date(2022, 10, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "2022-10-20T16:22:46.123+09:00",
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, jst),
		},
		{
			input:    "2022-10-20T16:22:46.123",
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			input:    "2022-10-20T16:22:46+09:00",
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, jst),
		},
		{
			input:    "2022-10-20T16:22:46",
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, time.UTC),
		},
		{
			input:    "2022-10-20T16:22+09:00",
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, jst),
		},
		{
			input:    "2022-10-20T16:22",
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, time.UTC),
		},
		{
			input:    "2022-10-20T16+09:00",
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, jst),
		},
		{
			input:    "2022-10-20T16",
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, time.UTC),
		},
		{
			input:    "2022-10-20+09:00",
			expected: time.Date(2022, 10, 20, 0, 0, 0, 0, jst),
		},
		{
			input:    "2022-10-20",
			expected: time.Date(2022, 10, 20, 0, 0, 0, 0, time.UTC),
		},
	}
	expected := time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC)

	t.Logf("%d", expected.UnixMilli())

	for _, testCase := range cases {
		parsed, err := flextime.RFC3339orUnixMilli.ParseInLocation(testCase.input, time.UTC)
		assert.NoError(t, err)
		assert.Conditionf(
			t,
			func() (success bool) { return testCase.expected.Equal(parsed) },
			cmp.Diff(testCase.expected, parsed),
		)
		parsed, err = flextime.RFC3339orUnixMilli.Parse(testCase.input)

		var expectedInLocal time.Time
		if testCase.expected.Location() == time.UTC {
			expectedInLocal = testCase.expected.In(time.Local)
		} else {
			expectedInLocal = testCase.expected
		}
		assert.NoError(t, err)
		assert.Conditionf(
			t,
			func() (success bool) { return expectedInLocal.Equal(parsed) },
			cmp.Diff(expectedInLocal, parsed),
		)
	}
}
