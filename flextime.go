package flextime

import (
	"time"
)

type Option func(f *Flextime)

type Flextime struct {
	layouts *LayoutSet
}

func NewFlextime(layouts *LayoutSet) *Flextime {
	return &Flextime{
		layouts: layouts,
	}
}

func (f *Flextime) parse(value string, parser func(layout, value string) (time.Time, error)) (time.Time, error) {
	var lastErr error
	for _, layout := range f.layouts.Layout() {
		t, err := parser(layout, value)
		if err != nil {
			lastErr = err
		} else {
			return t, nil
		}
	}
	return time.Time{}, lastErr
}

func (f *Flextime) Parse(value string) (time.Time, error) {
	return f.parse(
		value,
		func(layout, value string) (time.Time, error) { return time.Parse(layout, value) },
	)
}

func (f *Flextime) ParseInLocation(value string, loc *time.Location) (time.Time, error) {
	return f.parse(
		value,
		func(
			layout, value string,
		) (time.Time, error) {
			return time.ParseInLocation(layout, value, loc)
		},
	)
}

func (p *Flextime) LayoutSet() *LayoutSet {
	return p.layouts
}

func (p *Flextime) AddLayout(other *LayoutSet) *Flextime {
	return &Flextime{
		layouts: p.layouts.AddLayout(other),
	}
}
