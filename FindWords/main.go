package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	content := `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
	Integer fermentum auctor ornare. Mauris vestibulum sit amet nisl non tincidunt. 
	Quisque lorem est, rhoncus' eu fringilla vel, viverra sit amet augue. 
	Aenean feugiat ultrices orci, sed tempor arcu faucibus quis.`

	wordCount := getWordCount(content)
	words := make([]string, 0, wordCount)
	var sb strings.Builder
	for _, c := range content {
		if unicode.IsLetter(c) {
			sb.WriteRune(c)
		} else if unicode.IsSpace(c) {
			if sb.Len() > 0 {
				words = append(words, sb.String())
				sb.Reset()
			}
		}
	}

	if sb.Len() > 0 {
		words = append(words, sb.String())
	}

	fmt.Printf("%q", words)
}

func getWordCount(content string) int {
	count := 0
	hasLetter := false
	for _, c := range content {
		if unicode.IsLetter(c) {
			hasLetter = true
		} else if unicode.IsSpace(c) {
			if hasLetter {
				count++
				hasLetter = false
			}
		}
	}

	if hasLetter {
		count++
	}

	return count
}
