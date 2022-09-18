package optionalstring

import (
	"fmt"

	"github.com/pkg/errors"
	parsec "github.com/prataprc/goparsec"
)

const (
	OPENSQR           = "OPENSQR"
	CLOSESQR          = "CLOSESQR"
	SQUOTE            = "SQUOTE"
	ESCAPEDCHAR       = "ESCAPEDCHAR"
	NORMALCHARS       = "NORMALCHARS"
	CHAR              = "CHAR"
	CHARS             = "CHARS"
	CHARWITHINESCAPE  = "CHARWITHINESCAPE"
	CHARSWITHINESCAPE = "CHARSWITHINESCAPE"
	ESCAPED           = "ESCAPED"
	ITEM              = "ITEM"
	ITEMS             = "ITEMS"
	OPTIONAL          = "OPTIONAL"
	OPTIONALSTRING    = "OPTIONALSTRING"
)

var (
	opensqr     parsec.Parser = parsec.Atom(`[`, OPENSQR)
	closesqr                  = parsec.Atom(`]`, CLOSESQR)
	squote                    = parsec.Atom(`'`, SQUOTE)
	escapedchar               = parsec.Token(`\\.`, ESCAPEDCHAR)
	normalchars               = parsec.Token(`[^\[\]\\']+`, NORMALCHARS)
)

func MakeOptionalStringParser(ast *parsec.AST) parsec.Parser {
	char := ast.OrdChoice(CHAR, nil, escapedchar, normalchars)
	chars := ast.Many(CHARS, nil, char)
	charWithinEscape := ast.OrdChoice(CHARWITHINESCAPE, nil, escapedchar, normalchars, opensqr, closesqr)
	charsWithinEscape := ast.Many(CHARSWITHINESCAPE, nil, charWithinEscape)

	var optional parsec.Parser
	escaped := ast.And(ESCAPED, nil, squote, charsWithinEscape, squote)
	item := ast.OrdChoice(ITEM, nil, chars, escaped, &optional)
	items := ast.Kleene(ITEMS, nil, item)
	optional = ast.And(OPTIONAL, nil, opensqr, items, closesqr)
	return ast.Kleene(OPTIONALSTRING, nil, ast.OrdChoice("items", nil, optional, chars))
}

type SyntaxError struct {
	Input    string
	ParsedAs string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf(
		"syntax error: maybe no opening/closing sqrt? parsed result = %s, input = %s",
		e.ParsedAs,
		e.Input,
	)
}

func EnumerateOptionalStringRaw(optionalString string) (enumerated []rawString, err error) {
	var node parsec.Queryable
	func() {
		defer func() {
			if rcv := recover(); rcv != nil {
				err = errors.Errorf("%+v", rcv)
			}
		}()

		ast := parsec.NewAST("optionalString", 100)
		p := MakeOptionalStringParser(ast)
		s := parsec.NewScanner([]byte(optionalString))
		node, _ = ast.Parsewith(p, s)
	}()

	if err != nil {
		return
	}

	if parsedAs := node.GetValue(); len(parsedAs) != len(optionalString) {
		return []rawString{}, &SyntaxError{
			Input:    optionalString,
			ParsedAs: parsedAs,
		}
	}

	root := decode(node)

	return root.Flatten(), nil
}

func EnumerateOptionalString(optionalString string) (enumerated []string, err error) {
	raw, err := EnumerateOptionalStringRaw(optionalString)
	if err != nil {
		return []string{}, err
	}

	out := make([]string, len(raw))
	for idx, v := range raw {
		out[idx] = v.String()
	}
	return out, nil
}
