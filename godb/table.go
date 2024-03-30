package godb

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type serialization struct{
	pos uint32
	internalBytePos uint8
}

type ColunType struct{
	index uint16
	Name string
	Type BasicTypes
	TypeOptions int
	serialize *serialization
}

func (ct *ColunType) TypeToByteArray(value any) []byte{
	switch ct.Type{
	case uint8Type:
		v,ok := value.(uint8)
		if !ok{
			panic("value is not uint8")
		}
		return []byte{v}
	case uint16Type:
		v,ok := value.(uint16)
		if !ok{
			panic("value is not uint16")
		}
		return binary.BigEndian.AppendUint16(make([]byte, 0), v)
	case uint32Type:
		v,ok := value.(uint32)
		if !ok{
			panic("value is not uint32")
		}
		return binary.BigEndian.AppendUint32(make([]byte, 0), v)
	case uint64Type:
		v,ok := value.(uint64)
		if !ok{
			panic("value is not uint64")
		}
		return binary.BigEndian.AppendUint64(make([]byte, 0), v)
	case int8Type:
		v,ok := value.(int8)
		if !ok{
			panic("value is not int8")
		}
		return []byte{uint8(v)}
	case int16Type:
		v,ok := value.(int16)
		if !ok{
			panic("value is not int16")
		}
		return binary.BigEndian.AppendUint16(make([]byte, 0), uint16(v))
	case int32Type:
		v,ok := value.(int32)
		if !ok{
			panic("value is not int32")
		}
		return binary.BigEndian.AppendUint32(make([]byte, 0), uint32(v))
	case int64Type:
		v,ok := value.(int64)
		if !ok{
			panic("value is not int64")
		}
		return binary.BigEndian.AppendUint64(make([]byte, 0), uint64(v))
	case stringType:
		v,ok := value.(string)
		if !ok{
			panic("value is not string")
		}
		array := make([]byte, ct.TypeOptions)
		copy(array,[]byte(v))
		return array
	case binaryType:
		v,ok := value.([]byte)
		if !ok{
			panic("value is not byte array")
		}
		array := make([]byte, ct.TypeOptions)
		copy(array,v)
		return array
	case boolType:
		v,ok := value.(bool)
		if !ok{
			panic("value is not byte array")
		}
		return []byte{bool2int(v) << ct.serialize.internalBytePos}
	default:
		panic("value not found???????")
	}
}

func (ct *ColunType) ByteSize() uint32{
	switch ct.Type{
	case uint8Type:
		return 1
	case uint16Type:
		return 2
	case uint32Type:
		return 4
	case uint64Type:
		return 8
	case int8Type:
		return 1
	case int16Type:
		return 2
	case int32Type:
		return 4
	case int64Type:
		return 8
	case stringType:
		return uint32(ct.TypeOptions)
	case binaryType:
		return uint32(ct.TypeOptions)
	case boolType:
		return 1
	default:
		panic("value not found???????")
	}
}
func (ct *ColunType) ByteToType(rowByte []byte) any{
	switch ct.Type{
	case uint8Type:
		return uint8(rowByte[ct.serialize.pos])
	case uint16Type:
		
		return GetUint16At(rowByte, uint64(ct.serialize.pos))
	case uint32Type:

		return GetUint32At(rowByte, uint64(ct.serialize.pos))
	case uint64Type:
		
		return GetUint64At(rowByte, uint64(ct.serialize.pos))
	case int8Type:
		
		return int8( GetByteAt(rowByte, uint64(ct.serialize.pos)))
	case int16Type:
		return int16(GetUint16At(rowByte, uint64(ct.serialize.pos)))
	case int32Type:
		return int32(GetUint32At(rowByte, uint64(ct.serialize.pos)))
	case int64Type:
		return int64(GetUint64At(rowByte, uint64(ct.serialize.pos)))
	case stringType:
		return GetStringAt(rowByte, uint64(ct.serialize.pos), int16(ct.TypeOptions))
	case binaryType:
		return GetBytesAt(rowByte, uint64(ct.serialize.pos), int16(ct.TypeOptions))
	case boolType:
		B := GetByteAt(rowByte, uint64(ct.serialize.pos))
		return (B>>ct.serialize.internalBytePos & 1) == 1
	default:
		panic("value not found???????")
	}
}

func (ct *ColunType) ByteToString(rowByte []byte) string{
	anyValue := ct.ByteToType(rowByte)
	switch ct.Type{
	case uint8Type:
		return fmt.Sprint(anyValue.(uint8))
	case uint16Type:
		
		return fmt.Sprint(anyValue.(uint16))
	case uint32Type:

		return fmt.Sprint(anyValue.(uint32))
	case uint64Type:
		
		return fmt.Sprint(anyValue.(uint64))
	case int8Type:
		
		return fmt.Sprint(anyValue.(int8))
	case int16Type:
		return fmt.Sprint(anyValue.(int16))
	case int32Type:
		return fmt.Sprint(anyValue.(int32))
	case int64Type:
		return fmt.Sprint(anyValue.(int64))
	case stringType:
		return fmt.Sprint(anyValue.(string))
	case binaryType:
		return fmt.Sprint(anyValue.([]byte))
	case boolType:
		
		return fmt.Sprint( anyValue.(bool))
	default:
		panic("value not found???????")
	}
}

type Table struct{
	TableName string
	Coluns []ColunType
	File *os.File
	rowSize uint32
	headSize uint32
}

