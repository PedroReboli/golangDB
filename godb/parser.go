package godb

import (
	"fmt"
	"strconv"
	"strings"
)

type BasicTypes int
const(
	uint8Type BasicTypes = iota
	uint16Type 
	uint32Type 
	uint64Type 
	int8Type 
	int16Type 
	int32Type 
	int64Type 
	stringType 
	binaryType 
	boolType 
)

type tokenEnum int
const(
	hummm tokenEnum = iota
	openBraces 
	closeBraces
	openParentheses
	closeParentheses
	newLine
	text
	number
	colon
	EOF
	None
)

var tnum = []string{"hummm",
	"openBraces",
	"closeBraces",
	"openParentheses",
	"closeParentheses",
	"newLine",
	"text",
	"number",
	"colon",
	"EOF",
	"None"}


type Token struct{
	Token tokenEnum
	Option string
}

// type ColunType struct{
// 	Name string
// 	Type BasicTypes
// 	TypeOptions int
// }

// type Table struct{
// 	TableName string
// 	Coluns []ColunType
// }

func getBasicType(basicType string) BasicTypes{
	basicType = strings.TrimSpace(basicType)
	switch basicType{
	case "uint8":
		return uint8Type
	case "uint16":
		return uint16Type
	case "uint32":
		return uint32Type
	case "uint64":
		return uint64Type
	case "int8":
		return int8Type
	case "int16":
		return int16Type
	case "int32":
		return int32Type
	case "int64":
		return int64Type
	case "string":
		return stringType
	case "binary":
		return binaryType
	case "bool":
		return boolType
	default:
		panic(fmt.Sprintf("type %s does not exist", basicType))
	}
}

// func getTypeBasicType(basicType string) BasicTypes{
// 	basicType = strings.TrimSpace(basicType)
// 	switch basicType{
// 	case "uint8":
// 		return uint8Type
// 	case "uint16":
// 		return uint16Type
// 	case "uint32":
// 		return uint32Type
// 	case "uint64":
// 		return uint64Type
// 	case "int8":
// 		return int8Type
// 	case "int16":
// 		return int16Type
// 	case "int32":
// 		return int32Type
// 	case "int64":
// 		return int64Type
// 	case "string":
// 		return stringType
// 	case "binary":
// 		return binaryType
// 	case "bool":
// 		return boolType
// 	default:
// 		panic(fmt.Sprintf("type %s does not exist", basicType))
// 	}
// }

func ReadUntilNextSepToken(code string) (string, tokenEnum, string){
	hasSeparator := strings.ContainsAny(code, "{}():\n")
	if !hasSeparator{
		return code, EOF, ""
	}
	var builder strings.Builder
	builder.Grow(len(code))
	for i,r := range code{
		i = i +1
		switch r{
		case '{':
			return builder.String(), openBraces, code[i:]
		case '}':
			return builder.String(), closeBraces, code[i:]
		case '(':
			return builder.String(), openParentheses, code[i:]
		case ')':
			return builder.String(), closeParentheses, code[i:]
		case ':':
			return builder.String(), colon, code[i:]
		case '\n':
			return builder.String(), newLine, code[i:]
		default:
			builder.WriteRune(r)
		}
	}
	panic("unexpected EOF?")
}

func Tokenize(code string) []Token{
	x := make([]Token, 0)
	val, token, rest := ReadUntilNextSepToken(code)
	for token != EOF{
		val = strings.TrimSpace(val)
		if val != ""{
			var tokenType tokenEnum
			if strings.ContainsAny(strings.TrimSpace(string(val[0])), "0123456789"){
				_, err := strconv.Atoi(strings.TrimSpace(val))
				if err == nil{
					tokenType = number
				}else{
					println("could not parse to int")
					panic("could not parse to int")
				}
				
			}else{
				tokenType = text
			}
			x = append(x, Token{Token: tokenType, Option: val})
		}
		// println(strings.TrimSpace(val))
		// if len(x) > 0{
		// 	println(tnum[x[len(x)-1].Token])
		// }
		x = append(x, Token{Token:token, Option: ""})
		val, token, rest = ReadUntilNextSepToken(rest)
	}
	// for _, t := range x{
	// 	println(tnum[t.Token])
	// }
	return x
}

func expectToken(tokens []Token, expect tokenEnum, skipToken tokenEnum) ([]Token, Token) {
	for i, token := range tokens{
		if token.Token == expect{
			return tokens[i+1:], token
		}
		if token.Token == skipToken{
			continue
		}

		println("expected",tnum[expect])
		println("found",tnum[token.Token])
		panic("unexpected token")
	}
	panic("token not found")
}

func expectColumn(tokens []Token) ([]Token, ColunType){
	var name string
	var columnType BasicTypes
	var columnTypeSize int
	newtokens, token := expectToken(tokens, text, newLine)
	name = token.Option
	newtokens, _ = expectToken(newtokens, colon, None)
	newtokens, token = expectToken(newtokens, text, None)
	columnType = getBasicType(token.Option)
	// println("type ", columnType)
	if columnType == stringType || columnType == binaryType{
		var value Token
		newtokens, _ = expectToken(newtokens, openParentheses, None)
		newtokens, value = expectToken(newtokens, number, None)
		columnTypeSize,_ = strconv.Atoi(value.Option)
		newtokens, _ = expectToken(newtokens, closeParentheses, None)
	}
	newtokens, _ = expectToken(newtokens, newLine, None)
	return newtokens, ColunType{Name: name, Type: columnType, TypeOptions: columnTypeSize}
}

func expectTable(tokens []Token) ([]Token, Table){
	var tableName string
	columns := make([]ColunType, 0)
	newtokens, token  := expectToken(tokens, text, newLine)
	tableName = token.Option
	newtokens, _ = expectToken(newtokens, openBraces, newLine)
	for{
		var column ColunType
		newtokens, column = expectColumn(newtokens)
		columns = append(columns, column)
		if newtokens[0].Token == closeBraces{
			break
		}
		// println(tnum[newtokens[0].Token])
	}
	return newtokens, Table{
		TableName: tableName,
		Coluns: columns,
	}
}

func TokenParser(tokens []Token) Table{
	_, tt := expectTable(tokens)
	return tt
	// println(tt.TableName)
	// println(tt.Coluns[6].Type)
}