package Output

import (
	"fmt"
	"./color"
)

func PrintLog(log ...interface{}) {
	fmt.Println(color.LightGreen("[*]"), log)
}
