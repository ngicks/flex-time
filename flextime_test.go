package flextime_test

import (
	"testing"
	"time"

	"github.com/ngicks/flextime"
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

func TestFlextime(t *testing.T) {
	l, err := flextime.NewLayoutSet(`YYYY-MM-DD[THH[:mm[:ss.SSS]]][Z]`)
	require.NoError(t, err)
	p := flextime.NewFlextime(l)
	var parsed time.Time
	parsed, err = p.Parse("2022-10-20T23:16:22.168+09:00")
	require.NoError(t, err)
	require.Condition(t, func() (success bool) {
		return time.Date(2022, time.October, 20, 23, 16, 22, 168000000, jst).Equal(parsed)
	})
}
