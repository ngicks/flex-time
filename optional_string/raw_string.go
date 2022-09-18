package optionalstring

type valueType int

const (
	Normal valueType = iota
	SingleQuoteEscaped
	SlashEscaped
)

type Value struct {
	typ   valueType
	value string
}

func (v Value) Typ() valueType {
	return v.typ
}

func (v Value) Len() int {
	return len(v.value)
}

func (v Value) Value() string {
	return v.value
}

type rawString []Value

func newRawString() rawString {
	return make(rawString, 0)
}

func (rs rawString) Clone() rawString {
	cloned := make(rawString, len(rs))
	copy(cloned, rs)
	return cloned
}

func (rs rawString) Append(v rawString) rawString {
	return append(rs, v...)
}

func (rs rawString) String() string {
	var out string
	for _, v := range rs {
		out += v.value
	}
	return out
}
