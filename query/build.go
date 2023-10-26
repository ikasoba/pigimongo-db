package query

import (
	"encoding/json"
	"errors"
	"fmt"
)

type BuildContext struct {
	Values    []any
	Query     string
	args      []any
	argsIndex uint
}

func NewBuildContext(args ...any) *BuildContext {
	return &BuildContext{args: args, argsIndex: 0}
}

func (ctx *BuildContext) BuildQueryToWhere(n any) error {
	switch node := n.(type) {
	case BinaryOperator:
		{
			var op_type string
			switch node.op_type {
			case OpAND:
				op_type = "AND"
			case OpOR:
				op_type = "OR"
			case OpEQL:
				op_type = "="
			case OpNEQ:
				op_type = "<>"
			case OpLT:
				op_type = "<"
			case OpGT:
				op_type = ">"
			case OpLE:
				op_type = "<="
			case OpGE:
				op_type = ">="
			default:
				return errors.New("invalid operator type.")
			}

			err := ctx.BuildQueryToWhere(node.left)
			if err != nil {
				return err
			}

			ctx.Query += " " + op_type + " "

			err = ctx.BuildQueryToWhere(node.right)
			if err != nil {
				return err
			}

			return nil
		}

	case IdentLit:
		if string(node) == "Id_" || string(node) == "id_" {
			ctx.Query += `id`
		} else {
			ctx.Query += `data -> ?`
			ctx.Values = append(ctx.Values, "$."+string(node))
		}
		return nil

	case StrLit:
		ctx.Query += "?"
		data, err := json.Marshal(string(node))
		if err != nil {
			return err
		}

		ctx.Values = append(ctx.Values, string(data))
		return nil

	case NumLit:
		ctx.Query += "?"
		data, err := json.Marshal(float64(node))
		if err != nil {
			return err
		}

		ctx.Values = append(ctx.Values, string(data))
		return nil

	case ArrLit:
		return errors.New("ArrLit is not implemented.")

	case ObjLit:
		return errors.New("ObjLit is not implemented.")

	case PlaceholderLit:
		ctx.Query += "?"
		data, err := json.Marshal(ctx.args[ctx.argsIndex])
		if err != nil {
			return err
		}

		ctx.Values = append(ctx.Values, string(data))

		ctx.argsIndex += 1
		return nil

	default:
		return errors.New(fmt.Sprintf("invalid node type \"%T\".", node))
	}
}
