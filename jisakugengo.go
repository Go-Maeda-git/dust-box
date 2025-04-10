package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// トークンの種類
type TokenType int

const (
	TOKEN_EOF TokenType = iota
	TOKEN_IDENTIFIER
	TOKEN_NUMBER_FLOAT
	TOKEN_OPERATOR
	TOKEN_KEYWORD
	TOKEN_PUNCTUATOR // 記号 (括弧、中括弧、セミコロンなど)
)

// トークン構造体
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// 字句解析器
func lex(source string) ([]Token, error) {
	var tokens []Token
	lines := strings.Split(source, "\n")
	lineNumber := 1
	for _, line := range lines {
		columnNumber := 1
		for columnNumber <= len(line) {
			char := string(line[columnNumber-1])

			// 空白をスキップ
			if strings.Contains(" \t", char) {
				columnNumber++
				continue
			}

			// コメント (単一行)
			if strings.HasPrefix(line[columnNumber-1:], "//") {
				break // 行末までスキップ
			}

			// コメント (複数行)
			if strings.HasPrefix(line[columnNumber-1:], "/*") {
				endComment := strings.Index(line[columnNumber-1:], "*/")
				if endComment == -1 {
					// 複数行にわたるコメントの処理 (ここでは簡単のためエラーとする)
					return nil, fmt.Errorf("unterminated multi-line comment at line %d, column %d", lineNumber, columnNumber)
				}
				columnNumber += endComment + 2
				continue
			}

			// 数値 (浮動小数点数)
			if matched, _ := regexp.MatchString(`^[0-9]+\.[0-9]+`, line[columnNumber-1:]); matched {
				match := regexp.MustCompile(`^[0-9]+\.[0-9]+`).FindString(line[columnNumber-1:])
				tokens = append(tokens, Token{Type: TOKEN_NUMBER_FLOAT, Literal: match, Line: lineNumber, Column: columnNumber})
				columnNumber += len(match)
				continue
			}
			if matched, _ := regexp.MatchString(`^[0-9]+`, line[columnNumber-1:]); matched {
				match := regexp.MustCompile(`^[0-9]+`).FindString(line[columnNumber-1:])
				// ここでは一旦整数として扱う。型推論時に Float64 に変換を検討
				tokens = append(tokens, Token{Type: TOKEN_NUMBER_FLOAT, Literal: match + ".0", Line: lineNumber, Column: columnNumber})
				columnNumber += len(match)
				continue
			}

			// 識別子とキーワード
			if matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*`, line[columnNumber-1:]); matched {
				match := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`).FindString(line[columnNumber-1:])
				tokenType := TOKEN_IDENTIFIER
				switch match {
				case "if", "else", "for", "func", "class", "let", "return", "new": // 仮のキーワード
					tokenType = TOKEN_KEYWORD
				}
				tokens = append(tokens, Token{Type: tokenType, Literal: match, Line: lineNumber, Column: columnNumber})
				columnNumber += len(match)
				continue
			}

			// 演算子と記号 (暫定的な定義)
			operators := "+-*/=<>!."
			punctuators := "(){};,"
			if strings.Contains(operators, char) {
				tokens = append(tokens, Token{Type: TOKEN_OPERATOR, Literal: char, Line: lineNumber, Column: columnNumber})
				columnNumber++
				continue
			}
			if strings.Contains(punctuators, char) {
				tokens = append(tokens, Token{Type: TOKEN_PUNCTUATOR, Literal: char, Line: lineNumber, Column: columnNumber})
				columnNumber++
				continue
			}

			return nil, fmt.Errorf("unexpected character '%s' at line %d, column %d", char, lineNumber, columnNumber)
		}
		lineNumber++
	}
	tokens = append(tokens, Token{Type: TOKEN_EOF, Literal: "", Line: lineNumber, Column: 1})
	return tokens, nil
}

func main() {
	source := `
let pi = 3.14;
let radius = 2.0;
let area = pi * radius * radius; // Calculate area

if (area > 10.0) {
	print(area);
}

func circleArea(r) {
	let result = 3.14 * r * r;
	return result;
}

let myArea = circleArea(radius);
/*
This is a
multi-line comment.
*/
class Circle {
	let r;
	func init(radius) {
		this.r = radius;
	}
	func getArea() {
		return 3.14 * this.r * this.r;
	}
}

let c = new Circle(5.0);
let circle_area = c.getArea();
`

	tokens, err := lex(source)
	if err != nil {
		log.Fatal(err)
	}
	for _, tok := range tokens {
		fmt.Printf("%+v\n", tok)
	}
}