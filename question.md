```
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

```

clinic.json
```
{
  "name": "Heal Now",
  "Address": "3234 Rot Road, Singapore",
  "CreatedBy": "Doe Joe",
  "hours": {
    "mon": [
      {"1": {"from": "08:00", "to": "12:00"}},
      {"2": {"from": "13:00", "to": "15:30"}},
      {"3": {"from": "16:00", "to": "19:00"}}
    ],
    "tue": [
      {"1": {"from": "09:00", "to": "12:30"}},
      {"2": {"from": "13:00", "to": "18:00"}}
    ]
  }
}
```