func(t *Table) Init(file string) {
	byteOutput := make([]byte, 0)
	// x := len(t.TableName)
	byteOutput = binary.BigEndian.AppendUint16(byteOutput, uint16(len(t.TableName)))
	byteOutput = append(byteOutput, t.TableName...)
	byteOutput = binary.BigEndian.AppendUint16(byteOutput, uint16(len(t.Coluns)))
	for _, coluns := range t.Coluns{
		nameByte := []byte(coluns.Name)
		byteOutput = binary.BigEndian.AppendUint16(byteOutput, uint16(len(nameByte)))
		byteOutput = append(byteOutput, nameByte...)
		// byteOutput = append(byteOutput, )
		// byteOutput = binary.BigEndian.AppendUint16(byteOutput, uint16(coluns.Type))
		byteOutput = append(byteOutput, byte(coluns.Type))
		byteOutput = binary.BigEndian.AppendUint32(byteOutput, uint32(coluns.TypeOptions))
	}
	byteSize := make([]byte, 2)
	binary.BigEndian.PutUint16(byteSize, uint16(len(byteOutput))+1)
	t.headSize =  uint32(len(byteOutput)+1)
	byteOutput = append(byteSize, byteOutput...)
	// return byteOutput
	
	// fileHandler, err := os.OpenFile(file, os.O_CREATE & os.O_APPEND & os.O_RDWR, 0777)
	fileHandler, err := os.Create(file)
	if err != nil{
		panic(err)
	}
	_, err = fileHandler.Write(byteOutput)
	if err != nil{
		panic(err)
	}
	t.CalculateSerizalization()
	t.File = fileHandler
}
func OpenFile(file string) Table{
	fileHand, err := os.OpenFile(file, os.O_RDWR | os.O_APPEND, 0)
	if err != nil{
		panic(err)
	}
	headSizeByte := make([]byte, 2)
	fileHand.ReadAt(headSizeByte,0)
	
	ByteReader := NewByteReader(headSizeByte)
	headerSize :=ByteReader.GetUint16()
	// println(headerSize)
	// t.headSize =  uint16(len(byteOutput))+1
	header := make([]byte, headerSize)
	fileHand.ReadAt(header,2)
	table := ReadTable(header)
	
	// bufio.new
	table.File = fileHand
	table.headSize = uint32(headerSize) + 1
	// table.CalculateSerizalization()
	return table
}
// func 
func ReadTable(headerData []byte) Table{
	byteReader := NewByteReader(headerData)
	tableNameSize := byteReader.GetUint16()
	// println(tableNameSize)
	tableName := byteReader.GetString(int16(tableNameSize))
	colunsCount := byteReader.GetUint16()
	table := Table{TableName: tableName, Coluns: make([]ColunType, 0)}
	for i := uint16(0); i < colunsCount; i+=1{
		colunNameSize := byteReader.GetUint16()
		colunName := byteReader.GetString(int16(colunNameSize))
		colunType := byteReader.GetByte()
		colunOption := byteReader.GetUint32()
		colun := ColunType{
			Name: colunName,
			Type: BasicTypes(colunType),
			TypeOptions: int(colunOption),
		}
		table.Coluns = append(table.Coluns, colun)
	}
	table.CalculateSerizalization()
	return table
}

func (t* Table)CalculateSerizalization(){
	// coluns := make([]ColunType,0)
	var i uint16 = 0
	var size uint32 = 0
	has_bool := false
	for ii, c := range t.Coluns{
		if c.Type == boolType{
			has_bool = true
			continue
		}
		t.Coluns[ii].serialize = &serialization{
			pos: size,
			internalBytePos: 0,
		}
		t.Coluns[ii].index = i
		i += 1
		size += c.ByteSize()
	}
	if !has_bool{
		return
	}
	var bytePos uint8 = 0
	for ii, c := range t.Coluns{
		if c.Type != boolType{
			continue
		}
		t.Coluns[ii].serialize = &serialization{
			pos: size,
			internalBytePos: bytePos,
		}
		t.Coluns[ii].index = i

		i += 1
		bytePos += 1
		if bytePos == 8{
			bytePos = 0
			size += 1
		}
	}
	t.rowSize = size
}

func(t* Table) InsertRow(values []any){
	t.CalculateSerizalization()
	byt := make([]byte, t.rowSize+1)
	// println(len(values))
	// println(len(t.Coluns))
	size := 0
	for _,z := range Zip(t.Coluns, values){
		if z.First.Type == boolType{
			p := z.First.TypeToByteArray(z.Second)[0]
			x := byt[z.First.serialize.pos] | p
			byt[z.First.serialize.pos] = x
		}else{
			copy(byt[z.First.serialize.pos:z.First.serialize.pos+z.First.ByteSize()], z.First.TypeToByteArray(z.Second))
		}
		
		size += len(z.First.TypeToByteArray(z.Second))
		// byt = append(byt, z.First.TypeToByteArray(z.Second)...)
	}
	// println("expected len", t.rowSize)
	// println("found len", size)
	(*t.File).Seek(0,io.SeekEnd)
	_, _ = (*t.File).Write(byt)
}
func(t* Table) Close(){
	(*t.File).Close()
}

func(t* Table) ReadAllRows() {
	
	// t.row_size
	_, _ = t.File.Seek(int64(t.headSize),io.SeekStart)
	// println(i)
	// println(err)
	rowBytes := make([]byte, t.rowSize+1)
	// rowParsed := make([]any, len(t.Coluns))
	for{
		_, err := t.File.Read(rowBytes)
		println("-------------------------")
		if err != nil{
			// panic(err)
			return
		}
		// println(rowBytes)
		for _, col := range t.Coluns{
			xx := col.ByteToString(rowBytes)
			println(col.Name,":",xx)
		}
		// println()
	}
	

}