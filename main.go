package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	input := "$5\r\nAhmet\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, _ := reader.ReadByte()

	if b != '$' {
		fmt.Println("Invalid Type Expected bulk string")
	}
	size, _ := reader.ReadByte()
	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println("Name is", string(name))
}
