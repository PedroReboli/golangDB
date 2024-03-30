package godb

import (
	"fmt"
	"strconv"
	"strings"
)

type MacroOperation int
const (
	assignment MacroOperation = iota
	condition
	forLoop
	whileLoop
	ifCond
	elseCond
	panicExit
	constant 
	variable
	greaterComp 
	lessComp
	equalComp
	notEqualComp
	lessEqualComp
	greaterEqualComp
	andComp
	orComp
	notComp
)
// type conditionStruct struct{
// 	leftOp HighByteCodeType
// 	rightOp HighByteCodeType
// 	compType comp
// }
/*
select id, nome from tbl where id != 0 and nome = "andre"
*/

type operation struct{
	operationType MacroOperation
	option string 
}

func whereParser(tokens []Token) []operation{
	newToken, ops := conditionParser(tokens)
	var newOps = make([]operation, 0)
	for len(newToken) > 0{
		var token Token
		newToken, token = expectToken(newToken, text, newLine)
		if token.Token == text{
			if token.Option == "and"{
				ops = append(ops, operation{operationType: andComp, option: ""})
			} else if token.Option == "or" {
				ops = append(ops, operation{operationType: orComp, option: ""})
			}
		}
		newToken, newOps = conditionParser(newToken)
		ops = append(ops, newOps...)
		// newOps = make([]operation, 0)
	}
	return ops
}

func conditionParser(tokens []Token) ([]Token, []operation){
	ops := make([]operation,0)
	newToken, firstVariable := expectMultipleTokens(tokens, []tokenEnum{text, number}, []tokenEnum{newLine})
	newToken, firstComp := expectMultipleTokens(newToken,[]tokenEnum{equal, lessThan, greaterThan, exclamation}, []tokenEnum{newLine})
	newToken, secondCompOrVariable := expectMultipleTokens(newToken,[]tokenEnum{equal, text, number}, []tokenEnum{newLine})
	var SecondVariable Token
	var ComparisonType MacroOperation
	if firstComp.Token == equal || firstComp.Token == exclamation{
		if secondCompOrVariable.Token != equal{
			panic("equal was exepected found"+tnum[secondCompOrVariable.Token])
		}
		if firstComp.Token == equal{
			ComparisonType = equalComp
		}else{
			ComparisonType = notEqualComp
		}
		newToken, SecondVariable = expectMultipleTokens(newToken, []tokenEnum{text, number}, []tokenEnum{newLine})
	}else if secondCompOrVariable.Token == equal{
		switch firstComp.Token{ // this switch doesn't make sense
		case lessThan:
			ComparisonType = lessEqualComp
		case greaterThan:
			ComparisonType = greaterEqualComp
		default:
			panic("equal was exepected found"+tnum[firstComp.Token]) 
		}
		newToken, SecondVariable = expectMultipleTokens(newToken, []tokenEnum{text, number}, []tokenEnum{newLine})
	} else {
		switch firstComp.Token{ // this switch doesn't make sense
		case lessThan:
			ComparisonType = lessComp
		case greaterThan:
			ComparisonType = greaterComp
		}
		SecondVariable = secondCompOrVariable
	}
	if firstVariable.Token == text{
		ops = append(ops, operation{ operationType: variable, option: firstVariable.Option })
	}else{
		ops = append(ops, operation{ operationType: constant, option: firstVariable.Option })
	}
	ops = append(ops, operation{ operationType: ComparisonType, option: ""})
	if SecondVariable.Token == text{
		ops = append(ops, operation{ operationType: variable, option: SecondVariable.Option })
	}else{
		ops = append(ops, operation{ operationType: constant, option: SecondVariable.Option })
	}
	return newToken,ops
}


func QueryReadUntilNextSepToken(code string) (string, tokenEnum, string){
	hasSeparator := strings.ContainsAny(code, "{}():\n-+!=, ")
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
		case '-':
			return builder.String(), minus, code[i:]
		case '+':
			return builder.String(), plus, code[i:]
		case '!':
			return builder.String(), exclamation, code[i:]
		case '=':
			return builder.String(), equal, code[i:]
		case ',':
			return builder.String(), comma, code[i:]
		case ' ':
			return builder.String(), newLine, code[i:]
		case '<':
			return builder.String(), lessThan, code[i:]
		case '>':
			return builder.String(), greaterThan, code[i:]
		default:
			builder.WriteRune(r)
		}
	}
	panic("unexpected EOF?")
}
func QueryTokenize(code string) []Token{
	code = strings.TrimSpace(code)
	x := make([]Token, 0)
	val, token, rest := QueryReadUntilNextSepToken(code)
	for {
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
		if token == EOF{
			break
		}
		x = append(x, Token{Token:token, Option: ""})
		val, token, rest = QueryReadUntilNextSepToken(rest)
	}
	// for _, t := range x{
	// 	println(tnum[t.Token])
	// }
	return x
}

func ExpectQuery(tokens []Token) selectQuery{
	newtokens, token  := expectToken(tokens, text, newLine)
	// switch token.Option{
	// case "select":
		
	// case "insert":
	// 	panic("to be implemented")
	// }
	if token.Option != "select"{
		panic("select was expected")
	}
	coluns := make([]string,0)
	for{
		newtokens, token  = expectToken(newtokens, text, newLine)
		coluns = append(coluns, token.Option)
		if newtokens[0].Token != newLine && newtokens[0].Token != comma{
			panic(fmt.Sprintf("expected newLine or comma or text, found %s", tnum[newtokens[0].Token]))
		}
		if newtokens[0].Token == text{
			break
		}
		ntokens, nextToken := expectMultipleTokens(newtokens, []tokenEnum{text,comma}, []tokenEnum{newLine})
		if nextToken.Token == comma{
			newtokens = ntokens
		}else{
			break
		}
	}
	newtokens, token = expectToken(newtokens, text, newLine)
	if token.Option != "from"{
		panic("expected token FROM")
	}
	newtokens, tableToken := expectToken(newtokens, text, newLine)
	
	whereConditionTokens, whereToken := expectToken(newtokens, text, newLine)
	
	if whereToken.Option != "where"{
		panic("expected token where")
	}
	whereOp := whereParser(whereConditionTokens)
	return selectQuery{
		selectField: coluns,
		table: tableToken.Option,
		whereMacroCode: whereOp,
	}
}