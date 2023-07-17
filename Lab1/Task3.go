package main

import (
	"fmt"
	"strconv"
	"strings"
)

func luhnAlgo(word string) bool {

	stringNum := strings.ReplaceAll(word, " ", "")
	number, _ := strconv.Atoi(stringNum)

	if checksum(number) == 0 {
		return true
	} else {
		return false
	}
}

func checksum(number int) int {
	var luhn int

	for i := count(number) - 1; number > 0; i-- {
		cur := number % 10
		if i%2 == 0 {
			cur = cur * 2
			if cur >= 9 {
				cur = cur - 9
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func count(number int) int {
	count := 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count
}

func main() {
	number1 := "4539 3195 0343 6467"
	fmt.Println("The string is: ", number1)
	if luhnAlgo(number1) == true {
		fmt.Println("The string is Valid!")
	} else {
		fmt.Println("The string is Invalid!")
	}
	number2 := "8273 1232 7352 0569"
	fmt.Println("The string is: ", number2)
	if luhnAlgo(number2) == true {
		fmt.Println("The string is Valid!")
	} else {
		fmt.Println("The string is Invalid!")
	}
}
