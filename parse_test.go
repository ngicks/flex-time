package flextime_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/ngicks/flextime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	jst *time.Location
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
}

func TestParser(t *testing.T) {
	p, err := flextime.Compile(`YYYY-MM-DD[THH[:mm[:ss.SSS]]][Z]`)
	require.NoError(t, err)
	var parsed time.Time
	parsed, err = p.Parse("2022-10-20T23:16:22.168+09:00")
	require.NoError(t, err)
	require.Condition(t, func() (success bool) {
		return time.Date(2022, time.October, 20, 23, 16, 22, 168000000, jst).Equal(parsed)
	})

}

type replaceTimeTokenTestCase struct {
	input    string
	expected string
}

func TestReplaceTimeToken(t *testing.T) {
	cases := []replaceTimeTokenTestCase{
		{
			input:    "yyyy-MM-ddTHH:mm:ss.SSSSSSSSSZ07:00",
			expected: "2006-01-02T15:04:05.000000000Z07:00",
		},
		{
			input:    "YYYY-MM-DDTHH:mm:ss.999999999Z07:00",
			expected: "2006-01-02T15:04:05.999999999Z07:00",
		},
		{
			input:    `YYYY-MM-DD\[T\]HH:mm:ss`,
			expected: `2006-01-02[T]15:04:05`,
		},
	}

	for _, testCase := range cases {
		out, err := flextime.ReplaceTimeToken(testCase.input)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, out)
	}

}
