package main

import (
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Collection struct {
	id                int
	name              string
	loadStat          entity.LoadState
	consistency_level string
	nums              string
	shardNum          int32
	schema            *Schema
}

func NewCollection(name string) *Collection {
	return &Collection{
		name: name,
	}
}
func (c *Collection) show() {
	var status string
	switch c.loadStat {
	case entity.LoadStateNotExist:
		status = "NotExist"
	case entity.LoadStateNotLoad:
		status = "NotLoad"
	case entity.LoadStateLoading:
		status = "Loading"
	case entity.LoadStateLoaded:
		status = "Loaded"
	}
	fmt.Printf("\n%d.collection:%s\n", c.id, c.name)
	fmt.Printf("loadStat:%s\n", status)
	fmt.Printf("consistency_level:%s\n", c.consistency_level)
	fmt.Printf("nums: %s\n", c.nums)
	fmt.Printf("shardNum:%d\n", c.shardNum)
	c.schema.show()
}

func (s *Collection) Add(schema *Schema) {
	s.schema = schema
}
