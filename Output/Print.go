package Output

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func PrintLog(log ...interface{}) {
	fmt.Println(aurora.Green("[*]"), log)
}
