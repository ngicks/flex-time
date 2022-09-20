package flextime

import (
	"sort"
	"strings"

	optionalstring "github.com/ngicks/flextime/optional_string"
	"github.com/ngicks/type-param-common/set"
)

type LayoutSet struct {
	layouts []string
}

func newLayoutSet(layouts []string) *LayoutSet {
	sort.Slice(layouts, func(i, j int) bool {
		iLen := len(layouts[i])
		jLen := len(layouts[j])
		if iLen != jLen {
			return iLen > jLen
		} else {
			return strings.Compare(layouts[i], layouts[j]) == -1
		}
	})

	return &LayoutSet{
		layouts: layouts,
	}
}

func NewLayoutSet(optionalStr string) (*LayoutSet, error) {
	rawFormats, err := optionalstring.EnumerateOptionalStringRaw(optionalStr)
	if err != nil {
		return nil, err
	}

	layouts := make([]string, len(rawFormats))
	for i := 0; i < len(rawFormats); i++ {
		replaced, err := ReplaceTimeTokenRaw(rawFormats[i])
		if err != nil {
			return nil, err
		}
		layouts[i] = replaced
	}

	return newLayoutSet(layouts), nil
}

func NewSingleLayout(layout string) (*LayoutSet, error) {
	replaed, err := ReplaceTimeToken(layout)
	if err != nil {
		return nil, err
	}
	return &LayoutSet{
		layouts: []string{replaed},
	}, nil
}

func (l *LayoutSet) CloneLayout() []string {
	cloend := make([]string, len(l.layouts))
	copy(cloend, l.layouts)
	return cloend
}

func (l *LayoutSet) Layout() []string {
	return l.layouts
}

func (l *LayoutSet) AddLayout(other *LayoutSet) *LayoutSet {
	setLayout := set.New[string]()
	for _, v := range l.layouts {
		setLayout.Add(v)
	}
	for _, v := range other.layouts {
		setLayout.Add(v)
	}

	return newLayoutSet(setLayout.Values().Collect())
}
