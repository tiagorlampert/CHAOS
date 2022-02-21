package ui

import (
	"fmt"
	"unicode/utf8"
)

const spaceLength = 60

func ShowMenu(version, host, port string) {
	if len(port) > 0 {
		port = fmt.Sprint(":", port)
	}

	fmt.Printf(`
 ┌%s┐ 
 │%s│ 
 │%s│ 
 └%s┘
`,
		fillSpace("", "─"),
		fillSpace(fmt.Sprintf("CHAOS (%s)", version), " "),
		fillSpace(host+port, " "),
		fillSpace("", "─"),
	)
}

func fillSpace(content, filler string) string {
	spaceToFillSize := spaceLength - utf8.RuneCountInString(content)
	spaceBothSide := spaceToFillSize / 2

	var finalStr string
	for i := 0; i < spaceBothSide; i++ {
		finalStr += filler
	}
	finalStr += content
	for i := 0; i < spaceBothSide; i++ {
		finalStr += filler
	}

	finalStrCount := utf8.RuneCountInString(finalStr)
	if finalStrCount < spaceLength {
		diff := spaceLength - finalStrCount
		for i := 0; i < diff; i++ {
			finalStr += filler
		}
	}
	return finalStr
}
