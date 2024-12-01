package core

import (
	"fmt"
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

	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	valueBytes, err := contractState.Get([]byte("FOO"))
	value := deserializeInt64(valueBytes)
	assert.Nil(t, err)
	assert.Equal(t, value, int64(5))
	fmt.Printf("%v", vm.contractState)

}

func TestVMDiv(t *testing.T) {

	contract := []byte{0x04, byte(InstrPushInt), 0x02, byte(InstrPushInt), byte(InstrDiv)}

	contractState := NewState()
	vm := NewVM(contract, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)
	assert.Equal(t, 2, result)
}

func TestVMMul(t *testing.T) {

	contract := []byte{0x04, byte(InstrPushInt), 0x02, byte(InstrPushInt), byte(InstrMul)}

	contractState := NewState()
	vm := NewVM(contract, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)
	assert.Equal(t, 8, result)
}
