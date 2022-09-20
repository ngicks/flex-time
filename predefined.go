package flextime

import "time"

// RFC3339Optinal is LayoutSet where year, month, date is mandatory.
// And lower parts (hours, minutes, seconds, nanoseconds) and timezone offset are optional.
var RFC3339Optinal *LayoutSet = Must(
	func() (*LayoutSet, error) { return NewLayoutSet(`YYYY-MM-DD[THH[:mm[:ss.999999999]]][Z]`) },
)

var RFC3339orUnixMilli *CombinedFlextime = NewCombined([]*Flextime{NewFlextime(RFC3339Optinal)}, time.UnixMilli)

func Must[T any](newFunc func() (T, error)) T {
	v, err := newFunc()
	if err != nil {
		panic(err)
	}
	return v
}
