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
	inputFile := flag.String("i", "", "Expressions.txt")
	outputFile := flag.String("o", "", "Results.txt")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		return
	}
	// Manejo de errores para abrir los archivos recomendados por copilot.
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

// validateExpression valida una expresión matemática
func validateExpression(expr string) bool {
	stack := []rune{}
	prevToken := ""
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
			// Verificar que no haya corchete dentro de paréntesis
			if ch == '[' && len(stack) > 0 && stack[len(stack)-1] == '(' {
				return false
			}
			// Verificar que no haya número o cierre de paréntesis/corchete antes sin operador
			if prevToken == "number" || prevToken == ")" || prevToken == "]" {
				return false
			}
			stack = append(stack, ch)
			prevToken = string(ch)
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
			// Verificar que no haya número o cierre de paréntesis/corchete después sin operador
			if i+1 < length {
				nextCh := rune(expr[i+1])
				if unicode.IsDigit(nextCh) || nextCh == '(' || nextCh == '[' {
					return false
				}
			}
			prevToken = string(ch)
			i++
			continue
		}

		// Validar operadores
		if isOperator(ch) {
			// Checar operador potencia **
			if ch == '*' && i+1 < length && expr[i+1] == '*' {
				// Verificar que no esté al final
				if i+2 >= length {
					return false
				}
				next := rune(expr[i+2])
				if isOperator(next) || next == ')' || next == ']' {
					return false
				}
				// Asegura que antes haya número o cierre de paréntesis o corchete
				if prevToken != "number" && prevToken != ")" && prevToken != "]" {
					return false
				}
				prevToken = "operator"
				i += 2
				continue
			}

			// Verificar que no esté al principio o al final
			if i == 0 || i == length-1 {
				return false
			}

			// Verificar que antes haya número o cierre de paréntesis o corchete
			if prevToken != "number" && prevToken != ")" && prevToken != "]" {
				return false
			}

			// Verificar que después haya un número o apertura de paréntesis/corchete
			j := i + 1
			for j < length && unicode.IsSpace(rune(expr[j])) {
				j++
			}
			if j >= length {
				return false
			}
			next := rune(expr[j])
			if !(unicode.IsDigit(next) || next == '(' || next == '[' || next == '-') {
				return false
			}

			prevToken = "operator"
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

			dotSeen := false // se valida el uso de un solo punto despues de un num
			for i < length && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
				if expr[i] == '.' {
					if dotSeen {
						return false
					}
					dotSeen = true
				}
				i++
			}

			// Verificar si es un numero válido con regex
			number := expr[start:i]
			matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, number)
			if !matched {
				return false
			}

			// Checa que no haya paréntesis o corchete después sin operador
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

func isOperator(r rune) bool { //booleano con operadores váalidos
	return r == '+' || r == '-' || r == '*' || r == '/'
}
