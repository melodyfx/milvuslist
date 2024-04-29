package main

import (
	"fmt"
)

type Field struct {
	name  string
	ftype string
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) show() {
	fmt.Printf("%s %s\n", f.name, f.ftype)
}
