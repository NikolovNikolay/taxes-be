package parser

import (
	"fmt"
	"taxes-be/internal/core"
)

type Factory interface {
	Build(t core.StatementType) (Parser, error)
}

type parserFactory struct{}

func NewParserFactory() Factory {
	return &parserFactory{}
}

func (f *parserFactory) Build(t core.StatementType) (Parser, error) {
	switch t {
	case core.Revolut:
		return newRevolutStatementParser(), nil
	case core.EToro:
		return newEToroStatementParserLinked(), nil
	default:
		return nil, fmt.Errorf("unknown statement type")
	}
}
