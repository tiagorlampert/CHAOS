package ui

import "fmt"

func ShowMenu(port string) {
	fmt.Printf(`
 ┌───────────────────────────────────────────────────┐ 
 │                     CHAOS                         │ 
 │             http://127.0.0.1:%s│ 
 │                by tiagorlampert                   │ 
 └───────────────────────────────────────────────────┘
  `, fillSpace(port))
}

func fillSpace(v string) string {
	spaceSize := 20 - len(v)
	for i := 0; i <= spaceSize; i++ {
		v += " "
	}
	return v
}
