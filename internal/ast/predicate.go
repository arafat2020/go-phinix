package ast

type Op int

const (
	OpEq Op = iota
	OpNeq
	OpGT
	OpGte
	OpLT
	OpLte
)

type Predicate struct {
	Column string
	Op     Op
	Value  any
}
