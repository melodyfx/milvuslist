package main

import "fmt"

type Database struct {
	name  string
	colls []*Collection
}

func NewDatabase(name string) *Database {
	return &Database{
		name: name,
	}
}

func (d *Database) show() {
	fmt.Printf("===database:%s===\n", d.name)
	for _, c := range d.colls {
		c.show()
	}
}

func (d *Database) Add(c *Collection) {
	d.colls = append(d.colls, c)
}
