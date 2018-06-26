package main

import (
	"log"
	"strconv"
)

var (
	fileIndex int = -1
)

type Input struct {
	hasHeader bool
	columnNum int
	separator string
	path      string
	file      Reader
	name      string
	header    []string
	firstRow  []string
}

func NewInput(hasHeader bool, columnNum int, separator, path string) (*Input, error) {
	fileIndex++
	rd, e := NewReader(path)
	if e != nil {
		log.Fatalln(e)
		return nil, e
	}
	input := &Input{hasHeader, columnNum, separator, path, rd, "t" + strconv.Itoa(fileIndex), nil, nil}
	input.firstRow = input.readLine()
	if columnNum <= 0 {
		input.columnNum = len(input.firstRow)
	}

	if hasHeader {
		input.header = input.firstRow
		input.firstRow = nil
	} else {
		input.header = make([]string, input.columnNum)
		for i := 0; i < input.columnNum; i++ {
			input.header[i] = "c" + strconv.Itoa(i)
		}
	}

	return input, nil
}

func (txt *Input) readLine() (res []string) {
	res = txt.file.read(txt.separator)
	if res == nil {
		return
	}

	if txt.columnNum > 0 {
		if txt.columnNum > len(res) {
			for i := len(res); i < txt.columnNum; i++ {
				res = append(res, "")
			}
		} else if txt.columnNum < len(res) {
			str := res[txt.columnNum-1]
			for i := txt.columnNum; i < len(res); i++ {
				str += res[i]
			}
			res = res[0:txt.columnNum]
			res[txt.columnNum-1] = str
		}
	}
	return res
}

func (txt *Input) Name() string {
	return txt.name
}

func (txt *Input) SetName(name string) {
	txt.name = name
}

// ReadRecord reads a single record from the CSV. Always returns successfully.
// If the record is empty, an empty []string is returned.
// Record expand to match the current row size, adding blank fields as needed.
// Records never return less then the number of fields in the first row.
// Returns nil on EOF
// In the event of a parse error due to an invalid record, it is logged, and
// an empty []string is returned with the number of fields in the first row,
// as if the record were empty.
func (txt *Input) ReadRecord() []string {
	if txt.firstRow != nil {
		row := txt.firstRow
		txt.firstRow = nil
		return row
	}

	return txt.readLine()
}

// Header returns the header of the csvInput. Either the first row if a header
// set in the options, or c#, where # is the column number, starting with 0.
func (txt *Input) Header() []string {
	return txt.header
}
