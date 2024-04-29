package main

import "fmt"

type Collection struct {
	id                int
	name              string
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
	fmt.Printf("\n%d.collection:%s\n", c.id, c.name)
	fmt.Printf("consistency_level:%s\n", c.name)
	fmt.Printf("nums: %s\n", c.nums)
	fmt.Printf("shardNum:%d\n", c.shardNum)
	c.schema.show()
}

func (s *Collection) Add(schema *Schema) {
	s.schema = schema
}
