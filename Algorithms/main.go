package main

import (
	LL "Algorithms/LinkedList"
)

func main() {
	lk_list := LL.LinkedList{
		Head: nil,
	}

	lk_list.InsertNode(1)
	lk_list.InsertNode(3)
	lk_list.InsertNode(5)
	lk_list.InsertNode(7)

	lk_list.ReadList()
	println("--------------")
	lk_list.ReverseList()
	lk_list.ReadList()

	reversePhrase("Hola Mundo")
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	sumDiagonal(matrix)

	//a
	//ab
	//bc
	//abc
}

func sortArray(strArr []string) {
	sortedStr := make([]string, len(strArr))
	currentStr := sortedStr[0]

	for i := 1; i < len(strArr); i++ {
		modStr := len(strArr[i]) % 3
		modCurrentStr := len(currentStr) % 3

		if modStr > modCurrentStr {

		}
	}
}
