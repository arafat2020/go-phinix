package compiler

import (
	"fmt"

	"github.com/arafat2020/go-phoenix/internal/ast"
)

var opSymbols = map[ast.Op]string{
	ast.OpEq:  "=",
	ast.OpNeq: "!=",
	ast.OpGT:  ">",
	ast.OpGte: ">=",
	ast.OpLT:  "<",
	ast.OpLte: "<=",
}

func CompilePredicate(p ast.Predicate, position int) (sql string, arg any, err error) {
	symbol, ok := opSymbols[p.Op]

	if !ok {
		return "", nil, fmt.Errorf("compiler: unknown op %v", p.Op)
	}

	return fmt.Sprintf("%s %s $%d", p.Column, symbol, position), p.Value, nil
}
