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

type BinaryOperator struct {
	priority uint
	op_type  OperatorType

	left  any
	right any
}
