package util

import "strings"

var (
	DelimiterString string = "\n"
	DelimiterByte   byte   = '\n'
)

func SplitAfterIndex(str string, index byte) string {
	return str[strings.IndexByte(str, index)+1:]
}
