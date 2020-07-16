package analysis

import "fmt"

type ErrorLevel uint

const (
	Major ErrorLevel = iota
	Minor ErrorLevel = iota
	Fatal ErrorLevel = iota
)

func (l ErrorLevel) ToString() string {
	switch l {
	case Fatal:
		return "Fatal"
	case Major:
		return "Major"
	case Minor:
		return "Minor"
	default:
		return ""
	}
}

type Error interface {
	IsFatal() bool
	ToString() string
}

func NewError(level ErrorLevel, msg string) *TypeError {
	return &TypeError{
		Descriptor: msg,
		Level:      level,
	}
}
type TypeError struct {
	Line       int
	Offset     int
	Descriptor string
	Level ErrorLevel
}

func (t TypeError) IsFatal() bool {
	return t.Level != Fatal
}
//as of now the compiler don't keep track of line numbers
func (t TypeError) ToString() string {
	return fmt.Sprintf("%s Error: %s", t.Level.ToString(), t.Descriptor)
}
