package Output

import (
	"fmt"
	"eplc/src/libepl/Output/color"
	"os"
)

const startmsg = "Starting aplc v0.001:"

func PrintLog(log ...interface{}) {
	fmt.Print(color.BLightGreen("[*] "))
	fmt.Println(log...)
}

func PrintErr(seg string, msg...interface{}) {
	fmt.Print(color.BLightRed("Fatal: "), "<"+seg+">")
	fmt.Print(msg...)
	fmt.Println()
	fmt.Println(color.BLightBlue("Quiting!!"))
	os.Exit(-1)
}

func PrintStartMSG() {
	fmt.Println(color.LightBlue(startmsg))
}