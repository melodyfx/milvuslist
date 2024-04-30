package main

import "fmt"

type Schema struct {
	desc         string
	dynamicField bool
	fields       []*Field
}

func NewSchema() *Schema {
	return &Schema{}
}
func (s *Schema) show() {
	fmt.Println("字段:")
	fmt.Printf("EnableDynamicField:%t\n", s.dynamicField)
	for _, f := range s.fields {
		f.show()
	}
}

func (s *Schema) Add(f *Field) {
	s.fields = append(s.fields, f)
}
