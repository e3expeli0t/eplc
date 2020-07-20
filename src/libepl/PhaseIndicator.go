package libepl


type PhaseIndicator uint

const (
	_ = iota
	Parser
	Lexer
	IntermediateCodeGenerator
	TypeChecker
	Analysis
	TargetCodeGenerator
)
