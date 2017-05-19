// usage: go run main.go -email baz@gmail.com -password 12345 -name "josh"

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/schema"
)

var dbPath = "/tmp/db.boltdb"

func init() {
	schema.RegisterType("Admin", Admin{})
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
	email := flag.String("email", "", "Admin's email")
	password := flag.String("password", "", "Admin's password")
	name := flag.String("name", "", "Admin's name")
	flag.Parse()

	if *email == "" || *password == "" || *name == "" {
		fmt.Println("Arguments must include email,  password, and name")
		os.Exit(0)
	}

	store := initializeAndOpenGraph(dbPath)

	hash, err := hashPassword(*password)
	checkErr(err)

	err = Insert(store, Admin{
		Name:           *name,
		Email:          *email,
		HashedPassword: hash,
		Password:       "foobar",
		LoggedIn:       "hey",
	})

	checkErr(err)

	fmt.Println("Admin was created.")
}

// helper functions

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Insert(h *cayley.Handle, o interface{}) error {
	qw := graph.NewWriter(h)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	_, err := schema.WriteAsQuads(qw, o)
	return err
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
