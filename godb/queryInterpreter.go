package godb

import (
	"io"
	"strconv"
)
type selectQuery struct{
	selectField []string
	table string
	whereMacroCode []operation
}

func (sq *selectQuery) Execute(table *Table) string{
	// table := table2
	_, err := table.File.Seek(int64(table.headSize),io.SeekStart)
	if err != nil{
		panic(err)
	}
	rowBytes := make([]byte, table.rowSize+1)
	colunsShow := make([]ColunType,0)
	colunsSelect := make(map[string]ColunType,0)
	for _,field := range sq.selectField{
		for _, col := range table.Coluns{
			if col.Name == field{
				colunsShow = append(colunsShow, col)
				colunsSelect[col.Name] = col
			}
		}
	}
	for _,field := range variablesOfMacroCode(sq.whereMacroCode){
		for _, col := range table.Coluns{
			if col.Name == field{
				colunsSelect[col.Name] = col
			}
		}
	}

	variables := make(map[string]string)
	output := "" 
	for{
		_, err := table.File.Read(rowBytes)
		if err != nil{
			break
		}
		for _, col := range colunsSelect{
			variables[col.Name] = col.ByteToString(rowBytes)
		}

		if executeMacroCode(&variables, sq.whereMacroCode){
			for _,cs := range colunsShow{
				v := variables[cs.Name]
				output += (cs.Name+ " : "+v+"\n")
				// println(cs.Name, " : "+v)
			}
		}
		
	}
	return output
}

func variablesOfMacroCode(MacroCode []operation) []string{
	variables := make([]string,0)
	for _,Mc := range MacroCode{
		if Mc.operationType == variable{
			variables = append(variables, Mc.option)
		}
	}
	return variables
}

func executeCompMacroCode(variable *map[string]string, MacroCode []operation) bool{
	var left int
	var right int
	 
	if MacroCode[0].operationType == constant{
		left,_ = strconv.Atoi(MacroCode[0].option)  	
	} else{
		variable ,ok := (*variable)[MacroCode[0].option]
		if !ok{
			panic("variable not found: "+variable)
		}
		left,_ = strconv.Atoi(variable)
	}
	if MacroCode[2].operationType == constant{
		right,_ = strconv.Atoi(MacroCode[2].option)  	
	} else{
		variable ,ok := (*variable)[MacroCode[2].option]
		if !ok{
			panic("variable not found: "+variable)
		}
		right,_ = strconv.Atoi(variable)
	}

	switch MacroCode[1].operationType{
	case equalComp:
		return left == right
	case lessComp:
		return left < right
	case greaterComp:
		return left > right
	case greaterEqualComp:
		return left >= right
	case lessEqualComp:
		return left <= right
	}
	panic("unreachable code reach")
}

func executeMacroCode(variable *map[string]string, MacroCode []operation) bool{
	i := 0
	lastOP := true
	BoolOp := andComp
	for{
		thisOP := executeCompMacroCode(variable, MacroCode[i:i+3])
		if BoolOp == andComp{
			lastOP = lastOP && thisOP
		}
		if BoolOp == orComp{
			lastOP = lastOP || thisOP
		}
		i += 3
		if len(MacroCode) <= i+1{
			return lastOP
		}
		BoolOp = MacroCode[i].operationType
		i += 1
	}
}
