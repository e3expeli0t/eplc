package ast

type Visitor interface {
	Visit(node Node) (v Visitor)
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
