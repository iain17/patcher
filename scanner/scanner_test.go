package scanner

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"container/list"
)

//Test string sig to linked list byte data structure.
func TestSigToBytes(t *testing.T) {
	expect := list.New()
	expect.PushBack(0xC7)
	expect.PushBack(0xE8)
	expect.PushBack(0xff)
	expect.PushBack(0xff)
	expect.PushBack(0xff)
	expect.PushBack(0xff)
	expect.PushBack(0x48)
	expect.PushBack(0x89)

	input := "C7 E8 ? ? ? ? 48 89"
	assert.Equal(t, sigToBytes(input), expect)
}

//Integration test with the atom executable.
//Find the main
func TestScan(t *testing.T) {
	/*
0000000100000f00         push       rbp
0000000100000f01         mov        rbp, rsp
0000000100000f04         push       r14
0000000100000f06         push       rbx
0000000100000f07         mov        r14, rsi
0000000100000f0a         mov        ebx, edi
0000000100000f0c         call       sub_100000f30
0000000100000f11         mov        edi, ebx
	 */
	find := "55 48 89 E5 41 56 53 49 89 F6 89 FB E8 ? ? ? ? 89 DF 4C"
	expected := int64(3840)
	address, err := Scan(find, "../tests/Atom")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, address, expected)
}