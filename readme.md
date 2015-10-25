# Cayley graph database

```
cayley init --config=cayley.cfg
cayley load --config=cayley.cfg --quads=canada.nq
cayley repl --config=cayley.cfg
cayley http --config=cayley.cfg
```

```
:a "justin trudeau" "in love with" "Sophie Grégoire" .

:a Krissy "lives in" "United States" .
:a Krissy "in love with" "justin trudeau" .
:a Krissy "moves to" "Canada" .

:a Sara "lives in" "Canada" .
:a Sara "in love with" "justin trudeau" .
:a Sara "votes for" "justin trudeau" .

:a Tyler "lives in" "Canada" .
:a Tyler "pissed with" "Sara" .
:a Tyler "moves to" "United States" .

 g.V("Justin Trudeau").Out().All()
g.V("Krissy").Out().All()
g.V("Krissy").Out("in love with").All()
```
