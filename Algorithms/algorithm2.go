package main

import (
	"fmt"
)

func reverseString(text string) string {
	byteSlice := []rune(text)
	left := 0
	for i := len(byteSlice) - 1; i >= left; i-- {
		temp := byteSlice[i]
		byteSlice[i] = byteSlice[left]
		byteSlice[left] = temp
		left++
	}
	return string(byteSlice)
}

func palindromo2(text string) {
	if text == reverseString(text) {
		fmt.Println("Es palindromo")
	} else {
		fmt.Println("No es palindromo")
	}
}

func getminandmax(numbers []int) (int, int) {
	max := numbers[0]
	min := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i] > max {
			max = numbers[i]
		} else if numbers[i] < min {
			min = numbers[i]
		}
	}
	return min, max
}
