package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"context"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	uuid "github.com/satori/go.uuid"
)

var dbPath = "db.boltdb"

type Admin struct {
	Name           string `json:"name" quad:"name"`
	Email          string `json:"email" quad:"email"`
	HashedPassword string `json:"hashedPassword"  quad:"hashed_password"`
}

type Clinic struct {
	ID        quad.IRI       `json:"id" quad:"@id"`
	Name      string         `json:"name" quad:"name"`
	Address1  string         `json:"address" quad:"address"`
	CreatedBy quad.IRI       `quad:"createdBy"`
	OfficeTel string         `json:"officeTel" quad:"officeTel"`
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

	id, err := insert(store, a)
	checkErr(err)

	adminId, err := findAdminID(store, a.Email)
	checkErr(err)

	const (
		Monday = quad.IRI("http://schema.org/Monday")
	)

	c := loadJSON("clinic.json")
	c.CreatedBy = adminId

	// var id quad.Value
	id, err = insert(store, c)
	checkErr(err)

	c = loadJSON("updated-clinic.json")

	err = update(store, c, id)
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

// LoadJson loads json of a clinic and returns a Clinic struct
func loadJSON(JSONFile string) *Clinic {
	raw, err := ioutil.ReadFile(JSONFile)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Clinic

	err = json.Unmarshal(raw, &c)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return &c
}

func insert(h *cayley.Handle, o interface{}) (quad.Value, error) {
	qw := graph.NewWriter(h)
	defer qw.Close() // don't forget to close a writer; it has some internal buffering
	id, err := schema.WriteAsQuads(qw, o)
	return id, err
}

func update(h *cayley.Handle, o interface{}, id quad.Value) error {
	// find a clinic with a given name
	// read fields from json file and update clinic
	// save clinic

	var clinic Clinic
	t := cayley.NewTransaction()
	checkErr(schema.LoadTo(nil, h, &clinic, id))
	clinic.Address1 = "Seasame st"    // TODO: read this from the json file
	clinic.OfficeTel = "999 999 9999" // TODO: read this from the json file
	id, err := insert(h, clinic)
	checkErr(err)
	err = h.ApplyTransaction(t)
	checkErr(err)

	return nil
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

	ctx := context.TODO()

	for it.Next(ctx) {
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
		fmt.Println("Address:", c.Address1)
		fmt.Println("OfficeTel:", c.OfficeTel)

		for _, h := range c.Hours {
			fmt.Println("---")
			fmt.Println("h", h)
			fmt.Println("---")
			// fmt.Println("Day", strings.Split(string(h.DayOfWeek), "/")[3])
			fmt.Println("Day", string(h.DayOfWeek))
			fmt.Println("Slot", h.Slot)
			fmt.Println("Opens", h.Opens)
			fmt.Println("Closes", h.Closes)
		}
		fmt.Println("----------------------------")
	}

	fmt.Println()
}
