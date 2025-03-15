package golang

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/gopls/internal/cache"
	"golang.org/x/tools/gopls/internal/cache/methodsets"
)

func getArgComposite(arg ast.Expr) *ast.CompositeLit {
	var ret *ast.CompositeLit
	ret, ok := arg.(*ast.CompositeLit)
	if ok {
		return ret
	}

	// ptr
	unaryExpr, ok := arg.(*ast.UnaryExpr)
	if !ok {
		return nil
	}
	ret, ok = unaryExpr.X.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	return ret

}

func callExprArgMatch(pkg *cache.Package, methodStructType types.Type, methodType types.Type, arg ast.Expr) bool {
	if arg == nil {
		return false
	}
	if methodStructType == nil {
		return false
	}
	var argIdent *ast.Ident

	argIdent, ok := arg.(*ast.Ident)
	if !ok {
		compositeLit := getArgComposite(arg)
		if compositeLit == nil {
			return false
		}
		argIdent, ok = compositeLit.Type.(*ast.Ident)
		if !ok {
			return false
		}
	}

	argType, ok := pkg.TypesInfo().Uses[argIdent]
	if !ok {
		return false
	}
	if argType == nil {
		return false
	}
	argIdentType := argType.Type()
	if argIdentType == nil {
		return false
	}
	argIdentPtrType := methodsets.EnsurePointer(argIdentType)
	return argIdentPtrType.String() == methodStructType.String()
}

func callExprParamMatch(methodType types.Type, methodName string, param *types.Var) bool {
	if param == nil {
		return false
	}
	paramType := param.Type()
	paramTypeNamed, ok := paramType.(*types.Named)
	if !ok || paramTypeNamed == nil {
		return false
	}
	underLying := paramTypeNamed.Underlying()
	underLyintInterface, ok := underLying.(*types.Interface)
	if !ok || underLyintInterface == nil {
		return false
	}

	for _method := range underLyintInterface.Methods() {
		if _method.Name() != methodName {
			continue
		}
		match := types.Identical(_method.Type(), methodType)
		if !match {
			continue
		}
		return true
	}
	return false
}

func CallExprMatches(pkg *cache.Package, methodStructType types.Type, methodType types.Type, methodName string, callExpr *ast.CallExpr) bool {
	if callExpr == nil {
		return false
	}
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	selExprtype, ok := pkg.TypesInfo().Uses[selExpr.Sel]
	if !ok {
		return false
	}
	funcType, ok := selExprtype.(*types.Func)
	if !ok {
		return false
	}

	sinatureType := funcType.Signature()
	if sinatureType == nil {
		return false
	}
	params := sinatureType.Params()
	// fmt.Printf("SelectorExpr expr %s\n", selExpr.Sel.Name)
	for i, arg := range callExpr.Args {
		if !callExprArgMatch(pkg, methodStructType, methodType, arg) {
			continue
		}
		param := params.At(i)
		if callExprParamMatch(methodType, methodName, param) {
			return true
		}
	}

	return false
}
