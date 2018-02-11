# How-to guide

## How to update a clinic

In this scenario we want to update the address and the phone of a clinic. Later we'll do it using a JSON file.


Run the following:
```
go get
go run main.go
```

Here are the interesting lines:
```
t := cayley.NewTransaction()

t.RemoveQuad(quad.Make(id, quad.IRI("address"), "3234 Rot Road, Singapore", nil))
t.AddQuad(quad.Make(id, quad.IRI("address"), "3235 Rot Road, Singapore", nil))

t.RemoveQuad(quad.Make(id, quad.IRI("officeTel"), "65 6100 0939", nil))
t.AddQuad(quad.Make(id, quad.IRI("officeTel"), "75 6100 0939", nil))

err := h.ApplyTransaction(t)
```
