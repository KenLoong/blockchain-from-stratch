package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)
}

func TestVM(t *testing.T) {
	// 1 + 2 = 3
	// 1
	// push stack
	// 2
	// push stack
	// add
	// 3
	// push stack

	//	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	//              2    push   2   push    add
	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, 4, vm.stack.data[vm.stack.sp-1])
	result := vm.stack.Pop().(int)

	assert.Equal(t, 4, result)

	// assert.Equal(t, "FOO", string(result))
}
