package compiler_test

import (
	"testing"

	"github.com/arafat2020/go-phoenix/internal/ast"
	"github.com/arafat2020/go-phoenix/internal/compiler"
)

func TestCompilePredicate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   ast.Predicate
		position    int
		expectedSQL string
		expectedArg any
		wantErr     bool
	}{
		{
			name: "greater than predicate at position 1",
			predicate: ast.Predicate{
				Column: "age",
				Op:     ast.OpGT,
				Value:  18,
			},
			position:    1,
			expectedSQL: "age > $1",
			expectedArg: 18,
			wantErr:     false,
		},
		{
			name: "equals predicate at position 2",
			predicate: ast.Predicate{
				Column: "active",
				Op:     ast.OpEq,
				Value:  true,
			},
			position:    2,
			expectedSQL: "active = $2",
			expectedArg: true,
			wantErr:     false,
		},
		{
			name: "invalid op value returns error",
			predicate: ast.Predicate{
				Column: "invalid_col",
				Op:     ast.Op(999),
				Value:  "val",
			},
			position:    1,
			expectedSQL: "",
			expectedArg: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, arg, err := compiler.CompilePredicate(tt.predicate, tt.position)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CompilePredicate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if sql != tt.expectedSQL {
				t.Errorf("CompilePredicate() sql = %q, want %q", sql, tt.expectedSQL)
			}
			if arg != tt.expectedArg {
				t.Errorf("CompilePredicate() arg = %v, want %v", arg, tt.expectedArg)
			}
		})
	}
}
