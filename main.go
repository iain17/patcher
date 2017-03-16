package main

import (
	"github.com/iain17/patcher/scanner"
	"fmt"
)

func main() {
	find := "71 92 ? ? ? ? AF 2F 41 B3 F8 FB E3 70 CD E6"//6624
	address, err := scanner.Scan(find, "tests/Atom")
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature is at %#08x (%d)\n", address, address)
}

//patched, err := os.Create("tests/Bee_pathed")
//
//if err != nil {
//	panic(err.Error())
//}
//
//defer patched.Close()
//
//reader := bufio.NewReader(file)
//scanner := bufio.NewScanner(reader)
//scanner.Split(bufio.ScanBytes)
//
//memory := int64(0)
//for scanner.Scan() {
//	hexstring := fmt.Sprintf("%#08x = %x", memory, scanner.Text())
//	println(hexstring)
//
//	patched.Write(scanner.Bytes())
//
//	//Up the memory address
//	memory += int64( len(scanner.Bytes()) )
//}
//return memory