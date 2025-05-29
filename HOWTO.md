# Documentación del proyecto

Programa que simula un compialdor con el comportamiento de un automata de pila (PDA), que analiza la sintaxis de diferentes expresiones matemáticas. Apilando los corchetes y parentesís segun el orden de jerarquia por convención. Asi mismo balancea el uso de estos simbolos dentro de las cadenas. Donde se asegura la apertura y el cierre correspondiente de cada uno.

## Funcionamiento

El programa lee un archivo "expressions.txt" con la lista de expresiones que se analizaran y devolvera un archivo "results.txt" verificando si las cadenas son válidas o inválidas.

## Consideraciones
1. Toma como invalidas cadenas con operadores que no se defineron dentro del código.
2. Respeta la jerarquia entre parentesis y corchetes.
3. Operadores duplicados o mal colocados.
4. Parentesis o corchetes desbalanceados o vacíos,
5. Debe haber siempre digitos o parentesis/corchetes a los lados de cualquier operador excepto que sea un signo negativo.

## Instrucciones para compilar el programa:

Corre la línea en tu editor: `go run math_pda_compiler.go -i expressions.txt -o results.txt`
Esto incluye los archivos de entrada y salida que leerá el programa.


