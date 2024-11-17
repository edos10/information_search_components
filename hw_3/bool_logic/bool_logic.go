package boollogic

type Operation uint8

const (
	And Operation = iota
	Or
)

type Node struct {
	Words     []string
	Nodes     []*Node
	Operation Operation
	Value     uint8
}

func New(operation Operation, words []string, node []*Node) *Node {
	return &Node{
		Operation: operation,
		Words:     words,
		Nodes:     node,
	}
}

func Search(node Node) []uint32 {
	return []uint32{}
}

func RecursiveInit(node *Node, words []string, nodes []Node) *Node {
	return nil
}
