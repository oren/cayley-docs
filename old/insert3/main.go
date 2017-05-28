package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
)

var dbPath = "/tmp/db.boltdb"

type User struct {
	ID             quad.IRI   `quad:"@id"`
	Name           string     `json:"name" quad:"name"`
	Email          string     `json:"email" quad:"email"`
	HashedPassword string     `json:"hashedPassword"  quad:"hashed_password"`
	Follows        []quad.IRI `quad:"follows"`
}

type Admin struct {
	ID             quad.IRI `quad:"@id"`
	Name           string   `json:"name" quad:"name"`
	Email          string   `json:"email" quad:"email"`
	HashedPassword string   `json:"hashedPassword"  quad:"hashed_password"`
}

type Post struct {
	ID      quad.IRI  `quad:"@id"`
	Author  quad.IRI  `quad:"author"`
	Message string    `quad:"msg"`
	Created time.Time `quad:"created"`
}

type UserAndPosts struct {
	ID   quad.IRI `quad:"@id"`
	Name string   `json:"name" quad:"name"`
	// reverse link to posts (In instead of Out)
	// do not load users without posts (require link)
	Posts []Post `quad:"author < *,required"`
}

func init() {
	schema.RegisterType("Admin", Admin{})
	schema.RegisterType("User", User{})
	schema.RegisterType("Post", Post{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Insert(h *cayley.Handle, o interface{}) error {
	qw := graph.NewWriter(h)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	_, err := schema.WriteAsQuads(qw, o)
	return err
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

	u1 := User{
		ID:             genID(),
		Name:           "bob",
		Email:          "bob@gmail.com",
		HashedPassword: "123",
	}

	u2 := User{
		ID:             genID(),
		Name:           "alice",
		Email:          "alice@gmail.com",
		HashedPassword: "abc",
	}
	u2.Follows = append(u2.Follows, u1.ID)

	checkErr(Insert(store, u1))
	checkErr(Insert(store, u2))

	checkErr(Insert(store, Post{
		ID:      genID(),
		Author:  u1.ID,
		Message: "This is an interesting post about graphs",
		Created: time.Now(),
	}))

	checkErr(Insert(store, Post{
		ID:      genID(),
		Author:  u1.ID,
		Message: "This is a second post",
		Created: time.Now(),
	}))

	checkErr(Insert(store, Post{
		ID:      genID(),
		Author:  u2.ID,
		Message: "Hi!",
		Created: time.Now(),
	}))

	//---------------------------------------------------------

	it := store.QuadsAllIterator()
	defer it.Close()
	fmt.Println("\nquads:")
	for it.Next() {
		fmt.Println(store.Quad(it.Result()))
	}
	fmt.Println()

	// get all admins
	var admins []Admin
	checkErr(schema.LoadTo(nil, store, &admins))
	fmt.Println("admins:")
	for _, a := range admins {
		fmt.Printf("%+v\n", a)
	}
	fmt.Println()

	// get all users
	var users []Admin
	checkErr(schema.LoadTo(nil, store, &users))
	fmt.Println("users:")
	for _, a := range users {
		fmt.Printf("%+v\n", a)
	}
	fmt.Println()

	// get all posts
	var posts []Post
	checkErr(schema.LoadTo(nil, store, &posts))
	fmt.Println("posts:")
	for _, p := range posts {
		fmt.Printf("%+v\n", p)
	}
	fmt.Println()

	// get posts grouped by individual users:
	var uposts []UserAndPosts
	checkErr(schema.LoadTo(nil, store, &uposts))
	fmt.Println("posts by users:")
	for _, p := range uposts {
		fmt.Printf("%+v\n", p)
	}
	fmt.Println()

	//---------------------------------------------------------

	// get user by email
	var u User
	p := path.StartPath(store).Has(quad.IRI("email"), quad.String("alice@gmail.com"))
	checkErr(schema.LoadPathTo(nil, store, &u, p))
	fmt.Println("user by email:")
	fmt.Printf("%+v\n\n", u)

	// get posts of people followed by alice (news feed?)
	var news []Post
	p = path.StartPath(store).Has(quad.IRI("name"), quad.String("alice")).Out(quad.IRI("follows")).In(quad.IRI("author"))
	checkErr(schema.LoadPathTo(nil, store, &news, p))
	fmt.Println("news feed of alice:")
	for _, p := range news {
		fmt.Printf("%+v\n", p)
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
