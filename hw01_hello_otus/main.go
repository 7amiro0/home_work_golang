package main

import "fmt"

func main() {
	var text string = "Hello, OTUS!"
	runeText := []rune(text)
	for i, j := 0, len(runeText)-1; i < len(runeText)/2; i, j = i+1, j-1 {
		runeText[i], runeText[j] = runeText[j], runeText[i]
	}
	fmt.Println(string(runeText))
}
