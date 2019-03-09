package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const cellSize = 8

type stack struct {
	tos  int
	data []byte
}

var s = &stack{
	-cellSize,
	make([]byte, 124*cellSize),
}

var r = &stack{
	-cellSize,
	make([]byte, 124*cellSize),
}

func push(val uint64) {
	s.tos += cellSize
	binary.BigEndian.PutUint64(s.data[s.tos:], val)
}

func pop() {
	s.tos -= cellSize
}

func rpush(val uint64) {
	r.tos += cellSize
	binary.BigEndian.PutUint64(r.data[r.tos:], val)
}

func rpop() {
	r.tos -= cellSize
}

func fetch() {
	ridx := int(binary.BigEndian.Uint64(s.data[s.tos:]))
	_ = copy(s.data[s.tos:s.tos+cellSize], r.data[ridx:ridx+cellSize])
}

func store() {
	ridx := int(binary.BigEndian.Uint64(s.data[s.tos:]))
	_ = copy(r.data[ridx:ridx+cellSize], s.data[s.tos-cellSize:s.tos])
	s.tos -= 2 * cellSize
}

func prn() string {
	return fmt.Sprintf("%s", hex.Dump(s.data[:s.tos+8]))
}
