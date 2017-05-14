package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
)

var dbPath = "/tmp/db.boltdb"

type Admin struct {
	ID             quad.IRI `quad:"@id"`
	Name           string   `json:"name" quad:"name"`
	Email          string   `json:"email" quad:"email"`
	HashedPassword string   `json:"hashedPassword"  quad:"hashed_password"`
}

func init() {
	schema.RegisterType("Admin", Admin{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func genID() quad.IRI {
	return quad.IRI(fmt.Sprintf("%x", rand.Intn(0xffff)))
}

func main() {
	store := initializeAndOpenGraph(dbPath)

	checkErr(Insert(store, Admin{
		ID:             genID(),
		Name:           "admin1",
		Email:          "foo@gmail.com",
		HashedPassword: "435iue8uou9eu",
	}))

	printAllQuads(store)
	printAllAdmins(store)

}

// helper functions

func Insert(h *cayley.Handle, o interface{}) error {
	qw := graph.NewWriter(h)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	_, err := schema.WriteAsQuads(qw, o)
	return err
}

func printAllQuads(store *cayley.Handle) {
	it := store.QuadsAllIterator()
	defer it.Close()
	fmt.Println("\nquads:")
	for it.Next() {
		fmt.Println(store.Quad(it.Result()))
	}
	fmt.Println()
}

func printAllAdmins(store *cayley.Handle) {
	// get all admins
	var admins []Admin
	checkErr(schema.LoadTo(nil, store, &admins))
	fmt.Println("admins:")
	for _, a := range admins {
		fmt.Printf("%+v\n", a)
	}
	fmt.Println()
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
