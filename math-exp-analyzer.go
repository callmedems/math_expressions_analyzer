package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	inputFile := flag.String("i", "", "Expressions.txt")
	outputFile := flag.String("o", "", "Results.txt")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Uso: go run math-exp-analyzer.go -i expressions.txt -o results.txt")
		return
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error al abrir archivo de entrada:", err)
		return
	}
	defer file.Close()

	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error al crear archivo de salida:", err)
		return
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		valid := validateExpression(line)
		result := "Invalid"
		if valid {
			result = "Valid"
		}
		outFile.WriteString(fmt.Sprintf("Expression %d: %-40s -  %s\n", lineNumber, line, result))
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}

func validateExpression(expr string) bool {
	expr = strings.ReplaceAll(expr, " ", "")
	if expr == "" {
		return false
	}

	stack := []rune{}
	prevToken := ""
	i := 0
	length := len(expr)

	for i < length {
		ch := rune(expr[i])

		// Manejar paréntesis y corchetes
		if ch == '(' || ch == '[' {
			// Verificar que no haya corchetes dentro de paréntesis
			if ch == '[' {
				for j := len(stack) - 1; j >= 0; j-- {
					if stack[j] == '(' {
						return false
					} else if stack[j] == '[' {
						break
					}
				}
			}
			stack = append(stack, ch)
			prevToken = string(ch)
			i++
			continue
		}

		if ch == ')' || ch == ']' {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if (ch == ')' && top != '(') || (ch == ']' && top != '[') {
				return false
			}
			stack = stack[:len(stack)-1]
			prevToken = string(ch)
			i++
			continue
		}

		// Manejar operadores
		if isOperator(ch) {
			// Manejar '**' para exponentiación
			if ch == '*' && i+1 < length && expr[i+1] == '*' {
				if prevToken == "" || isOperator(rune(prevToken[0])) || prevToken == "(" || prevToken == "[" {
					return false
				}
				i += 2
				prevToken = "**"
				continue
			}

			if prevToken == "" || isOperator(rune(prevToken[0])) || prevToken == "(" || prevToken == "[" {
				return false
			}
			i++
			prevToken = string(ch)
			continue
		}

		// Manejar números
		if unicode.IsDigit(ch) || ch == '.' || ch == '-' {
			start := i
			if ch == '-' {
				i++
				if i >= length || (!unicode.IsDigit(rune(expr[i])) && expr[i] != '.') {
					return false
				}
			}

			dotSeen := false
			for i < length && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
				if expr[i] == '.' {
					if dotSeen {
						return false
					}
					dotSeen = true
				}
				i++
			}

			number := expr[start:i]
			matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, number)
			if !matched {
				return false
			}

			// Verificar que no haya concatenación implícita con paréntesis o corchetes
			if i < length {
				nextCh := rune(expr[i])
				if nextCh == '(' || nextCh == '[' {
					return false
				}
			}
			prevToken = "number"
			continue
		}

		// Caracter inválido
		return false
	}

	return len(stack) == 0
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}
