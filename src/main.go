package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("../tests/hello.rad")
	if err != nil {
		fmt.Println(err.Error()) // Handle error
	}
	fmt.Printf("%s\n", string(file))
	fmt.Printf("Lexer Output: \n")
	tokens := Lex(string(file))
	for _, val := range tokens {
		fmt.Println(val)
	}
}
