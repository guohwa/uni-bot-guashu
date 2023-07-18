package pongo2gin

import (
	"github.com/flosch/pongo2"
)

type tagPsetNode struct {
	name       string
	expression pongo2.IEvaluator
}

func (node *tagPsetNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	// Evaluate expression
	value, err := node.expression.Evaluate(ctx)
	if err != nil {
		return err
	}

	ctx.Public[node.name] = value
	return nil
}

func tagPsetParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &tagPsetNode{}

	// Parse variable name
	typeToken := arguments.MatchType(pongo2.TokenIdentifier)
	if typeToken == nil {
		return nil, arguments.Error("Expected an identifier.", nil)
	}
	node.name = typeToken.Val

	if arguments.Match(pongo2.TokenSymbol, "=") == nil {
		return nil, arguments.Error("Expected '='.", nil)
	}

	// Variable expression
	keyExpression, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.expression = keyExpression

	// Remaining arguments
	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed 'pset'-tag arguments.", nil)
	}

	return node, nil
}
