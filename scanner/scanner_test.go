package scanner

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"container/list"
)

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