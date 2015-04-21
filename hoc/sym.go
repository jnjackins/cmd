package main

var symbols map[string]Symbol

type Symbol struct {
	val float64               // if VAR
	fn  func(float64) float64 // if BUILTIN
}

func init() {
	symbols = make(map[string]Symbol)
}

func newVar(val float64) Symbol {
	return Symbol{val: val}
}
