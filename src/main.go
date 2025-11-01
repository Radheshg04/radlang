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
	fmt.Println(string(file))
	fmt.Println(Lex(string(file)))
}
