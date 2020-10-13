package util

import "strings"

var (
	DelimiterString string = "\n"
	DelimiterByte   byte   = '\n'
)

func SplitAfterIndex(str string, index byte) string {
	return str[strings.IndexByte(str, index)+1:]
}

func Contains(v []string, str string) bool {
	var has bool
	for _, param := range v {
		if strings.Contains(param, str) {
			has = true
			break
		}
	}
	return has
}

func Find(v []string, str string) string {
	for _, param := range v {
		if strings.Contains(param, str) {
			return param
		}
	}
	return ""
}
