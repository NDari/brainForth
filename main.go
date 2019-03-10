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

type VM struct {
	s *stack
	r *stack
}

func NewVM() *VM {
	return &VM{
		&stack{
			0,
			make([]byte, 124*cellSize),
		},
		&stack{
			0,
			make([]byte, 124*cellSize),
		},
	}
}

func (v *VM) push(val uint64) {
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], val)
	v.s.tos += cellSize
}

func (v *VM) push2(belowTop, top uint64) {
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], belowTop)
	v.s.tos += cellSize
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], top)
	v.s.tos += cellSize
}

func (v *VM) pop() {
	v.s.tos -= cellSize
}

func (v *VM) rpush() {
	v.s.tos -= cellSize
	_ = copy(v.r.data[v.r.tos:v.r.tos+cellSize], v.s.data[v.s.tos:v.s.tos+cellSize])
	v.r.tos += cellSize
}

func (v *VM) rpush2() {
	_ = copy(v.r.data[v.r.tos:v.r.tos+2*cellSize], v.s.data[v.s.tos-2*cellSize:v.s.tos])
	v.r.tos += 2 * cellSize
	v.s.tos -= 2 * cellSize
}

func (v *VM) rpop() {
	v.r.tos -= cellSize
}

func (v *VM) fetch() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:]))
	_ = copy(v.s.data[v.s.tos-cellSize:v.s.tos], v.r.data[ridx*cellSize:(ridx+1)*cellSize])
}

func (v *VM) fetch2() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:]))
	v.s.tos -= cellSize
	_ = copy(v.s.data[v.s.tos:(v.s.tos+2)*cellSize], v.r.data[ridx*cellSize:(ridx+2)*cellSize])
	v.s.tos += 2 * cellSize
}

func (v *VM) store() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:]))
	_ = copy(v.r.data[ridx*cellSize:(ridx+1)*cellSize], v.s.data[v.s.tos-2*cellSize:v.s.tos-cellSize])
	v.s.tos -= 2 * cellSize
}

func (v *VM) store2() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-cellSize:]))
	v.s.tos -= cellSize
	_ = copy(v.r.data[ridx*cellSize:(ridx+2)*cellSize], v.s.data[v.s.tos-2*cellSize:v.s.tos])
	v.s.tos -= 2 * cellSize
}

func (v *VM) prn() string {
	return fmt.Sprintf("%s", hex.Dump(v.s.data[:v.s.tos]))
}
