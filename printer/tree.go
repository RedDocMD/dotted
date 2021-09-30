package printer

import "fmt"

type TreeNode interface {
	Render() string
	IsLeaf() bool
	Children() []TreeNode
}

type operation = int

const (
	root operation = iota
	last
	other
)

func treePrintRecursive(node TreeNode, ops []operation) {
	for _, op := range ops[:len(ops)-1] {
		if op == other {
			fmt.Print("\u2502  ")
		} else if op == last {
			fmt.Print("   ")
		}
	}
	lastOp := ops[len(ops)-1]
	if lastOp == last {
		fmt.Print("\u2514\u2500\u2500")
	} else if lastOp == other {
		fmt.Print("\u251C\u2500\u2500")
	}
	fmt.Print(node.Render())
	children := node.Children()
	for i, child := range children {
		var newOp operation
		if i == len(children)-1 {
			newOp = last
		} else {
			newOp = other
		}
		ops = append(ops, newOp)
		treePrintRecursive(child, ops)
		ops = ops[:len(ops)-1]
	}
}

func TreePrint(node TreeNode) {
	treePrintRecursive(node, []operation{root})
}
