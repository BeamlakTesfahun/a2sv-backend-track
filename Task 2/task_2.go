package main

import (
	"fmt"
	"regexp"
	"strings"
)

func WordFrequency(s string) map[string]int {
	s = strings.ToLower(s)

	// remove punctuation
	re := regexp.MustCompile(`[^\w\s]`)
	s = re.ReplaceAllString(s, "")

	words := strings.Fields(s)
	counts := make(map[string]int)

	for _, w := range words {
		counts[w]++
	}

	return counts
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)

	re := regexp.MustCompile(`[^a-z0-9]`)
	s = re.ReplaceAllString(s, "")

	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	fmt.Println(WordFrequency("Hello, hello! World world..."))
	fmt.Println(IsPalindrome("never odd or even"))
	fmt.Println(IsPalindrome("hello"))
}
