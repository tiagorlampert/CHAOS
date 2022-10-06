package ui

import (
	"fmt"
	"time"
	"unicode/utf8"
)

const spaceLength = 60

func ShowMenu(version, port string) {
	fmt.Printf(`
 ┌%s┐ 
 │%s│ 
 │%s│ 
 │%s│ 
 └%s┘
`,
		fillSpace("", "─"),
		fillSpace(fmt.Sprintf("CHAOS (%s)", version), " "),
		fillSpace("http://127.0.0.1:"+port, " "),
		fillSpace(fmt.Sprintf("by tiagorlampert (%d)", time.Now().UTC().Year()), " "),
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
