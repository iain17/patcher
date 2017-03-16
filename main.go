package main

import (
	"github.com/iain17/cracker/scanner"
	"fmt"
)

func main() {
	find := "C7 E8 DC 78 12 00 48 89"//6624
	address := scanner.Scan(find)
	fmt.Printf("signature is at %#08x (%d)\n", address, address)
}
