package main

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	reset(t)
	push(10)
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(s.data[s.tos:])), "push should put element on tos")
	assert.Equal(t, 0, s.tos, "tos should point at the starting byte of current element")
	push(12)
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(s.data[s.tos:])), "pushed element should replace tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(s.data[s.tos-cellSize:])), "push should keep old data")
	assert.Equal(t, s.tos, cellSize, "tos should advance by 'cellSize' each push at the starting byte of current element")
}

func TestRPush(t *testing.T) {
	reset(t)
	rpush(10)
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(r.data[r.tos:])), "rpush should put element on tos")
	assert.Equal(t, 0, r.tos, "tos should point at the starting byte of current element")
	rpush(12)
	assert.Equal(t, 12, int(binary.BigEndian.Uint64(r.data[r.tos:])), "rpushed element should replace tos")
	assert.Equal(t, 10, int(binary.BigEndian.Uint64(r.data[r.tos-cellSize:])), "rpush should keep old data")
	assert.Equal(t, r.tos, cellSize, "tos should advance by 'cellSize' each push at the starting byte of current element")
}

func TestPop(t *testing.T) {
	reset(t)
	push(10)
	pop()
	assert.Equal(t, -cellSize, s.tos, "tos should start at -cellSize")
	push(1234)
	push(12)
	pop()
	assert.Equal(t, 1234, int(binary.BigEndian.Uint64(s.data[s.tos:])), "pop should remove tos")
	assert.Equal(t, 0, s.tos, "tos should descrement 'cellSize' each pop")
}

func TestRPop(t *testing.T) {
	reset(t)
	rpush(10)
	rpop()
	assert.Equal(t, -cellSize, r.tos, "tos should start at -cellSize")
	rpush(1234)
	rpush(12)
	rpop()
	assert.Equal(t, 1234, int(binary.BigEndian.Uint64(r.data[r.tos:])), "rpop should remove tos")
	assert.Equal(t, 0, r.tos, "tos should descrement 'cellSize' each rpop")
}

func TestStore(t *testing.T) {
	reset(t)
	val, idx := uint64(1122), uint64(3)
	push(val)
	push(idx)
	store()
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(r.data[idx:])), "store should put element at index")
	assert.Equal(t, -cellSize, r.tos, "store should keep tos unchanged")
	assert.Equal(t, -cellSize, s.tos, "store should remove two elements from S")
}

func TestFetch(t *testing.T) {
	reset(t)
	val, idx := uint64(1122), uint64(3)
	push(val)
	push(idx)
	store()
	push(idx)
	fetch()
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(s.data[s.tos:])), "fetch should put the correct data on tos")
	assert.Equal(t, 1122, int(binary.BigEndian.Uint64(r.data[idx:])), "fetch should keep the data in R")
}

func TestPrn(t *testing.T) {
	reset(t)
	push(10)
	str := prn()
	expected := "00000000  00 00 00 00 00 00 00 0a                           |........|\n"
	assert.Equal(t, expected, str, "prn should print S correctly")
}

func reset(t *testing.T) {
	t.Helper()
	s = &stack{
		-cellSize,
		make([]byte, 124*8),
	}
	r = &stack{
		-cellSize,
		make([]byte, 124*8),
	}
}
