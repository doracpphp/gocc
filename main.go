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
type NodeKind int

const (
	NDADD NodeKind = iota
	NDSUB          // -
	NDMUL          // *
	NDDIV          // /
	NDNUM          // 整数
)

type Node struct {
	kind NodeKind
	lhs  *Node
	rhs  *Node
	val  int
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	node := &Node{kind: kind, lhs: lhs, rhs: rhs}
	return node
}
func NewNodeNum(val int) *Node {
	node := &Node{kind: NDNUM, val: val}
	return node
}
func expr(token *Token) (*Node, *Token) {
	node, token := mul(token)
	for {
		if ok, tok := consume('+', token); ok {
			rnode, tok := mul(tok)
			node = NewNode(NDADD, node, rnode)
			token = tok
		} else if ok, tok := consume('-', token); ok {
			rnode, tok := mul(tok)
			node = NewNode(NDSUB, node, rnode)
			token = tok
		} else {
			return node, token
		}
	}

}
func mul(token *Token) (*Node, *Token) {
	node, token := primary(token)
	for {
		if ok, tok := consume('*', token); ok {
			rnode, tok := primary(tok)
			node = NewNode(NDMUL, node, rnode)
			token = tok
		} else if ok, tok := consume('/', token); ok {
			rnode, tok := primary(tok)
			node = NewNode(NDDIV, node, rnode)
			token = tok
		} else {
			return node, token
		}
	}
}
func primary(token *Token) (*Node, *Token) {
	if ok, tok := consume('(', token); ok {
		node, token := expr(tok)
		token = expect(')', token)
		return node, token
	}
	n, token := expect_number(token)
	return NewNodeNum(n), token
}
func gen(node *Node) {
	if node.kind == NDNUM {
		fmt.Printf("  push %d\n", node.val)
		return
	}
	gen(node.lhs)
	gen(node.rhs)
	fmt.Printf("  pop rdi\n")
	fmt.Printf("  pop rax\n")
	switch node.kind {
	case NDADD:
		fmt.Printf("  add rax, rdi\n")
		break
	case NDSUB:
		fmt.Printf("  sub rax, rdi\n")
		break
	case NDMUL:
		fmt.Printf("  imul rax, rdi\n")
		break
	case NDDIV:
		fmt.Printf("  cqo\n")
		fmt.Printf("  idiv rdi\n")
		break
	}
	fmt.Printf("  push rax\n")
}
func consume(op byte, tok *Token) (bool, *Token) {
	if tok.kind != TKRESERVED || tok.str[0] != op {
		return false, tok
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
		if p == '+' || p == '-' || p == '*' || p == '/' || p == '(' || p == ')' {
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
	node, token := expr(token)
	fmt.Printf(".intel_syntax noprefix\n")
	fmt.Printf(".global main\n")
	fmt.Printf("main:\n")
	gen(node)
	fmt.Printf("  pop rax\n")
	fmt.Printf("  ret\n")
}
