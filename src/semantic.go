package main

import (
	"fmt"
	"strings"
)

type Symbol struct {
	Type     TokenType
	Declared bool
}

var builtins = map[string]bool{
	"print": true,
}

func inferType(expr Expression, symbols map[string]Symbol) (TokenType, error) {
	switch e := expr.(type) {
	case *Number_lit:
		if strings.Contains(e.Value, ".") {
			return FLOAT, nil
		}
		return INT, nil

	case *Identifier_expr:
		symbol, exists := symbols[e.Name]
		if !exists {
			return 0, fmt.Errorf("variable %v not declared", e.Name)
		}
		return symbol.Type, nil

	// case *Call_expr:
	// 	// can be used later when return is added

	case *Binary_expr:
		leftType, err := inferType(e.Left, symbols)
		if err != nil {
			return 0, err
		}
		rightType, err := inferType(e.Right, symbols)
		if err != nil {
			return 0, err
		}
		if leftType == rightType {
			return leftType, nil
		}
		return 0, fmt.Errorf("")
	}
	return 0, fmt.Errorf("unknown type")
}

func Analyze(p *Program) error {

	// Initialize semantic tables
	SymbolTable := make(map[string]Symbol)
	FunctionTable := make(map[string]*Func_Decl)

	// Check for redeclared functions
	for _, function := range p.Functions {
		if _, exists := FunctionTable[function.Name]; exists {
			return fmt.Errorf("%v redaclared in this block", function.Name)
		}
		FunctionTable[function.Name] = function

		for _, statement := range function.Body.Statement_Group {
			switch s := statement.(type) {
			case *Decl_stmt:
				// check for redeclared variables
				if _, exists := SymbolTable[s.Name]; exists {
					return fmt.Errorf("%v variable redaclared in this block", s.Name)
				}

				// populate symbol table
				symbol := Symbol{Type: s.Type, Declared: true}
				SymbolTable[s.Name] = symbol

			case *Assign_stmt:
				// check for assignment before declaration
				symbol, exists := SymbolTable[s.Target]
				if !exists {
					return fmt.Errorf("Assignment to undeclared variable %v", s.Target)
				}

				// type check
				assignType, err := inferType(s.Value, SymbolTable)
				if err != nil {
					return err
				}
				if assignType != symbol.Type {
					return fmt.Errorf("Cant assign %s to variable %v of type %s", assignType, s.Value, symbol.Type.String())
				}

			case *Update_stmt:
				// check for assignment before declaration
				symbol, exists := SymbolTable[s.Target]
				if !exists {
					return fmt.Errorf("Update to undeclared variable %v", s.Target)
				}
				if symbol.Type != INT && symbol.Type != FLOAT {
					return fmt.Errorf("Cannot perform update on variable of type %v", symbol.Type.String())
				}

			case *Expr_stmt:
				switch e := s.Expr.(type) {
				case *Call_expr:
					if builtins[e.Name] {
						break
					}
					if _, exists := FunctionTable[e.Name]; !exists {
						return fmt.Errorf("Call to undeclared function %v", e.Name)
					}
				}
			}
		}

	}

	return nil
}
