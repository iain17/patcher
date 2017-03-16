package scanner

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"io"
	"container/list"
)

const joker = byte(0xff) // which is '?'

func Scan(sig string, filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return int64(0), err
	}
	defer file.Close()

	sequence := sigToBytes(sig)
	return find(sequence, file), nil
}

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
func sigToBytes(sig string) *list.List {
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