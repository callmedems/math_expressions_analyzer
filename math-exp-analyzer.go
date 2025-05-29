package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"unicode"
)

func main() {
	// para leer los archivos
	inputFile := flag.String("i", "", "Expressions.txt")
	outputFile := flag.String("o", "", "Results.txt")
	flag.Parse()

	// verificamos que ambos archivos hayan sido especificados
	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Uso: go run math-exp-analyzer.go -i expressions.txt -o results.txt")
		return
	}

	// abrir el archivo de entrada para leer
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error al abrir archivo de entrada:", err)
		return
	}
	defer file.Close()

	// crear o sobrescribir el archivo de salida
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
		outFile.WriteString(fmt.Sprintf("Expression %d: %-20s -  %s\n", lineNumber, line, result))
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}

// validateExpression valida una expresión matemática
func validateExpression(expr string) bool {
	stack := []rune{}
	i := 0
	length := len(expr)

	for i < length {
		ch := rune(expr[i])

		// Saltar espacios
		if unicode.IsSpace(ch) {
			i++
			continue
		}

		// Abrir paréntesis o corchetes
		if ch == '(' || ch == '[' {
			stack = append(stack, ch)
			i++
			continue
		}

		// Cerrar paréntesis o corchetes
		if ch == ')' || ch == ']' {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if (ch == ')' && top != '(') || (ch == ']' && top != '[') {
				return false
			}
			stack = stack[:len(stack)-1]
			i++
			continue
		}

		// Validar operadores
		if isOperator(ch) {
			// Checar doble operador **
			if ch == '*' && i+1 < length && expr[i+1] == '*' {
				i += 2
				continue
			}

			// Validar que no sea operador al final o duplicado
			if i == length-1 || isOperator(rune(expr[i+1])) {
				return false
			}
			i++
			continue
		}

		// Validar número (real o entero)
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

			// Verificar si es un número válido con regex
			number := expr[start:i]
			matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, number)
			if !matched {
				return false
			}
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
