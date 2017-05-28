// How to insert administrator user into Cayley

package main

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
)

var dbPath = "/tmp/db.boltdb"

type Admin struct {
	Email          string `json:"email" quad:"email"`
	HashedPassword string `json:"hashedPassword"  quad:"hashed_password"`
}

func init() {
	schema.RegisterType("Admin", Admin{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store := initializeAndOpenGraph(dbPath)

	qw := graph.NewWriter(store)

	admin := Admin{
		"foo@gmail.com",
		"435iue8uou9eu",
	}

	var id quad.Value

	id, err := schema.WriteAsQuads(qw, admin)
	checkErr(err)
	fmt.Println("admin", admin) // {foo@gmail.com 435iue8uou9eu}
	fmt.Println("id", id)       // _:n7425755028193093394

	var results []Admin
	err = schema.LoadTo(nil, store, &results)
	fmt.Println("results", results) // [{foo@gmail.com 435iue8uou9eu}]
}

func initializeAndOpenGraph(dbFile string) *cayley.Handle {
	graph.InitQuadStore("bolt", dbFile, nil)

	// Open and use the database
	store, err := cayley.NewGraph("bolt", dbFile, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return store
}
