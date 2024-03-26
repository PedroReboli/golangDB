package main

import "godb/godb"

func main() {
	
	// parser.
	t := godb.Tokenize(`
	person{
		name: string(32)
		age: int8
		optional: binary(128)
		alive: bool
		happy: bool
		working: bool
		single: bool
	}`)
	tt := godb.TokenParser(t)
	tt.Init("table.godb")
	tt.Close()

	tt = godb.OpenFile("table.godb")
	tt.InsertRow([]any{"Pedro",int8(21), []byte("a test"), false, true, false, true })
	tt.InsertRow([]any{"Mike",int8(32), []byte("another tetr"), false, false, false, false })

	tt.ReadAllRows()
	tt.Close()
}

/*
types
int8, int16, int32, int64
uint8, uint16, uint32, uint64
string(size)
bool
binary(size)

person{
	name: string(32)
	age: int8
	optional: binary(128)
	alive: bool
	happy: bool
	working: bool
	single: bool
} 

*/
