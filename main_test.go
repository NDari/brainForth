package main

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	v := NewVM()
	v.push(10)
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "push should put element on tos")
	assert.Equal(t, cellSize, v.s.tos, "tos should point at the starting byte of current element")
	v.push(12)
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "pushed element should replace tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-2*cellSize:])), "push should keep old data")
	assert.Equal(t, 2*cellSize, v.s.tos, "tos should advance by 'cellSize' each push at the starting byte of current element")
}

func TestPush2(t *testing.T) {
	v := NewVM()
	v.push2(10, 12)
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "push2 should put 2nd element on tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-2*cellSize:])), "push2 should put 1st element below tos")
	assert.Equal(t, 2*cellSize, v.s.tos, "push2 should increment tos by 2 cellSize")
}

func TestRPush(t *testing.T) {
	v := NewVM()
	v.push(10)
	v.rpush()
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-cellSize:])), "rpush should put element on tos")
	assert.Equal(t, cellSize, v.r.tos, "tos should point at the starting byte of current element")
	assert.Equal(t, 0, v.s.tos, "element should be removed from S after rpush")
	v.push(12)
	v.rpush()
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-cellSize:])), "rpushed element should replace tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-2*cellSize:])), "rpush should keep old data")
	assert.Equal(t, 2*cellSize, v.r.tos, "tos should advance by 'cellSize' each push at the starting byte of current element")
}
func TestRPush2(t *testing.T) {
	v := NewVM()
	v.push2(10, 12)
	v.rpush2()
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-cellSize:])), "rpush2 should put 2nd element on tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-2*cellSize:])), "rpush2 should put 1st element below tos")
	assert.Equal(t, 2*cellSize, v.r.tos, "rpush2 should increment tos twice")
	assert.Equal(t, 0, v.s.tos, "element should be removed from S after rpush2")
}

func TestPop(t *testing.T) {
	v := NewVM()
	v.push(10)
	v.pop()
	assert.Equal(t, 0, v.s.tos, "tos should start at -cellSize")
	v.push(1234)
	v.push(12)
	v.pop()
	assert.Equal(t, 1234, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "pop should remove tos")
	assert.Equal(t, cellSize, v.s.tos, "tos should descrement 'cellSize' each pop")
}

func TestRPop(t *testing.T) {
	v := NewVM()
	v.push(10)
	v.rpush()
	v.rpop()
	assert.Equal(t, 0, v.r.tos, "tos of empty stack should be 0")
	v.push(1234)
	v.rpush()
	v.push(12)
	v.rpush()
	v.rpop()
	assert.Equal(t, 1234, int(binary.BigEndian.Uint64(v.r.data[v.r.tos-cellSize:])), "rpop should remove tos")
	assert.Equal(t, cellSize, v.r.tos, "tos should descrement 'cellSize' each rpop")
}

func TestStore(t *testing.T) {
	v := NewVM()
	val, idx := uint64(1122), uint64(3)
	v.push(val)
	v.push(idx)
	v.store()
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(v.r.data[idx*cellSize:])), "store should put element at index")
	assert.Equal(t, 0, v.r.tos, "store should keep tos unchanged")
	assert.Equal(t, 0, v.s.tos, "store should remove two elements from S")
}

func TestFetch(t *testing.T) {
	v := NewVM()
	val, idx := uint64(1122), uint64(3)
	v.push(val)
	v.push(idx)
	v.store()
	v.push(idx)
	v.fetch()
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "fetch should put the correct data on tos")
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(v.r.data[idx*cellSize:])), "fetch should keep the data in R")
}

func TestStore2(t *testing.T) {
	v := NewVM()
	val1, val2, idx := uint64(112), uint64(311), uint64(5)
	v.push(val1)
	v.push(val2)
	v.push(idx)
	v.store2()
	assert.Equal(t, 0, v.s.tos, "tos should be at the start should be empty")
	assert.Equal(t, 112, int(binary.BigEndian.Uint64(v.r.data[int(idx*cellSize):int((idx+1)*cellSize)])), "val1 should be at start of address")
	assert.Equal(t, 311, int(binary.BigEndian.Uint64(v.r.data[int((idx+1)*cellSize):int((idx+2)*cellSize)])), "val2 should be one cellsize after address")
}

func TestFetch2(t *testing.T) {
	v := NewVM()
	val1, val2, idx := uint64(112), uint64(311), uint64(5)
	v.push(val1)
	v.push(val2)
	v.push(idx)
	v.store2()
	assert.Equal(t, 0, v.s.tos, "tos should be at the start should be empty")
	v.push(idx)
	v.fetch2()
	assert.Equal(t, 311, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:])), "val1 should be below tos")
	assert.Equal(t, 112, int(binary.BigEndian.Uint64(v.s.data[v.s.tos-2*cellSize:])), "val2 should be on tos")
}

func TestPrn(t *testing.T) {
	v := NewVM()
	v.push(10)
	str := v.prn()
	expected := "00000000  00 00 00 00 00 00 00 0a                           |........|\n"
	assert.Equal(t, expected, str, "prn should print S correctly")
}
