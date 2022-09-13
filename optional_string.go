package flextime

import (
	"fmt"

	"github.com/ngicks/type-param-common/slice"
	"github.com/pkg/errors"
	parsec "github.com/prataprc/goparsec"
)

const (
	OPENSQR     = "OPENSQR"
	CLOSESQR    = "CLOSESQR"
	ESCAPEDCHAR = "ESCAPEDCHAR"
	NORMALCHARS = "NORMALCHARS"
)

var (
	opensqr     parsec.Parser = parsec.Atom(`[`, OPENSQR)
	closesqr                  = parsec.Atom(`]`, CLOSESQR)
	escapedchar               = parsec.Token(`\\.`, ESCAPEDCHAR)
	normalchars               = parsec.Token(`[^\[\]\\]+`, NORMALCHARS)
)

func MakeOptionalStringParser(ast *parsec.AST) parsec.Parser {
	char := ast.OrdChoice("char", nil, escapedchar, normalchars)
	chars := ast.Many("chars", nil, char)

	var optional parsec.Parser
	item := ast.OrdChoice("item", nil, chars, &optional)
	items := ast.Kleene("items", nil, item)
	optional = ast.And("optional", nil, opensqr, items, closesqr)
	return ast.Kleene("optionalString", nil, ast.OrdChoice("items", nil, optional, chars))
}

type SyntaxError struct {
	org      string
	parsedAs string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf(
		"syntax error: maybe no opening/closing sqrt? parsed result = %s, input = %s",
		e.parsedAs,
		e.org,
	)
}

func EnumerateOptionalString(optionalString string) (enumerated []string, err error) {
	defer func() {
		if rcv := recover(); rcv != nil {
			err = errors.Errorf("%+v", rcv)
		}
	}()

	ast := parsec.NewAST("optionalString", 100)
	p := MakeOptionalStringParser(ast)
	s := parsec.NewScanner([]byte(optionalString))
	node, _ := ast.Parsewith(p, s)

	if parsedAs := node.GetValue(); len(parsedAs) != len(optionalString) {
		return []string{}, &SyntaxError{
			org:      optionalString,
			parsedAs: parsedAs,
		}
	}

	root := optStrContext{
		node: &treeNode{},
	}
	recursiveDecode(node, &root)

	return root.Trunc(), nil
}

type optStrContext struct {
	node *treeNode
}

func (ctx *optStrContext) Open() {
	ctx.node.open = true
}
func (ctx *optStrContext) IsOpened() bool {
	return ctx.node.open
}
func (ctx *optStrContext) Close() {
	ctx.node.open = false
}

func (ctx *optStrContext) IsRoot() bool {
	return ctx.node.parent == nil
}

func (ctx *optStrContext) StepOptional() *optStrContext {
	if ctx.node.optional == nil {
		ctx.node.optional = &treeNode{parent: ctx.node}
	}
	return &optStrContext{
		node: ctx.node.optional,
	}
}

func (ctx *optStrContext) StepNext() *optStrContext {
	if ctx.node.next == nil {
		ctx.node.next = &treeNode{parent: ctx.node}
	}
	return &optStrContext{
		node: ctx.node.next,
	}
}

func (ctx *optStrContext) StepBack() *optStrContext {
	return &optStrContext{
		node: ctx.node.parent,
	}
}

func (ctx *optStrContext) AddString(str string) {
	ctx.node.value += str
}

func (ctx *optStrContext) Trunc() []string {
	return ctx.trunc()
}

func (ctx *optStrContext) trunc() []string {
	cur := ctx.node.value
	ret := []string{cur}

	if ctx.node.optional != nil {
		for _, s := range ctx.StepOptional().trunc() {
			ret = append(ret, cur+s)
		}
	}

	if ctx.node.next != nil {
		nextStr := ctx.StepNext().trunc()
		_, hasNonEmpty := slice.Find(nextStr, func(v string) bool { return v != "" })
		if hasNonEmpty {
			org := make([]string, len(ret))
			copy(org, ret)
			ret = ret[:0]
			for _, nn := range nextStr {
				for _, oo := range org {
					ret = append(ret, oo+nn)
				}
			}
		}
	}

	return ret
}

type treeNode struct {
	parent   *treeNode
	open     bool
	value    string
	optional *treeNode
	next     *treeNode
}

func recursiveDecode(node parsec.Queryable, ctx *optStrContext) {
	for _, child := range node.GetChildren() {
		switch child.GetName() {
		case OPENSQR:
			if ctx.IsOpened() {
				ctx = ctx.StepNext()
			}
			ctx.Open()
			ctx = ctx.StepOptional()
		case CLOSESQR:
			ctx = ctx.StepBack()
		case ESCAPEDCHAR, NORMALCHARS:
			if ctx.IsOpened() {
				ctx = ctx.StepNext()
			}
			ctx.AddString(child.GetValue())
		}
		recursiveDecode(child, ctx)
	}
}
