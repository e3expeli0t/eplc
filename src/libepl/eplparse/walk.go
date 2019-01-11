package eplparse

import "eplc/src/libepl/eplparse/ast"

type Visitor interface {

	Visit(node ast.Node) (v Visitor)
}

func Walk(v Visitor, node ast.Node) {

}
