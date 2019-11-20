package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// regex to match numbers
var numRegex = regexp.MustCompile("^[-+]?[0-9]+.?[[0-9]*]?$")

const cellSize = 8

type num struct {
	ptr bool
	val int
}

type stack struct {
	tos  int
	data []byte
}

type vmfunc func(*VM) error

type VM struct {
	words map[string]vmfunc
	d     stack
	r     stack
	s     stack
	input string
}

func main() {
	v := NewVM()
	rd := bufio.NewReader(os.Stdin)
	lineNum := 1
	for {
		fmt.Print(fmt.Sprintf("[%d]=> ", lineNum))
		str, _ := rd.ReadString('\n')
		if strings.TrimSpace(str) == ":quit" || strings.TrimSpace(str) == ":q" {
			fmt.Println("Goodbye!")
			break
		}
		if strings.TrimSpace(str) == "" {
			continue
		}
		v.input = str + "\n"
		err := parseLine(v)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("data:", v.d.data[:v.d.tos])
			fmt.Println("return:", v.r.data[:v.r.tos])
		}
		lineNum++
	}
}

func NewVM() *VM {
	return &VM{
		makeCoreWords(),
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
		"",
	}
}

func parseLine(v *VM) error {
	items := strings.Fields(v.input)
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if err := parseItem(v, item); err != nil {
			return err
		}
	}
	return nil
}

func parseItem(v *VM, s string) error {
	if strings.HasPrefix(s, "\\") {
		b := []byte(s)
		c := copy(v.r.data[v.r.tos:], b)
		if c != len(b) {
			return fmt.Errorf("quote copy failed: copied %d, should be %d", c, len(b))
		}
		v.r.tos += len(b)
		binary.BigEndian.PutUint64(v.d.data[v.d.tos:], uint64(v.d.tos))
		v.d.tos += cellSize
		binary.BigEndian.PutUint64(v.d.data[v.d.tos:], uint64(len(b)))
		v.d.tos += cellSize
		return nil

	}
	if numRegex.Match([]byte(s)) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint64(v.d.data[v.d.tos:], uint64(i))
		v.d.tos += cellSize
		return nil
	}
	word, ok := v.words[s]
	if !ok {
		return fmt.Errorf("unknown word: %s", s)
	}
	err := word(v)
	if err != nil {
		return err
	}
	return nil
}

func prnDataStack(v *VM) error {
	fmt.Println(fmt.Sprintf("%s", hex.Dump(v.s.data[:v.s.tos])))
	return nil
}

func makeCoreWords() map[string]vmfunc {
	d := make(map[string]vmfunc)
	d["prn"] = prnDataStack
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
