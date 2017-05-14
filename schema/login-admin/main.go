// usage: go run main.go -email baz@gmail.com -password 12345
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/graph/path"
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
	email := flag.String("email", "", "Admin's email")
	password := flag.String("password", "", "Admin's password")
	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Println("Arguments must include email and password")
		os.Exit(0)
	}

	store := initializeAndOpenGraph(dbPath)
	found, a := findAdmin(store, *email)

	if !found {
		fmt.Println("Admin not found")
		os.Exit(0)
	}

	fmt.Printf("Admin by email: %+v\n\n", a)
	match := checkPasswordHash(*password, a.HashedPassword)
	fmt.Printf("Login?: %+v\n\n", match)

}

// helper functions
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func findAdmin(store *cayley.Handle, email string) (bool, Admin) {
	var a Admin
	p := path.StartPath(store).Has(quad.IRI("email"), quad.String(email))
	err := schema.LoadPathTo(nil, store, &a, p)

	if err != nil {
		return false, a
	}

	return true, a
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
