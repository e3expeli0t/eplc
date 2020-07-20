package libepl

func NewInfoStruct (filename string) InfoStruct {
	return InfoStruct{
		Filename:  filename,
		SystemCPU: ResolveBits(),
	}
}

//todo: temporary solution
func ResolveBits() int {
	const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

	if is64Bit {
		return 64
	} else {
		return 32
	}
}

type InfoStruct struct {
	Filename string
	SystemCPU int
}
