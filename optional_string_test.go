package flextime_test

import (
	"sort"
	"testing"

	"github.com/ngicks/flextime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYYYMMDDOptional(t *testing.T) {
	input := `YYYY-MM-DD[THH[:mm[:ss.SSS]]][Z]`
	output := []string{
		`YYYY-MM-DDTHH:mm:ss.SSSZ`,
		`YYYY-MM-DDTHH:mm:ss.SSS`,
		`YYYY-MM-DDTHH:mmZ`,
		`YYYY-MM-DDTHH:mm`,
		`YYYY-MM-DDTHHZ`,
		`YYYY-MM-DDTHH`,
		`YYYY-MM-DDZ`,
		`YYYY-MM-DD`,
	}

	result, err := flextime.EnumerateOptionalString(input)
	require.NoError(t, err)
	sort.Strings(result)
	sort.Strings(output)
	assert.Equal(t, output, result)
}
