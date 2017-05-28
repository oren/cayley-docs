I have a few question regarding [this code](https://github.com/oren/test-cayley/blob/f172d111804dcf380b9d42535a7714eaec662e35/insert2/main.go) that inserts admins into my cayley:

1. why transactions? so if something fail it will rollback?

So an "entity" (admin) is either added entirely or not, rather than just like adding is_a "admin" but no hashed_password. It also increases performance on many backends.

2. why IRI for subjects and predicates - `	t.AddQuad(quad.Make(quad.IRI(uuid), quad.IRI("is_a"), quad.String("admin"), nil))
`

* Basically, as @dennwc pointed out -- it is a differentiation of kind. Consider the triple "Denn". That is representing that Entity: Robert -> Entity: Likes -> String: Denn. In that case, I like the NAME "Denn" -- not Denn in specific. The IRI tells you it is an addressable unique entity. .
* It is also in the RDF N-Quads standard.

3. Is the quad.String() needed? why not just string? `t.AddQuad(quad.Make(quad.IRI(uuid), quad.IRI("is_a"), quad.String("admin"), nil))`

4. Should the admin object be an IRI instead of a string? how to decide?

If it is an entity that might have .Outs from it, it should probably be an IRI, if it is something pointed to, like a string, or a value like the int 25 -- it should not be an IRI, because like 25 shouldn't have .Outs because you have no idea what is going to point .In to it, because all values of 25 will point it it.

5. `Out()`, `In()`, and `Has()` - is it true that the first argument is always a predicate?

Yes. Also .Out() and .In() can have multiple predicates in them as multiple predicates to follow. .Has() takes one predicate ban can take multiple nodes, so you can say .Has("name", "Denn", "Robert") which is basically Has the predicate name with the value "Denn" or "Robert".

6. Is it accurate to say that `Regex()` and `Has()` is a way to fliter the numbers of nodes on my path object?

Yes.

7. Can I see the path while I am building the query? (to help me troubleshoot issues)

You can iterate it at any point by basically commenting out the rest, or even returning a path at each step rather than chaining.

```
p := cayley.StartPath(h).Out(quad.IRI("email")).Regex(email).In(quad.IRI("email")).Has(quad.IRI("is_a"), quad.String("admin")).Tag("id").Save(quad.IRI("email"), "email").Save(quad.IRI("hashed_password"),("hashed_password")
```

can be like, assuming you wrote a printPath function

```
p :=  cayley.StartPath(h)
printPath(p)
p = p.Out(quad.IRI("email"))
printPath(p)
p = p.Regex(email)
printPath(p)
p = p.In(quad.IRI("email"))
printPath(p)
p = p.Has(quad.IRI("is_a"), quad.String("admin"))
printPath(p)
p = p.Tag("id").Save(quad.IRI("email"), "email").Save(quad.IRI("hashed_password"),("hashed_password")
printPath(p) // path here will be over same nodes, as last step but with new Tags assos with it.
```
8. Can someone explain [the comment in line 89](https://github.com/oren/test-cayley/blob/f172d111804dcf380b9d42535a7714eaec662e35/insert2/main.go#L89) - Use lock to make sure between check and write we don't have one slip in.

It references the line above it -- 88 -- in this. We need to make sure we don't insert duplicate emails. This email will be related to a unique UUID, so it won't conflict on quad insertion. So to ensure we don't insert duplicate emails, we have to check if email is in system, and if it isn't we need to insert it. Inbetween the check of email, and the insertion of email, we could have someone else also check and insert an email, hence making ours a duplicate! So, there are a lot of ways to deal with this, but the easy one is sort of a global lock around this bit of code, lock code, check email, insert email, unlock code.


Thanks!
