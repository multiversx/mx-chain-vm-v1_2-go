package mandosjsonparse

import (
	ei "github.com/multiversx/mx-chain-vm-v1_2-go/mandos-go/expression/interpreter"
	fr "github.com/multiversx/mx-chain-vm-v1_2-go/mandos-go/fileresolver"
)

// Parser performs parsing of both json tests (older) and scenarios (new).
type Parser struct {
	ExprInterpreter ei.ExprInterpreter
}

// NewParser provides a new Parser instance.
func NewParser(fileResolver fr.FileResolver) Parser {
	return Parser{
		ExprInterpreter: ei.ExprInterpreter{
			FileResolver: fileResolver,
		},
	}
}
