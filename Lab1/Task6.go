package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func checkCharacter(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func checkContain(text string) {
	lowerText := strings.ToLower(text)
	textArray := strings.Split(lowerText, "")
	var checkContain []string
	for i := 0; i < len(text); i++ {
		if !(checkCharacter(checkContain, textArray[i])) {
			checkContain = append(checkContain, textArray[i])
		}
	}
	for i := 0; i < len(checkContain); i++ {
		go worker(checkContain[i], lowerText)
	}
	time.Sleep(time.Second)
}

func worker(character, text string) {
	res := strings.Count(text, character)
	fmt.Println(character, ": ", res)
}

func main() {
	text := "Net-centric Programming"
	fmt.Println("The string given here is:", text)
	checkContain(text)
	fmt.Println("\n")
	content, err := os.ReadFile("Lab1/Task6text.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
	checkContain(string(content))
}
