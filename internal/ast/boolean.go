package ast

type And struct {
	Children []Node
}

func (And) isNode() {}

type Or struct {
	Children []Node
}

func (Or) isNode() {}
