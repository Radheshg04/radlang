package semantic

import "fmt"

type Severity int
type DiagnosticCode string

const (
	Error Severity = iota
	Warn
)

var diagnosticMessages = map[DiagnosticCode]string{
	ErrRedeclaredBuiltinFunction: "cannot redeclare builtin function",
	ErrRedeclaredSameFunc:        "function already declared in this scope",
	ErrArgsMustMatchParams:       "argument count does not match parameter count",
	ErrIfExpressionNotBool:       "if condition must be a boolean expression",
	ErrExpectedOneExpr:           "expected exactly one expression",
	ErrRedeclaredVariables:       "variable already declared in this scope",
	ErrUnusedDeclaredVariable:    "variable declared but not used",
	ErrIdentNotDeclared:          "identifier not declared",
	ErrNoNewVariablesOnWalrus:    "no new variables on left side of :=",
	ErrTooManyReturnValues:       "too many return values",
	ErrNotEnoughReturnValues:     "not enough return values",
	ErrReturnOutsideBlock:        "return statement outside of function",
	ErrBadReturnType:             "return type does not match function signature",
	ErrJumpOutsideFor:            "break/continue outside of for loop",
	ErrInvalidExprStmt:           "expression is not a valid statement",
	ErrInvalidOperand:            "invalid operand for operator",
	ErrMismatchTypesInExpr:       "mismatched types in expression",
	ErrUndefined:                 "undefined reference",
	WarnUnreachableCode:          "unreachable code",
	ErrPostfixOnNonNumeric:       "cannot perform posfix op on non numeric",
}

type Diagnostic struct {
	Severity       Severity
	DiagnosticCode DiagnosticCode
	Span           *Span
}

type Span struct {
	line int
	col  int
}

const (
	// Functions
	ErrRedeclaredBuiltinFunction DiagnosticCode = "RL001"
	ErrRedeclaredSameFunc        DiagnosticCode = "RL002"

	// Statements
	ErrArgsMustMatchParams    DiagnosticCode = "RL401"
	ErrIfExpressionNotBool    DiagnosticCode = "RL402"
	ErrExpectedOneExpr        DiagnosticCode = "RL403"
	ErrRedeclaredVariables    DiagnosticCode = "RL404"
	ErrUnusedDeclaredVariable DiagnosticCode = "RL405"
	ErrIdentNotDeclared       DiagnosticCode = "RL406"
	ErrNoNewVariablesOnWalrus DiagnosticCode = "RL407"
	ErrTooManyReturnValues    DiagnosticCode = "RL408"
	ErrNotEnoughReturnValues  DiagnosticCode = "RL409"
	ErrReturnOutsideBlock     DiagnosticCode = "RL410"
	ErrBadReturnType          DiagnosticCode = "RL411"
	ErrJumpOutsideFor         DiagnosticCode = "RL412"
	ErrInvalidExprStmt        DiagnosticCode = "RL413"
	ErrPostfixOnNonNumeric    DiagnosticCode = "RL414"

	// Expressions
	ErrInvalidOperand      DiagnosticCode = "RL501"
	ErrMismatchTypesInExpr DiagnosticCode = "RL502"

	// Misc
	ErrUndefined        DiagnosticCode = "RL701"
	WarnUnreachableCode DiagnosticCode = "RL702"
)

func NewRLDiagnostic(code DiagnosticCode) *Diagnostic {

	diagnostic := &Diagnostic{DiagnosticCode: code}
	switch code {
	// Case Error
	case ErrRedeclaredBuiltinFunction, ErrRedeclaredSameFunc,
		ErrArgsMustMatchParams, ErrExpectedOneExpr, ErrIfExpressionNotBool,
		ErrUnusedDeclaredVariable, ErrRedeclaredVariables, ErrIdentNotDeclared,
		ErrNoNewVariablesOnWalrus, ErrTooManyReturnValues, ErrNotEnoughReturnValues,
		ErrReturnOutsideBlock, ErrBadReturnType, ErrJumpOutsideFor, ErrInvalidExprStmt,
		ErrPostfixOnNonNumeric, ErrInvalidOperand, ErrMismatchTypesInExpr, ErrUndefined:
		diagnostic.Severity = Error
	// Case Warn
	case WarnUnreachableCode:
		diagnostic.Severity = Warn
	default:
		return nil
	}
	return diagnostic
}

func Report(diagnostics []Diagnostic) {
	for _, diag := range diagnostics {
		severity := "error"
		if diag.Severity == Warn {
			severity = "warn"
		}
		msg, ok := diagnosticMessages[diag.DiagnosticCode]
		if !ok {
			msg = "unknown diagnostic"
		}
		if diag.Span != nil {
			fmt.Printf("%s [%s] %d:%d: %s\n", severity, diag.DiagnosticCode, diag.Span.line, diag.Span.col, msg)
		} else {
			fmt.Printf("%s [%s]: %s\n", severity, diag.DiagnosticCode, msg)
		}
	}
}
