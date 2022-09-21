package optionalstring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextNode(t *testing.T) {
	var tn TextNode

	tn = TextNode{typ: Normal, value: "aaaa"}
	assert.Equal(t, tn.Value(), "aaaa")
	assert.Equal(t, tn.Unescaped(), "aaaa")
	assert.Equal(t, tn.Len(), 4)
	assert.Equal(t, tn.Typ(), Normal)

	tn = TextNode{typ: SingleQuoteEscaped, value: `'aaaa'`}
	assert.Equal(t, tn.Value(), `'aaaa'`)
	assert.Equal(t, tn.Unescaped(), `aaaa`)
	assert.Equal(t, tn.Len(), 6)
	assert.Equal(t, tn.Typ(), SingleQuoteEscaped)

	tn = TextNode{typ: SlashEscaped, value: `\a`}
	assert.Equal(t, tn.Value(), `\a`)
	assert.Equal(t, tn.Unescaped(), `a`)
	assert.Equal(t, tn.Len(), 2)
	assert.Equal(t, tn.Typ(), SlashEscaped)
}
