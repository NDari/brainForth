package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// regex to match numbers
var number = regexp.MustCompile("^[-+]?[0-9]+.?[[0-9]*]?$")

type str struct {
	ptr int
	len int
}

type stack struct {
	tos  int
	data []byte
}

type vmfunc func(*VM) error

type VM struct {
	d      map[string]vmfunc
	s      stack
	r      stack
	e      stack
	stream *bufio.ReadWriter
}

func main() {
	b := make([]byte, 1000)
	buffer := bytes.NewBuffer(b)
	v := NewVM(bufio.NewReadWriter(buffer, bufio.NewWriter(os.Stdout)))
	lineNum := 0
	for {
		lineNum++
		fmt.Print(fmt.Sprintf("[%d]=> ", lineNum))
		str, _ := rd.ReadString('\n')
		s, err := e.Eval(str)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(s)
		}
	}
}

func NewVM(stream *bufio.ReadWriter) *VM {
	return &VM{
		makeDefaultDict(),
		stack{
			0,
			make([]byte, 124*8),
		},
		stack{
			0,
			make([]byte, 124*8),
		},
		stack{
			0,
			make([]byte, 124*8),
		},
		stream,
	}
}

func parseItem(v *VM) error {
	s, err := v.input.ReadString(' ')
	if err != nil {
		return err
	}
	if strings.HasPrefix(s, "\\") {
		b := []byte(s)
		c := copy(v.r.data[v.r.tos:], b)
		if c != len(b) {
			return fmt.Errorf("quote copy failed: copied %d, should be %d", c, len(b))
		}
		binary.BigEndian.PutUint64(v.s.data[v.s.tos:], uint64(v.s.tos))
		v.s.tos += 8
		binary.BigEndian.PutUint64(v.s.data[v.s.tos:], uint64(len(b)))
		v.s.tos += 8
		return nil

	}
	if number.Match([]byte(s)) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint64(v.s.data[v.s.tos:], uint64(i))
		return nil
	}
	val, ok := v.d[s]
	if !ok {
		return fmt.Errorf("unknown word: %s", s)
	}
	err = val(v)
	if err != nil {
		return err
	}
	return nil
}

func prnDataStack(v *VM) error {
	_, err := v.output.WriteString(fmt.Sprintf("%s", hex.Dump(v.s.data[:v.s.tos])))
	if err != nil {
		return err
	}
	return nil
}

func makeDefaultDict() map[string]vmfunc {
	d := make(map[string]vmfunc)
	d["prn"] = prnDataStack
	d["parse-item"] = parseItem
	return d
}

// func (v *VM) PushCell(val uint64) {
// 	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], val)
// 	v.s.tos += v.CellSize()
// }

// func (v *VM) Push2Cell(belowTop, top uint64) {
// 	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], belowTop)
// 	v.s.tos += v.CellSize()
// 	binary.BigEndian.PutUint64(v.s.data[v.s.tos:], top)
// 	v.s.tos += v.CellSize()
// }

// func (v *VM) PopCell() {
// 	v.s.tos -= v.CellSize()
// }

// func (v *VM) RPushCell() {
// 	v.s.tos -= v.CellSize()
// 	_ = copy(v.r.data[v.r.tos:v.r.tos+v.CellSize()], v.s.data[v.s.tos:v.s.tos+v.CellSize()])
// 	v.r.tos += v.CellSize()
// }

// func (v *VM) RPush2Cell() {
// 	_ = copy(v.r.data[v.r.tos:v.r.tos+2*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos])
// 	v.r.tos += 2 * v.CellSize()
// 	v.s.tos -= 2 * v.CellSize()
// }

// func (v *VM) RPopCell() {
// 	v.r.tos -= v.CellSize()
// }

// func (v *VM) FetchCell() {
// 	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
// 	_ = copy(v.s.data[v.s.tos-v.CellSize():v.s.tos], v.r.data[ridx*v.CellSize():(ridx+1)*v.CellSize()])
// }

// func (v *VM) Fetch2Cell() {
// 	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
// 	v.s.tos -= v.CellSize()
// 	_ = copy(v.s.data[v.s.tos:(v.s.tos+2)*v.CellSize()], v.r.data[ridx*v.CellSize():(ridx+2)*v.CellSize()])
// 	v.s.tos += 2 * v.CellSize()
// }

// func (v *VM) StoreCell() {
// 	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
// 	_ = copy(v.r.data[ridx*v.CellSize():(ridx+1)*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos-v.CellSize()])
// 	v.s.tos -= 2 * v.CellSize()
// }

// func (v *VM) Store2Cell() {
// 	ridx := int(binary.BigEndian.Uint64(v.s.data[v.s.tos-v.CellSize():]))
// 	v.s.tos -= v.CellSize()
// 	_ = copy(v.r.data[ridx*v.CellSize():(ridx+2)*v.CellSize()], v.s.data[v.s.tos-2*v.CellSize():v.s.tos])
// 	v.s.tos -= 2 * v.CellSize()
// }

// func (v *VM) PrnDataStack() string {
// 	return fmt.Sprintf("%s", hex.Dump(v.s.data[:v.s.tos]))
// }
