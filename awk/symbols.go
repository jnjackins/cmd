package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

// tvals
const (
	numVal = 1 << iota
	strVal
	varVal
)

var symbols map[string]*symbol

type symbol struct {
	name string  // name, for variables only
	sval string  // string value
	fval float64 // value as number
	tval int     // type info: numVal|strVal|varVal
}

func (s *symbol) String() string {
	name := s.name
	if s.name == "" {
		name = "<unnamed>"
	}
	if s.tval&numVal > 0 {
		return fmt.Sprintf("(%s | %p: %.2f)", name, s, s.fval)
	}
	return fmt.Sprintf("(%s | %p: %#v)", name, s, s.sval)
}

// TODO: bad name (same as setString and setNum but different behaviour)
func setVar(name, sval string) {
	sym := &symbol{name: name}
	if s, ok := symbols[sval]; ok {
		sym.sval = s.sval
		sym.fval = s.fval
		sym.tval = s.tval
	} else {
		sym.sval = sval
		sym.tval = strVal
	}
	symbols[name] = sym
}

func (s *symbol) setString(v string) *symbol {
	s.sval = v
	s.tval &= ^numVal // invalidate s as a number
	s.tval |= strVal  // s is now a valid string
	return s
}

func (s *symbol) setNum(v float64) *symbol {
	s.fval = v
	s.tval &= ^strVal // invalidate s as a string
	s.tval |= numVal  // s is now a valid number
	return s
}

func (s *symbol) getString() string {
	if (s.tval & (numVal | strVal)) == 0 {
		log.Printf("error: getString: symbol is neither numVal or strVal: %v", s)
		return ""
	}
	if s.tval&strVal > 0 {
		return s.sval
	}
	if _, frac := math.Modf(s.fval); frac == 0.0 {
		// it's integral
		s.sval = fmt.Sprintf("%.30g", s.fval)
	} else {
		s.sval = fmt.Sprintf("%.6g", s.fval)
	}
	s.tval |= strVal
	return s.sval
}

func (s *symbol) getNum() float64 {
	if (s.tval & (numVal | strVal)) == 0 {
		log.Print("error: getNum: symbol is neither numVal or strVal")
		return 0.0
	}
	if s.tval&numVal > 0 {
		return s.fval
	}
	f, err := atof(s.sval)
	if err == nil {
		s.fval = f
		s.tval |= numVal
	} else {
		log.Print(err)
	}
	return s.fval
}

func atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func initSymbols() {
	symbols = make(map[string]*symbol)
	symbols["NR"] = &symbol{name: "NR", tval: numVal}
	symbols["NF"] = &symbol{name: "NF", tval: numVal}
}

func updateSymbols(rec record) {
	s := symbols["NR"]
	s.setNum(s.getNum() + 1.0)
	s = symbols["NF"]
	s.setNum(float64(len(rec.fields)))
	for i := 1; i <= len(rec.fields); i++ {
		nm := fmt.Sprintf("$%d", i)
		s := getSymbol(nm)
		s.setString(rec.fields[i-1])
	}
	// TODO: clear the rest
	dprintf("record: NR: %.0f, NF: %.0f, fields: %v", symbols["NR"].fval, symbols["NF"].fval, rec.fields)
}

func getSymbol(name string) *symbol {
	if s, ok := symbols[name]; ok {
		dprintf("getSymbol: loaded existing symbol: %v", s)
		return s
	}
	s := &symbol{name: name, tval: strVal | numVal} // zero-value symbols print as an empty string or zero
	symbols[name] = s
	dprintf("getSymbol: created new symbol: %v", s)
	return s
}
