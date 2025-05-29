package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// para leer los archivos
	inputFile := flag.String("i", "", "Expressions.txt")
	outputFile := flag.String("o", "", "Results.txt")
	flag.Parse()

	input, err := os.Open(*inputFile)
	output, err = os.Create(*outputFile)

	scanner := bufio.NewScanener(input)
	writer := bufio.NewReadWriter(output)
	expresionNum := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		valid := validarExpresion(line)
		status := "Valid"
		if valid {
			status = "Valid"
		}

		writer.WriteString(fmt.Sprintf("Expresion %d: %s - %s\n", expresionNum, line, status))

		writer.Flush()
	}
}

func validarExpresion(expr string) bool {
	// usamos una pila para checar paréntesis y corchetes balaceads
	stack := []rune{}
	i := 0 //para iterar caracteres
	length := len(expr)
}
