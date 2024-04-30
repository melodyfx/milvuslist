package main

import (
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Field struct {
	name       string
	ftype      string
	primaryKey bool
	autoId     bool
	typeParam  map[string]string
	indexs     []entity.Index
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) show() {
	fmt.Printf("%s %s", f.name, f.ftype)
	if f.primaryKey {
		fmt.Printf(" primaryKey:%t", f.primaryKey)
	}
	if f.autoId {
		fmt.Printf(" autoid:%t", f.autoId)
	}
	for k, v := range f.typeParam {
		fmt.Printf(" %s:%s ", k, v)
	}
	for _, idx := range f.indexs {
		fmt.Printf(" 索引名称:%s", idx.Name())
		for k, v := range idx.Params() {
			fmt.Printf(" %s:%s ", k, v)
		}
	}
	fmt.Println()

}
