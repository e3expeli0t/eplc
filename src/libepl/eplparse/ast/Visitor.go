package ast

type Visitor interface {
	Visit(n Node) Visitor
}


func Walk(v Visitor, node Node) {
	if v := v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case Decl:
		Walk(v, n)
	case *Block:
		for _,sn := range *n.Nodes {
			Walk(v, sn)
		}
	case *ProgramFile:
		for _, decl := range *n.GlobalDecls {
			Walk(v, decl)
		}

	}
}

type Traveler func(Node) bool


func (f Traveler) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}


func Travel(node Node, f func(Node) bool) {
	Walk(Traveler(f), node)
}