package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Errorf("引数の個数が正しくありません\n")
	}
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	fmt.Printf("  mov rax, %d\n", os.Args[1][0]-'0')
	for i := 1; i < len(os.Args[1]); i++ {
		p := os.Args[1][i]
		if p == '+' {
			i++
			fmt.Printf("  add rax, %d\n", os.Args[1][i]-'0')
			continue
		}
		if p == '-' {
			i++
			fmt.Printf("  sub rax, %d\n", os.Args[1][i]-'0')
			continue
		}
		fmt.Errorf("予期しない文字です: '%c'\n", p)
		os.Exit(1)
	}
	fmt.Printf("  ret\n")
}
