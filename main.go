package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// regex to match numbers
var numRegex = regexp.MustCompile("^[-+]?[0-9]+.?[[0-9]*]?$")

type vmfunc func(*VM) error

type VM struct {
	words  map[string]vmfunc
	defs   map[string]string
	macros map[string]string
	data   []string
}

func main() {
	v := NewVM()
	rd := bufio.NewReader(os.Stdin)
	lineNum := 1
	verbose := false
	for {
		fmt.Print(fmt.Sprintf("[%d]=> ", lineNum))
		str, _ := rd.ReadString(';')
		str = strings.TrimSpace(str)
		str = strings.TrimSuffix(str, ";")
		if str == ":q" {
			fmt.Println("Goodbye!")
			break
		}
		if str == "" {
			continue
		}
		if str == ":v" {
			verbose = true
			continue
		}
		err := v.interpret(str)
		if err != nil {
			fmt.Println(err)
		}
		if verbose {
			fmt.Println("data:", v.data)
		}
		lineNum++
	}
}

func NewVM() *VM {
	return &VM{
		words:  makeCoreWords(),
		defs:   make(map[string]string),
		macros: make(map[string]string),
		data:   make([]string, 0),
	}
}

func (v *VM) interpret(s string) error {
	items := strings.Fields(s)
	if len(items) == 0 {
		return nil
	}
	if items[0] == "mac" {
		body := ""
		name := items[1]
		for _, item := range items[2:] {
			body = fmt.Sprintf("%s %s", body, item)
		}
		v.macros[name] = body
		return nil
	}
	for _, item := range items {
		if strings.HasPrefix(item, "'") {
			v.data = append(v.data, item)
			continue
		}
		w, ok := v.words[item]
		if ok {
			w(v)
			continue
		}
		body, ok := v.defs[item]
		if ok {
			err := v.interpret(body)
			if err != nil {
				return fmt.Errorf("error in calling user definition %s: %v", item, err)
			}
			continue
		}
		v.data = append(v.data, item)
	}
	return nil
}

func prnDataStack(v *VM) error {
	fmt.Println(v.data)
	return nil
}

func first(v *VM) error {
	fmt.Println(v.data[len(v.data)-1])
	return nil
}

func second(v *VM) error {
	fmt.Println(v.data[len(v.data)-2])
	return nil
}

func third(v *VM) error {
	fmt.Println(v.data[len(v.data)-3])
	return nil
}

func prnDefs(v *VM) error {
	fmt.Println(v.defs)
	return nil
}

func words(v *VM) error {
	for k := range v.words {
		fmt.Println(k)
	}
	return nil
}

func makeCoreWords() map[string]vmfunc {
	d := make(map[string]vmfunc)
	d["$0"] = prnDataStack
	d["$1"] = first
	d["$2"] = second
	d["$3"] = third
	d["defs"] = prnDefs
	d["words"] = words
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
