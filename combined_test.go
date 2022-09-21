package flextime_test

import (
	"math"
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
			input:    uint64(1666282966123),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			// 1666282966123
			input:    []byte("1666282966123"),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			input:    []byte(`"2022-10-20T16:22:46.123+09:00"`),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, jst),
		},
		{
			input:    []byte(`"2022-10-20T16:22:46.123"`),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 123000000, time.UTC),
		},
		{
			input:    []byte(`"2022-10-20T16:22:46+09:00"`),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, jst),
		},
		{
			input:    []byte(`"2022-10-20T16:22:46"`),
			expected: time.Date(2022, 10, 20, 16, 22, 46, 0, time.UTC),
		},
		{
			input:    []byte(`"2022-10-20T16:22+09:00"`),
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, jst),
		},
		{
			input:    []byte(`"2022-10-20T16:22"`),
			expected: time.Date(2022, 10, 20, 16, 22, 0, 0, time.UTC),
		},
		{
			input:    []byte(`"2022-10-20T16+09:00"`),
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, jst),
		},
		{
			input:    []byte(`"2022-10-20T16"`),
			expected: time.Date(2022, 10, 20, 16, 0, 0, 0, time.UTC),
		},
		{
			input:    []byte(`"2022-10-20+09:00"`),
			expected: time.Date(2022, 10, 20, 0, 0, 0, 0, jst),
		},
		{
			input:    []byte(`"2022-10-20"`),
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
			func() (success bool) {
				return expectedInLocal.Equal(parsed) && expectedInLocal.Format("Z07:00") == parsed.Format("Z07:00")
			},
			cmp.Diff(expectedInLocal, parsed),
		)
	}
}

func TestCombinedError(t *testing.T) {
	var err error

	var parseError *time.ParseError
	for _, invalid := range []any{"2022-10-20T16:2a2:46.123+09:00", "22-10-20T16:22:46.123+09:00"} {
		_, err = flextime.RFC3339orUnixMilli.ParseInLocation(invalid, nil)
		assert.ErrorAs(t, err, &parseError)
	}

	var unsupportedTypeError *flextime.UnsupportedTypeError
	for _, invalid := range []any{
		false,
		true,
		struct{}{},
		[3]int{1, 2, 3},
		[...]byte{34, 50, 48, 50, 50, 45, 49, 48, 45, 50, 48, 34},
	} {
		_, err = flextime.RFC3339orUnixMilli.ParseInLocation(invalid, nil)
		assert.ErrorAs(t, err, &unsupportedTypeError)
	}

	var unmarshalError *flextime.UnmarshalError
	for _, nonUnmarshalable := range []any{[]byte("foobar"), []byte("123q")} {
		_, err = flextime.RFC3339orUnixMilli.ParseInLocation(nonUnmarshalable, nil)
		assert.ErrorAs(t, err, &unmarshalError)
	}

	var valueOutOfRangeError *flextime.ValueOutOfRangeError
	for i := 1; i < 100; i++ {
		_, err = flextime.RFC3339orUnixMilli.ParseInLocation(uint64(math.MaxInt64)+uint64(i), nil)
		assert.ErrorAs(t, err, &valueOutOfRangeError)
	}

	p := flextime.NewCombined(
		[]*flextime.Flextime{flextime.NewFlextime(flextime.RFC3339Optinal)},
		nil,
	)

	_, err = p.Parse(1666282966123)
	assert.ErrorIs(t, err, flextime.ErrEmptyNumParser)
}
