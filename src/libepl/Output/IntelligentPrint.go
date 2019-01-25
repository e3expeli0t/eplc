package Output

import (
	"eplc/src/libepl/Output/color"
	"fmt"
	"strings"
)


type LexicalErrorDescriptor struct {
	fname string
	line uint
	lineOffset uint
	currentLine string
	errorMSG string
	ch rune
}

//Todo: Make this more elegant
func LexicalPrint(fname string, line uint,  lineOffset uint, currentLine string, errorMSG string, ch rune) {
	ed := LexicalErrorDescriptor{fname, line, lineOffset,
		currentLine, errorMSG, ch}
	LexerIntelligentError(ed)
}

type ErrorDescriptor struct {
	fname string
	line uint
	lineOffset uint
	currentLine string
	errorMSG string
	token string
}

func (e *ErrorDescriptor) basicInfoPrinter() string{
	return color.BGreen(fmt.Sprintf(":%s:%d:%d: ", e.fname, e.line, e.lineOffset))
}

func (e *LexicalErrorDescriptor) basicInfoPrinter() string{
	return color.BGreen(fmt.Sprintf(":%s:%d:%d: ", e.fname, e.line, e.lineOffset))
}

func prints(s string , n int) string{
	return strings.Repeat(s, n)
}

//Todo: make it support Cflags (there token don't have the @ char)
func (e *ErrorDescriptor) TokenMarker() string {
	str := e.currentLine+"\n"

	for _, token := range strings.Split(str, " ") {
		if token != e.token {
			str += prints(color.BLightPurple("~"), len(token))
		} else {
			str += prints(color.BLightPurple("^"), len(token))
		}
	}

	//pad
	str += prints(color.BLightPurple("~"), len(e.currentLine)-len(str))

	return str
}

func (e *LexicalErrorDescriptor) LexicalMarker() string {
	str := e.currentLine+"\n"

	for _, token := range e.currentLine{
		if token != e.ch {
			str += prints(color.BLightPurple("~"), 1)
		} else {
			str += prints(color.BLightPurple("^"), 1)
		}
	}

	if len(str) < len(e.currentLine) {
		str += prints(color.BLightPurple("~"), len(e.currentLine)-len(str))
	}

	return str
}

func LexerIntelligentError(err LexicalErrorDescriptor) {
	fmt.Print(color.BGreen("Error:")+err.basicInfoPrinter())
	fmt.Print(color.BGreen(err.errorMSG)+"\n")
	fmt.Println(err.LexicalMarker())
}

func ParserIntelligentError(err ErrorDescriptor) {
	fmt.Print(color.BGreen("Error:")+err.basicInfoPrinter())
	fmt.Println(color.BGreen(err.errorMSG)+"\n")
	fmt.Println(err.TokenMarker())
}
