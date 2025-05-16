package main

import "fmt"

func CallAlgorithm() {
	CallAlgorithm3()
}

func reverse(text string) {
	reverseText := ""

	for i := len(text) - 1; i >= 0; i-- {
		reverseText = reverseText + string(text[i])
	}
	fmt.Printf("Text: %s Reverse: %s", text, reverseText)
}

func palindromo(text string) {
	reverseText := ""

	for i := len(text) - 1; i >= 0; i-- {
		reverseText = reverseText + string(text[i])
	}

	if text == reverseText {
		fmt.Println("Si es palindromo")
	} else {
		fmt.Println("No es palindromo")
	}
}

func maxandmin(numbers []int) (int, int) {
	for i := 0; i < len(numbers); i++ {
		temp1 := numbers[i]
		for j := 1; j < len(numbers); j++ {
			temp2 := numbers[j]
			if temp1 > temp2 {
				numbers[i] = temp2
				numbers[j] = temp1
			}
		}
	}
	return numbers[0], numbers[len(numbers)-1]
}
