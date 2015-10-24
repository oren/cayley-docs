# Cayley graph database

```
cayley init --config cayley.cfg
cayley repl --config=cayley.cfg
cayley http --config=cayley.cfg
```

```
:a Krissy "lives in" "United States" .
:a Krissy "in love with" "justin trudeau" .
:a Krissy "moves to" "Canada" .
g.V("Krissy").Out().All()
```


