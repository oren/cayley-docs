// Example for adding struct and removing it using it's id

// usage:
// go run main.go
// go run main.go -id 9c2177ed-408c-11e7-af40-843a4b0f5a10

// You should see nothing after the second command and the db should be empty

package main

import (
	"flag"
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

func init() {
	schema.RegisterType("Admin", Admin{})
	schema.GenerateID = func(_ interface{}) quad.Value {
		return quad.IRI(uuid.NewV1().String())
	}
}

type Admin struct {
	Name           string `json:"name" quad:"name"`
	Email          string `json:"email" quad:"email"`
	HashedPassword string `json:"hashedPassword"  quad:"hashed_password"`
	Password       string `quad:"-"`
	LoggedIn       string `quad:"opt"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	adminId := flag.String("id", "", "Admin's Id")
	flag.Parse()

	store := initializeAndOpenGraph(dbPath)

	if *adminId != "" {
		id := quad.IRI(*adminId)

		err := store.RemoveNode(store.ValueOf(id))

		if err != nil {
			fmt.Println("Error removing the node", err)
		}

		printAllQuads(store)
		return
	}

	err := Insert(store, Admin{
		Name:           "josh",
		Email:          "josh@test.com",
		HashedPassword: "abc",
	})

	checkErr(err)

	fmt.Println("Admin was created.")
	printAllQuads(store)
}

// helper functions

func initializeAndOpenGraph(dbFile string) *cayley.Handle {
	graph.InitQuadStore("bolt", dbFile, nil)

	// Open and use the database
	store, err := cayley.NewGraph("bolt", dbFile, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return store
}

func Insert(store *cayley.Handle, o interface{}) error {
	qw := graph.NewWriter(store)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	_, err := schema.WriteAsQuads(qw, o)
	return err
}

func printAllQuads(store *cayley.Handle) {
	quads, err := readAllQuads(store)

	if err != nil {
		fmt.Println("error reading quads:", err)
	}

	for _, q := range quads {
		fmt.Println(q)
	}
}

func readAllQuads(store *cayley.Handle) ([]quad.Quad, error) {
	var results []quad.Quad
	it := store.QuadsAllIterator()
	defer it.Close()

	for it.Next() {
		results = append(results, store.Quad(it.Result()))
	}

	return results, nil
}
