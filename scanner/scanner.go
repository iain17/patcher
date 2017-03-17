package scanner

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"io"
	"container/list"
	"github.com/pkg/errors"
)

const joker = byte(0xff) // which is '?'

//TODO: (optimization) Allow an array of signatures and multiple results. So that we don't iterate through the whole file for every signature.
func Scan(sig string, filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return int64(0), err
	}
	defer file.Close()

	sequence := ByteSeqToByteList(sig)
	return find(sequence, file), nil
}

func GenSig(address int64, filePath string, length int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sig := gen(address, length, file)
	if len(sig) == 0 {
		return nil, errors.New("Could not find address")
	}
	return sig, nil
}

//Generates a sinature based on a memory addres
func gen(address int64, length int, rd io.Reader) []byte {
	signature := []byte{}
	reader := bufio.NewReader(rd)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)
	memory := int64(0)

	for scanner.Scan() {
		//Up the memory address
		memory += int64( len(scanner.Bytes()) )
		if memory > address {
			for _, byte := range scanner.Bytes() {
				signature = append(signature, byte)
			}
		}
		//Signature should not be bigger than 8
		if len(signature) >= length {
			break
		}
	}
	return signature
}

//Finds a signature in binary.
func find(find *list.List, rd io.Reader) int64 {
	reader := bufio.NewReader(rd)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)
	memory := int64(0)
	curr := find.Front()
	foundMemory := int64(0)//the memory the first time a instruction was right

	for scanner.Scan() {
		//Up the memory address
		memory += int64( len(scanner.Bytes()) )

		//hexstring := fmt.Sprintf("%#08x = %x", memory, scanner.Text())
		instruction := scanner.Bytes()[0]//We have scan type set to bytes. We can't have more than one!

		//If the current byte we are iterating through is the correct. Go to the next
		if curr.Value.(byte) == instruction || curr.Value.(byte) == joker {
			curr = curr.Next()
		} else {
		//If it isn't go back to the beginning.
			curr = find.Front()
			foundMemory = memory
		}

		//Found it!
		if curr == nil {
			return foundMemory
		}
	}
	return 0//Couldn't find it!
}

//Returns bytes based on the signature string.
func ByteSeqToByteList(sig string) *list.List {
	parts := strings.Split(sig, " ")
	bytes := list.New()
	for _, part := range parts {
		if part == "?" {
			bytes.PushBack(joker)
		} else {
			val, err := strconv.ParseInt(part, 16,64)
			if err != nil {
				panic(err)
			}
			bytes.PushBack(byte(val))
		}
	}
	return bytes
}