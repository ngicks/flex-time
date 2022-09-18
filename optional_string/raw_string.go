package optionalstring

import "github.com/ngicks/type-param-common/slice"

type valueType int

const (
	Normal valueType = iota
	SingleQuoteEscaped
	SlashEscaped
)

type TextNode struct {
	typ   valueType
	value string
}

func (v TextNode) Typ() valueType {
	return v.typ
}

func (v TextNode) Len() int {
	return len(v.value)
}

func (v TextNode) Value() string {
	return v.value
}

func (v TextNode) Unescaped() string {
	switch v.typ {
	case Normal:
		return v.value
	case SingleQuoteEscaped:
		return v.Value()[1 : v.Len()-1]
	case SlashEscaped:
		return v.Value()[1:]
	}
	panic("unknown")
}

type RawString slice.Deque[TextNode]

func NewRawString() RawString {
	return make(RawString, 0)
}

func (rs RawString) Append(str ...RawString) RawString {
	c := rs.Clone()
	for _, v := range str {
		(*slice.Deque[TextNode])(&c).Append(v...)
	}
	return c
}

func (rs RawString) Clone() RawString {
	cloned := (*slice.Deque[TextNode])(&rs).Clone()
	return RawString(cloned)
}

func (rs RawString) String() string {
	var out string
	for _, v := range rs {
		out += v.value
	}
	return out
}

func (rs RawString) Unescaped() string {
	var out string
	for _, v := range rs {
		out += v.Unescaped()
	}
	return out
}
