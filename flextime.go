package flextime

import (
	"time"
)

type Option func(f *Flextime)

func TryAllLayouts(f *Flextime) {
	f.tryAllLayouts = true
}

type Flextime struct {
	tryAllLayouts bool
	layouts       *LayoutSet
}

func NewMultiLayoutParser(layouts *LayoutSet, options ...Option) *Flextime {
	return &Flextime{
		layouts: layouts,
	}
}

func Compile(optionalStr string, options ...Option) (*Flextime, error) {
	layouts, err := NewLayoutSet(optionalStr)
	if err != nil {
		return nil, err
	}
	return NewMultiLayoutParser(layouts, options...), nil
}

func (p *Flextime) parse(value string, parser func(layout, value string) (time.Time, error)) (time.Time, error) {
	var lastErr error
	for _, layout := range p.layouts.Layout() {
		t, err := parser(layout, value)
		if err != nil {
			lastErr = err
		} else {
			return t, nil
		}
	}
	return time.Time{}, lastErr
}

func (p *Flextime) Parse(value string) (time.Time, error) {
	return p.parse(
		value,
		func(layout, value string) (time.Time, error) { return time.Parse(layout, value) },
	)
}

func (p *Flextime) ParseInLocation(value string, loc *time.Location) (time.Time, error) {
	return p.parse(
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
		layouts:       p.layouts.AddLayout(other),
		tryAllLayouts: p.tryAllLayouts,
	}
}
