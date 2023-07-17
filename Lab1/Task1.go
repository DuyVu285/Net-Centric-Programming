package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func compareDNA(a, b []string) int {
	hamming_distance := 0
	for i, v := range a {
		if v != b[i] {
			hamming_distance += 1
		}
	}
	return hamming_distance
}

func randomString(n int, alphabet []rune) string {
	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var alphabet []rune = []rune("ACGT")

	for i := 0; i < 1000; i++ {
		rs1 := randomString(17, alphabet)
		rs2 := randomString(17, alphabet)
		fmt.Println("The DNA pair number", i+1, " is: ")
		fmt.Println(rs1)
		fmt.Println(rs2)
		res1 := strings.Split(rs1, "")
		res2 := strings.Split(rs2, "")
		fmt.Println("The Hamming Distance of this pair: ")
		fmt.Println(compareDNA(res1, res2))
	}
}
