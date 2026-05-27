package semantic

type Severity int
type DiagnosticCode string

const (
	Error Severity = iota
	Warn
)

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
	ErrUnusedDeclaredVariable DiagnosticCode = "RL403"
	ErrRedeclaredVariables    DiagnosticCode = "RL404"
	ErrIdentNotDeclared       DiagnosticCode = "RL405"
	ErrNoNewVariablesOnWalrus DiagnosticCode = "RL406"
	ErrTooManyReturnValues    DiagnosticCode = "RL407"
	ErrNotEnoughReturnValues  DiagnosticCode = "RL408"
	ErrReturnOutsideBlock     DiagnosticCode = "RL409"
	ErrBadReturnType          DiagnosticCode = "RL410"
	ErrJumpOutsideFor         DiagnosticCode = "RL411"
	ErrInvalidExprStmt        DiagnosticCode = "RL412"
	WarnUnreachableCode       DiagnosticCode = "RL413"

	// Expressions
	ErrInvalidOperand DiagnosticCode = "RL501"

	// Misc
)

func NewRLDiagnostic(code DiagnosticCode) *Diagnostic {
	diagnostic := &Diagnostic{DiagnosticCode: code}
	switch code {
	// Case Error
	case ErrRedeclaredBuiltinFunction, ErrRedeclaredSameFunc,
		ErrArgsMustMatchParams, ErrIfExpressionNotBool, ErrUnusedDeclaredVariable,
		ErrRedeclaredVariables, ErrIdentNotDeclared, ErrNoNewVariablesOnWalrus,
		ErrTooManyReturnValues, ErrNotEnoughReturnValues, ErrReturnOutsideBlock,
		ErrBadReturnType, ErrJumpOutsideFor, ErrInvalidExprStmt, ErrInvalidOperand:
		diagnostic.Severity = Error
	// Case Warn
	case WarnUnreachableCode:
		diagnostic.Severity = Warn
	default:
		return nil
	}
	return diagnostic
}
