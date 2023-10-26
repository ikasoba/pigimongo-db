package query

type KeyChain struct {
	// Ident | StrLit | NumLit
	chain []any
}

type IdentLit string
type StrLit string
type NumLit float64
type ArrLit []any
type ObjLit map[string]any

type OperatorType uint

const (
	OpAND OperatorType = iota
	OpOR
	OpEQL
	OpNEQ
	OpLT
	OpGT
	OpLE
	OpGE
)

var OperatorPriority = map[OperatorType]uint{
	OpAND: 0,
	OpOR:  1,
	OpEQL: 2,
	OpNEQ: 2,
	OpLT:  3,
	OpGT:  3,
	OpLE:  3,
	OpGE:  3,
}

type BinaryOperator struct {
	priority uint
	op_type  OperatorType

	left  any
	right any
}

type PlaceholderLit struct{}
