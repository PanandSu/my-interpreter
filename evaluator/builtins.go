package builtins

import (
	"fmt"
	"my-interpreter/object"
)

var builtins = map[string]*object.Builtins{
	// len(string)或者len(array)
	"len": {Fn: func(params ...object.Object) object.Object {
		if len(params) != 1 {
			return newError("wrong number of arguments. got=%d, want=1",
				len(params),
			)
		}
		switch param := params[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(param.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(param.Elements))}
		default:
			return newError("argument to `len` not supported, got %s", param.Type())
		}
	}},
	//  first(array)
	"first": {Fn: func(params ...object.Object) object.Object {
		if len(params) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(params))
		}
		if params[0].Type() != object.ARRAY {
			return newError("argument to `first` must be INTEGER, got %s", params[0].Type())
		}
		arr := params[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return Nil
	}},
	// last(array)
	"last": {Fn: func(params ...object.Object) object.Object {
		if len(params) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(params))
		}
		if params[0].Type() != object.ARRAY {
			return newError("argument to `last` must be INTEGER, got %s", params[0].Type())
		}
		arr := params[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[len(arr.Elements)-1]
		}
		return Nil
	}},
	// push(array,element)
	"push": {Fn: func(params ...object.Object) object.Object {
		if len(params) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(params))
		}
		if params[0].Type() != object.ARRAY {
			return newError("argument to `push` not supported, got %s", params[0].Type())
		}
		arr := params[0].(*object.Array)
		length := len(arr.Elements)

		newArr := make([]object.Object, length+1)
		copy(newArr, arr.Elements)

		newArr[length] = params[1]
		return &object.Array{Elements: newArr}
	}},
	// pop(array)
	"pop": {Fn: func(params ...object.Object) object.Object {
		if len(params) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(params))
		}
		if params[0].Type() != object.ARRAY {
			return newError("argument to `pop` not supported, got %s", params[0].Type())
		}
		arr := params[0].(*object.Array)
		length := len(arr.Elements)

		newArr := make([]object.Object, length-1)
		copy(newArr, arr.Elements[:length-1])
		return &object.Array{Elements: newArr}
	}},
	// prints(任意数量任何类型的数据)
	"prints": {Fn: func(params ...object.Object) object.Object {
		for _, param := range params {
			fmt.Println(param.Inspect())
		}
		return Nil
	}},
}
