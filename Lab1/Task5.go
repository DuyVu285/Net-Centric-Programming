package main

import "fmt"

var openBrackets = map[rune]bool{'{': true, '(': true, '[': true}
var closeBrackets = map[rune]rune{'}': '{', ')': '(', ']': '['}

func Bracket(text string) bool {
	bks := newStack()
	for _, s := range text {
		if openBrackets[s] {
			bks.push(s)
		} else if b, ok := closeBrackets[s]; ok {
			if bks.isEmpty() || b != bks.pop() {
				return false
			}
		}
	}
	return bks.isEmpty()
}

type stack struct {
	element []rune
}

func newStack() *stack {
	return &stack{
		element: []rune{},
	}
}
func (s *stack) pop() rune {
	ln := len(s.element)
	res := s.element[ln-1]
	s.element = s.element[:ln-1]
	return res
}
func (s *stack) push(e rune) {
	s.element = append(s.element, e)
}
func (s *stack) isEmpty() bool {
	return len(s.element) == 0
}

func main() {
	string := "fmt.Println(a.TypeOf(xyz)){[ ]}"
	fmt.Println("The string is: ", string)
	if Bracket(string) {
		fmt.Println("Correct")
	} else {
		fmt.Println("Incorrect")
	}
}
