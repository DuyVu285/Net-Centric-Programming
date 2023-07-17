package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func createBoard(length, width, mines int) [][]string {
	m := make([][]string, width)
	for i := 0; i < width; i++ {
		for j := 0; j < length; j++ {
			m[i] = append(m[i], ".")
		}
	}

	for i := 0; i < mines; i++ {
		m[rand.Intn(width)][rand.Intn(length)] = "*"
	}
	return m
}

func Annotate(board [][]string) [][]string {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			board[i][j] = getValue(i, j, board)
		}
	}
	return board
}

func getValue(rowIndex, colIndex int, matrix [][]string) string {
	result := " "
	count := 0
	if matrix[rowIndex][colIndex] != "*" {
		if colIndex-1 >= 0 && matrix[rowIndex][colIndex-1] == "*" {
			count++
		}
		if colIndex+1 <= len(matrix[rowIndex])-1 && matrix[rowIndex][colIndex+1] == "*" {
			count++
		}
		if rowIndex-1 >= 0 && matrix[rowIndex-1][colIndex] == "*" {
			count++
		}
		if rowIndex+1 <= len(matrix)-1 && matrix[rowIndex+1][colIndex] == "*" {
			count++
		}
		if rowIndex-1 >= 0 && colIndex-1 >= 0 && matrix[rowIndex-1][colIndex-1] == "*" {
			count++
		}
		if rowIndex-1 >= 0 && colIndex+1 <= len(matrix[rowIndex])-1 && matrix[rowIndex-1][colIndex+1] == "*" {
			count++
		}
		if rowIndex+1 <= len(matrix)-1 && colIndex-1 >= 0 && matrix[rowIndex+1][colIndex-1] == "*" {
			count++
		}
		if rowIndex+1 <= len(matrix)-1 && colIndex+1 <= len(matrix[rowIndex])-1 && matrix[rowIndex+1][colIndex+1] == "*" {
			count++
		}
	} else {
		return "*"
	}
	if count > 0 {
		result = strconv.Itoa(count)
	}
	return result
}
func printGrid(board [][]string) {
	for row := 0; row < len(board); row++ {
		for column := 0; column < len(board[0]); column++ {
			fmt.Print(board[row][column], " ")
		}
		fmt.Print("\n")
	}
}
func main() {
	m := createBoard(20, 25, 99)
	fmt.Print("The Original Matrix: \n")
	printGrid(m)
	fmt.Print("\n")
	fmt.Print("After Narking The Board: \n")
	printGrid(Annotate(m))
}
