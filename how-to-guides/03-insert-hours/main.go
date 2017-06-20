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
	"strings"
)

var dbPath = "db.boltdb"

type Admin struct {
	Name           string `json:"name" quad:"name"`
	Email          string `json:"email" quad:"email"`
	HashedPassword string `json:"hashedPassword"  quad:"hashed_password"`
}

type Clinic struct {
	Name      string         `json:"name" quad:"name"`
	Address1  string         `json:"address" quad:"address"`
	CreatedBy quad.IRI       `quad:"createdBy"`
	Hours     []OpeningHours `quad:"schema:openingHoursSpecification"`
}

type OpeningHours struct {
	DayOfWeek quad.IRI `json:"day" quad:"schema:dayOfWeek"` // set to one of consts like the one above
	Slot      int      `json:"slot" quad:"slot"`
	Opens     string   `json:"opens" quad:"schema:opens"` // ex: 12:00 or 12:00:00
	Closes    string   `json:"closes" quad:"schema:closes"`
}

func init() {
	schema.RegisterType("Admin", Admin{})
	schema.RegisterType("Clinic", Clinic{})
	schema.RegisterType("schema:OpeningHoursSpecification", OpeningHours{})
	schema.GenerateID = func(_ interface{}) quad.Value {
		return quad.IRI(uuid.NewV1().String())
	}
}

func main() {
	store := initializeAndOpenGraph(dbPath)
	a := Admin{
		Name:           "Josh",
		Email:          "josh_f@gmail.com",
		HashedPassword: "435iue8uou9eu",
	}

	err := insert(store, a)
	checkErr(err)

	adminId, err := findAdminID(store, a.Email)
	checkErr(err)

	const (
		Monday = quad.IRI("http://schema.org/Monday")
	)

	mon1 := OpeningHours{
		DayOfWeek: quad.IRI(Monday),
		Slot:      1,
		Opens:     "8:00",
		Closes:    "12:00",
	}

	mon2 := OpeningHours{
		DayOfWeek: quad.IRI(Monday),
		Slot:      2,
		Opens:     "13:30",
		Closes:    "18:00",
	}

	var hours []OpeningHours
	hours = append(hours, mon1)
	hours = append(hours, mon2)

	c := Clinic{
		Name:      "Healthy Life",
		Address1:  "11 boar st, Singapore 11233",
		CreatedBy: adminId,
		Hours:     hours,
	}

	err = insert(store, c)
	checkErr(err)

	printAdmins(store)
	printClinics(store)
	printQuads(store)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

func insert(h *cayley.Handle, o interface{}) error {
	qw := graph.NewWriter(h)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	_, err := schema.WriteAsQuads(qw, o)
	return err
}

func findAdminID(store *cayley.Handle, email string) (quad.IRI, error) {
	p := cayley.StartPath(store).Has(quad.IRI("email"), quad.String(email))
	id, err := p.Iterate(nil).FirstValue(nil)

	if err != nil {
		return "", err
	}

	return id.(quad.IRI), nil
}

func printQuads(store *cayley.Handle) {
	// get all quads
	it := store.QuadsAllIterator()
	defer it.Close()

	fmt.Println("Quads:")
	fmt.Println("-----")

	for it.Next() {
		fmt.Println(store.Quad(it.Result()))
	}

	fmt.Println()
}

func printAdmins(store *cayley.Handle) {
	// get all admins
	var admins []Admin
	checkErr(schema.LoadTo(nil, store, &admins))

	fmt.Println("Admins:")
	fmt.Println("------")

	for _, a := range admins {
		fmt.Println("Name:", a.Name)
		fmt.Println("Email:", a.Email)
		fmt.Println("Hashed Password:", a.HashedPassword)
	}

	fmt.Println()
}

func printClinics(store *cayley.Handle) {
	// get all admins
	var clinics []Clinic
	checkErr(schema.LoadTo(nil, store, &clinics))

	fmt.Println("Clinics:")
	fmt.Println("-------")

	for _, c := range clinics {
		fmt.Println("Name:", c.Name)
		fmt.Println("Email:", c.Address1)

		for _, h := range c.Hours {
			fmt.Println("Day", strings.Split(string(h.DayOfWeek), "/")[3])
			fmt.Println("Slot", h.Slot)
			fmt.Println("Opens", h.Opens)
			fmt.Println("Closes", h.Closes)
		}
	}

	fmt.Println()
}
