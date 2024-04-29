package builtins

import (
	"fmt"
	"my-interpreter/ast"
	"my-interpreter/object"
)

var (
	Nil   = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expr, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.Return{Value: val}

	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StrLiteral:
		return &object.String{Value: node.Token.Literal}

	case *ast.BoolLiteral:
		return nativeBool2BooleanObject(node.Value)

	case *ast.ArrLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.MapLiteral:
		return evalMapLiteral(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IfExpression:
		return evalIfExpr(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpr(node.Token.Literal, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpr(node.Token.Literal, left, right)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		index := Eval(node.Index, env)
		return evalIndexExpr(left, index)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var res object.Object
	for _, statement := range program.Statements {
		res = Eval(statement, env)
		switch obj := res.(type) {
		case *object.Return:
			return obj.Value
		case *object.Error:
			return obj
		}
	}
	return res
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object
	for _, statement := range block.Statements {
		res = Eval(statement, env)
		if res != nil {
			typ := res.Type()
			if typ == object.RETURN || typ == object.ERROR {
				return res
			}
		}
	}
	return res
}

func evalIndexExpr(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrIndex(left, index)
	case left.Type() == object.MAP: //将在evalMapIndex函数中判断index是否可哈希
		return evalMapIndex(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrIndex(array object.Object, index object.Object) object.Object {
	arr := array.(*object.Array)
	idx := index.(*object.Integer).Value
	length := int64(len(arr.Elements))
	if idx < 0 || idx >= length {
		return newError("index out of range")
	}
	return arr.Elements[idx]
}

func evalMapIndex(m, index object.Object) object.Object {
	obj := m.(*object.Map)
	//将index断言为Hashable,即可哈希
	i, err := index.(object.Hashable)
	if !err {
		return newError("unhashable type : %s", index.Type())
	}
	/*
		Type assertions are used to check that a variable is of some type and return the underlying interface value.
		Type assertions work only for interfaces.
		For example, in the following code: 'var x interface{} = 42 t := x.(int)',
		'x' has the 'interface{}' type with the underlying int value ('42'),
		'int' is the concrete type that we want to check.
		If we print 't', the output will be '42'.
		Changing of the concrete type to 'string' ('t := x.(string)')
		will cause a runtime panic. Type assertions can return two values.
		For example, the expression 't, ok := x.(int)' has the boolean 'ok' that
		returns 'true' if the assertion is correct.
		If 'ok' is 'false', 't' is set to a zero value and no panic occurs.
	*/
	pairs, ok := obj.Mappings[i.Hash()]
	if !ok {
		//没找到
		return Nil
	}
	return pairs.Value
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var res []object.Object
	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case True:
		return true
	case False:
		return false
	case Nil:
		return false
	default:
		return false
	}
}

func evalIfExpr(node *ast.IfExpression, env *object.Environment) object.Object {
	condf := Eval(node.Condition, env)
	if isError(condf) {
		return condf
	}
	if isTruthy(condf) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	} else {
		//如果else语句的Alternative为空
		return Nil
	}
}

func nativeBool2BooleanObject(flag bool) *object.Boolean {
	if flag {
		return True
	}
	return False
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusOperatorExpr(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalBangOperatorExpr(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Nil:
		return False
	default:
		return False
	}
}

func evalMinusOperatorExpr(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	val := right.(*object.Integer).Value
	return &object.Integer{Value: val}
}

func evalInfixExpr(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == right.Type() && left.Type() == object.INTEGER:
		return evalIntInfixExpr(op, left, right)
	case left.Type() == right.Type() && left.Type() == object.STRING:
		return evalStringInfixExpr(op, left, right)
	case op == "==":
		return nativeBool2BooleanObject(left == right)
	case op == "!=":
		return nativeBool2BooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch :%s %s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpr(op string, left, right object.Object) object.Object {
	if op != "+" {
		return newError("unknown operator :%s %s %s", left.Type(), op, right.Type())
	}
	lft := left.(*object.String).Value
	rgt := right.(*object.String).Value
	return &object.String{Value: lft + rgt}
}

func evalIntInfixExpr(op string, left, right object.Object) object.Object {
	lft := left.(*object.Integer).Value
	rgt := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: rgt + lft}
	case "-":
		return &object.Integer{Value: rgt - lft}
	case "*":
		return &object.Integer{Value: rgt * lft}
	case "/":
		return &object.Integer{Value: rgt / lft}
	case ">":
		return &object.Boolean{Value: rgt > lft}
	case "<":
		return &object.Boolean{Value: rgt < lft}
	case "==":
		return &object.Boolean{Value: rgt == lft}
	case "!=":
		return &object.Boolean{Value: rgt != lft}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Token.Literal); ok {
		return val
	}
	if val, ok := builtins[node.Token.Literal]; ok {
		return val
	}
	return newError("identifier not found: %s", node.Token.Literal)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Builtins:
		return fn.Fn(args...)
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapFunctionReturn(evaluated)
	default:
		return newError("unknown function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Token.Literal, args[paramIdx])
	}
	return env
}

func unwrapFunctionReturn(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.Return); ok {
		return returnValue.Value
	}
	return obj
}

func evalMapLiteral(m *ast.MapLiteral, env *object.Environment) object.Object {
	res := make(map[object.HashKey]*object.Pair)
	for keyNode, valNode := range m.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unhashable type as hash key: %s", key.Type())
		}
		val := Eval(valNode, env)
		if isError(val) {
			return val
		}
		hashed := hashKey.Hash()
		res[hashed] = &object.Pair{
			Key:   key,
			Value: val,
		}
	}
	return &object.Map{Mappings: res}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
