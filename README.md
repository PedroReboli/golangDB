# goDB
A just for fun and basic implementation of a database in golang with suport to a basic SQL-like query
### why? EGO
brain farted in class and tried to compare using == in sql, friend called me stupid, created my own database that suports compare with ==
## Types
```
int8, int16, int32, int64
uint8, uint16, uint32, uint64
string(size)
bool
binary(size)
```
## Table
the table definition is
```
table name {
    column name: type(type size)
}
person{
    name: string(32)
    age: int8
    optional: binary(128)
    alive: bool
    happy: bool
    working: bool
    single: bool
}
```
## Query
currently, only `select` query with `int` comparison is supported
```sql
select name, happy from person where age == 30 or age >= 40 and 1==1 or age == age
```
## Insert
only by calling a function
```rust
table.InsertRow([]any{"Pedro", int8(21), []byte("a test"), false, true, false, true })
```
might add a SQL to that not sure

## Optmizations
booleans will be compressed into a single byte<br>
in select query, only needed columns will be deserialized<br>
~~that's no allot~~

## Limitations
only one table per file is supported<br>
no insert query<br>
select query must have where, if `where` not needed, it is possible to use `1==1`<br>

## Ideas that I might add

### Table validation
```c#
person{
	name: string(32)
	age: int8
	optional: binary(128)
	alive: bool
	happy: bool
	working: bool
	single: bool
} > {
	if (age < 5 and working){
		panic("kids should not work")
	}
	if (not alive){
		working = false
		happy = false // he shouldn't be working
	}
	if (name == ""){
		panic("name should not be empty")
	}
}
```

### Select with code
```c#
select (name, age, alive, working) > {
	if (alive){
		if (working){
			age += 2
		}else{
			age += 1
		}
		if (name == "kid"){
			age = 4
		}
	}
	return (name, age)
} from person where alive = true // is this too much vodoo?
```

## Or I'll just forget about this project, not sure





