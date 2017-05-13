package main

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	uuid "github.com/satori/go.uuid"
)

var dbPath = "/tmp/db.boltdb"

type Admin struct {
	Email          string `json:"email" quad:"email"`
	HashedPassword string `json:"hashedPassword"  quad:"hashed_password"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store := initializeAndOpenGraph(dbPath)

	schema.GenerateID = func(_ interface{}) quad.Value {
		return quad.IRI(uuid.NewV1().String())
	}

	qw := graph.NewWriter(store)

	email := "foo@gmail.com"
	hash := "435iue8uou9eu"

	admin := Admin{
		email,
		hash,
	}

	var id quad.Value

	id, err := schema.WriteAsQuads(qw, admin)
	checkErr(err)
	fmt.Println("admin", admin)
	fmt.Println("id", id)
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
