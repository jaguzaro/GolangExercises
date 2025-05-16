package main

import (
	"fmt"
	"strings"
)

func CallAlgorithm3() {
	fmt.Println(anagrama("amor", "roma"))
}

func anagrama(str1, str2 string) bool {
	if len(str1) != len(str2) {
		return false
	}

	original := make(map[rune]int)

	for i := 0; i < len(str1); i++ {
		original[rune(str1[i])]++
		original[rune(str2[i])]--
	}

	for _, value := range original {
		if value != 0 {
			return false
		}
	}
	return true
}

func reversePhrase(phrase string) {
	var b strings.Builder
	words := strings.Split(phrase, " ")

	for _, w := range words {
		b.WriteString(reverseString(w) + " ")
	}
	fmt.Println(b.String())
}

func sumDiagonal(matrix [][]int) {
	fmt.Println(len(matrix), matrix)
	var sum int
	for i := 0; i < len(matrix); i++ {
		sum += matrix[i][i]
	}
	fmt.Println(sum)
}
