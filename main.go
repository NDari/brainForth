package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type stack struct {
	tos  int
	data []byte
}

type VM struct {
	cellSize int
	s        *stack
	r        *stack
}

func NewVM(bytesPerCell int) *VM {
	return &VM{
		bytesPerCell,
		&stack{
			0,
			make([]byte, 124*bytesPerCell),
		},
		&stack{
			0,
			make([]byte, 124*bytesPerCell),
		},
	}
}

func (v *VM) CellSize() int {
	return v.cellSize
}

func (v *VM) PushCell(val uint64) {
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], val)
	v.s.tos += v.CellSize()
}

func (v *VM) Push2Cell(belowTop, top uint64) {
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], belowTop)
	v.s.tos += v.CellSize()
	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], top)
	v.s.tos += v.CellSize()
}

func (v *VM) PopCell() {
	v.s.tos -= v.CellSize()
}

func (v *VM) RPushCell() {
	v.s.tos -= v.CellSize()
	_ = copy(v.r.data[v.r.tos:v.r.tos+v.CellSize()], v.s.data[v.s.tos:v.s.tos+v.CellSize()])
	v.r.tos += v.CellSize()
}

func (v *VM) RPush2Cell() {
	_ = copy(v.r.data[v.r.tos:v.r.tos+2*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos])
	v.r.tos += 2 * v.CellSize()
	v.s.tos -= 2 * v.CellSize()
}

// func (v *VM) RStoreBytes(buf []byte) {
// 	l := len(buf)
// 	ridx := uint64(v.r.tos)
// 	length := uint64(l)
// 	_ = copy(v.r.data[v.r.tos:v.r.tos+l], buf)
// 	v.r.tos += l
// 	v.Push2Cell(ridx, length)
// }

// func (v *VM) FetchBytes

func (v *VM) RPopCell() {
	v.r.tos -= v.CellSize()
}

func (v *VM) FetchCell() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
	_ = copy(v.s.data[v.s.tos-v.CellSize():v.s.tos], v.r.data[ridx*v.CellSize():(ridx+1)*v.CellSize()])
}

func (v *VM) Fetch2Cell() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
	v.s.tos -= v.CellSize()
	_ = copy(v.s.data[v.s.tos:(v.s.tos+2)*v.CellSize()], v.r.data[ridx*v.CellSize():(ridx+2)*v.CellSize()])
	v.s.tos += 2 * v.CellSize()
}

func (v *VM) StoreCell() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
	_ = copy(v.r.data[ridx*v.CellSize():(ridx+1)*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos-v.CellSize()])
	v.s.tos -= 2 * v.CellSize()
}

func (v *VM) Store2Cell() {
	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
	v.s.tos -= v.CellSize()
	_ = copy(v.r.data[ridx*v.CellSize():(ridx+2)*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos])
	v.s.tos -= 2 * v.CellSize()
}

func (v *VM) PrnDataStack() string {
	return fmt.Sprintf("%s", hex.Dump(v.s.data[:v.s.tos]))
}
