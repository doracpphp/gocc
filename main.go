package main

import (
	"fmt"
	"os"
	"strconv"
)

//TokenKind TokenEnum
type TokenKind int

const (
	TKRESERVED TokenKind = iota // 記号
	TKNUM                       // 整数トークン
	TKEOF                       // 入力の終わりを表すトークン
)

type Token struct {
	kind TokenKind
	next *Token
	val  int
	str  []byte
}

func consume(op byte, tok *Token) (bool, *Token) {
	if tok.kind != TKRESERVED || tok.str[0] != op {
		return false, tok.next
	}
	return true, tok.next
}
func expect(op byte, tok *Token) *Token {
	if tok.kind != TKRESERVED || tok.str[0] != op {
		panic("ではありません")
	}
	return tok.next
}
func expect_number(tok *Token) (int, *Token) {
	if tok.kind != TKNUM {
		panic("数字ではありません")
	}
	var val int = tok.val
	return val, tok.next
}
func at_eof(tok *Token) bool {
	return tok.kind == TKEOF
}
func NewToken(kind TokenKind, cur *Token, str []byte) *Token {
	tok := &Token{kind: kind, str: str}
	cur.next = tok
	return tok
}
func tokenize(str string) *Token {
	var head Token
	head.next = nil
	cur := &head
	for i := 0; i < len(str); i++ {
		p := str[i]
		if isspace(p) {
			i++
			continue
		}
		if p == '+' || p == '-' {
			cur = NewToken(TKRESERVED, cur, []byte{p})
			continue
		}
		if isdigit(p) {
			s, index := strtol(str, i)
			i = index
			cur = NewToken(TKNUM, cur, s)
			n, err := strconv.Atoi(string(s))
			if err != nil {
				panic(err)
			}
			cur.val = n
			continue
		}
	}
	cur = NewToken(TKEOF, cur, []byte{})
	return head.next
}
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
	token := tokenize(os.Args[1])
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	n, token := expect_number(token)
	fmt.Printf("  mov rax, %d\n", n)
	for !at_eof(token) {
		if ok, tok := consume('+', token); ok {
			n, token = expect_number(tok)
			fmt.Printf("  add rax, %d\n", n)
			continue
		}
		tok := expect('-', token)
		n, token = expect_number(tok)
		fmt.Printf("  sub rax, %d\n", n)
	}
	fmt.Printf("  ret\n")
}
