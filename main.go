package main

import (
	"fmt"
	"os"
)

func isdigit(ch byte) bool {
	if '0' <= ch && '9' >= ch {
		return true
	}
	return false
}
func strtol(str string, index int) ([]byte, int) {
	var num []byte
	for i := index; i < len(str); i++ {
		if isdigit(str[i]) {
			num = append(num, str[i])
		} else {
			return num, i - 1
		}
	}
	return num, len(str) - 1
}
func isspace(ch byte) bool {
	if ch == ' ' {
		return true
	}
	return false
}
func main() {
	if len(os.Args) != 2 {
		fmt.Errorf("引数の個数が正しくありません\n")
	}
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	s, index := strtol(os.Args[1], 0)
	index++
	fmt.Printf("  mov rax, %s \n", s)
	for i := index; i < len(os.Args[1])-1; i++ {
		p := os.Args[1][i]
		if p == '+' {
			i++
			s, index := strtol(os.Args[1], i)
			i = index
			fmt.Printf("  add rax, %s\n", s)
			continue
		}
		if p == '-' {
			i++
			s, index := strtol(os.Args[1], i)
			i = index
			fmt.Printf("  sub rax, %s\n", s)
			continue
		}
		fmt.Errorf("予期しない文字です: '%c'\n", p)
		os.Exit(1)
	}
	fmt.Printf("  ret\n")
}
