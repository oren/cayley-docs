# How-to guide

## How to update a clinic using a JSON file

In this scenario we want to update the address and the phone of a clinic using a json file.
Our program will look at the name of the clinic, and if it finds it in our DB, it will update the DB with any field we have in the json file.

Run the following:
```
go get
go run main.go
```

Notice the output in the terminal:
```
Admins:
------
Name: Josh
Email: josh_f@gmail.com
Hashed Password: 435iue8uou9eu

Clinics:
-------
Name: Healthy Life
Email: 11 boar st, Singapore 11233

Quads:
-----
<831c71de-43eb-11e7-9cd0-843a4b0f5a10> -- <rdf:type> -> <Clinic>
<831bc569-43eb-11e7-9cd0-843a4b0f5a10> -- <rdf:type> -> <Admin>
<831bc569-43eb-11e7-9cd0-843a4b0f5a10> -- <email> -> "josh_f@gmail.com"
<831c71de-43eb-11e7-9cd0-843a4b0f5a10> -- <address> -> "11 boar st, Singapore 11233"
<831c71de-43eb-11e7-9cd0-843a4b0f5a10> -- <name> -> "Healthy Life"
<831bc569-43eb-11e7-9cd0-843a4b0f5a10> -- <name> -> "Josh"
<831bc569-43eb-11e7-9cd0-843a4b0f5a10> -- <hashed_password> -> "435iue8uou9eu"
<831c71de-43eb-11e7-9cd0-843a4b0f5a10> -- <createdBy> -> <831bc569-43eb-11e7-9cd0-843a4b0f5a10>
```

If you see something similar to the above output, you are doing fine!
You just created an administrator and a clinic and connected between them.

You should try [the second how-to guide](../02-visualize/README.md), and learn how to visualize your data.
