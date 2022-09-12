package flextime

import (
	"github.com/pkg/errors"
	parsec "github.com/prataprc/goparsec"
)

func MakeOptionalStringParser(ast *parsec.AST) parsec.Parser {
	opensqr := parsec.Atom(`[`, "OPENSQR")
	closesqr := parsec.Atom(`]`, "CLOSESQR")
	escapedchar := parsec.Token(`\\.`, "ESCAPEDCHAR")
	normalchar := parsec.Token(`[^\[\]\\]+`, "NORMALCHAR")

	char := ast.OrdChoice("char", nil, escapedchar, normalchar)
	chars := ast.Many("chars", nil, char)

	var optional parsec.Parser
	item := ast.OrdChoice("item", nil, chars, &optional)
	items := ast.Kleene("items", nil, item)
	optional = ast.And("optional", nil, opensqr, items, closesqr)
	return ast.Kleene("optionalString", nil, ast.OrdChoice("items", nil, optional, chars))
}

func EnumerateOptionalString(optionalString string) (enumerated []string, err error) {
	defer func() {
		if rcv := recover(); rcv != nil {
			err = errors.Errorf("%+v", rcv)
		}
	}()

	ast := parsec.NewAST("optionalString", 100)
	p := MakeOptionalStringParser(ast)
	s := parsec.NewScanner([]byte(optionalString)).TrackLineno()
	ast.Parsewith(p, s)

	return []string{}, nil
}
