package libeplc

type Flag uint

const (
	_ Flag = 1 << iota
	Verbose
	Profile
	Default
)

func (flag Flag) Contains(f Flag) bool {
	return f & flag != 0
}

func (flag* Flag) Add(f Flag) Flag {
	return *flag | f
}