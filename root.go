package main

import "fmt"

type Root struct {
	dbs []*Database
}

func (r *Root) show() {
	for _, d := range r.dbs {
		fmt.Println()
		d.show()
	}
}

func (r *Root) Add(d *Database) {
	r.dbs = append(r.dbs, d)
}
