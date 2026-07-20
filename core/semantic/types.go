package semantic

type ValueType int

const (
	InvalidType ValueType = iota
	IntType
	FloatType
	StringType
	BoolType
	ErrorType
)
