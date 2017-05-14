package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	dbPath       = "/tmp/db.boltdb"
	ErrBadFormat = errors.New("invalid email format")
	emailRegexp  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Admin struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"-"` // - so it doesn't get encoded to json ever
	HashedPassword string `json:"hashedPassword"`
}

func initializeAndOpenGraph(dbFile string) *cayley.Handle {
	graph.InitQuadStore("bolt", dbFile, nil)

	store, err := cayley.NewGraph("bolt", dbFile, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return store
}

func main() {
	h := initializeAndOpenGraph(dbPath)
	a1 := Admin{
		Email:    "me@fake.com",
		Password: "tobehashed",
	}
	a2 := Admin{
		Email:    "me2@fake.com",
		Password: "tobehashed2",
	}
	a3 := Admin{
		Email:    "dog@fake.com",
		Password: "tobehashed2",
	}

	err := CreateAdmin(h, a1)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateAdmin(h, a2)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateAdmin(h, a3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n==== All admins ====")
	As, err := ReadAdmins(h, regexp.MustCompile(".*"))
	if err != nil {
		log.Fatal(err)
	}
	PrintAdmins(As)

	fmt.Println("\n==== Below is only emails that start with me* ====")
	As, err = ReadAdmins(h, regexp.MustCompile("^me.*"))
	if err != nil {
		log.Fatal(err)
	}
	PrintAdmins(As)

}

// TODO: Check for duplicate email
// TODO: Use lock to make sure between check and write we don't have one slip in
func CreateAdmin(h *cayley.Handle, a Admin) error {
	err := validateEmail(a.Email)
	if err != nil {
		return err
	}

	if a.Password != "" && a.HashedPassword == "" { // if we have pw and no hash, hash it
		a.HashedPassword, err = hashPassword(a.Password)
		if err != nil {
			return err
		}
	}

	uuid := uuid.NewV1().String()

	// if one command fail, rollback
	t := cayley.NewTransaction()
	// both subject and predicate should be IRI. why?

	t.AddQuad(quad.Make("3232ueououueu", "is_a", "admin", nil))

	t.AddQuad(quad.Make(quad.IRI(uuid), quad.IRI("is_a"), quad.String("admin"), nil))
	t.AddQuad(quad.Make(quad.IRI(uuid), quad.IRI("email"), quad.String(a.Email), nil))
	t.AddQuad(quad.Make(quad.IRI(uuid), quad.IRI("hashed_password"), quad.String(a.HashedPassword), nil))
	err = h.ApplyTransaction(t)
	if err != nil {
		return err
	}

	return nil
}

func ReadAdmins(h *cayley.Handle, email *regexp.Regexp) ([]Admin, error) {
	p := cayley.StartPath(h).
		Out(quad.IRI("email")).Regex(email).In(quad.IRI("email")).Has(quad.IRI("is_a"), quad.String("admin")).
		Tag("id").
		Save(quad.IRI("email"), "email").
		Save(quad.IRI("hashed_password"), "hashed_password")

	fmt.Println("p", p)

	results := []Admin{}
	err := p.Iterate(nil).TagValues(nil, func(tags map[string]quad.Value) {
		// fmt.Println("tags", tags)

		// tags["id"] contain a subject node. it's an interface so we have to convert it first to IRI and than to String
		results = append(results, Admin{
			ID:             quad.NativeOf(tags["id"]).(quad.IRI).String(),
			Email:          quad.NativeOf(tags["email"]).(string),
			HashedPassword: quad.NativeOf(tags["hashed_password"]).(string),
		})
	})

	if err != nil {
		return []Admin{}, err
	}

	// return []Admin{}, nil
	return results, nil
}

func PrintAdmins(as []Admin) {
	for _, a := range as {
		fmt.Println("ID: ", a.ID)
		fmt.Println("\tEmail: ", a.Email)
		fmt.Println("\tHashedPassword: ", a.HashedPassword)
	}
}

func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrBadFormat
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
