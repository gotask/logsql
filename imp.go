package main

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/gotask/gost/stutil"
	"github.com/tealeg/xlsx"
)

type Reader interface {
	read(separator string) []string
}

func NewReader(path string) (Reader, error) {
	var rd Reader
	if strings.HasSuffix(path, ".csv") {
		rfile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		rd = &csvFile{csv.NewReader(rfile)}
	} else if strings.HasSuffix(path, ".xlsx") {
		xlFile, err := xlsx.OpenFile(path)
		if err != nil {
			return nil, err
		}
		rd = &xlsxFile{xlFile.Sheets, 0}
	} else {
		stfile, err := stutil.NewSTFile(path)
		if err != nil {
			return nil, err
		}
		rd = &txtFile{stfile}
	}
	return rd, nil
}

type txtFile struct {
	file *stutil.STFile
}

func (txt *txtFile) read(separator string) (res []string) {
	line, n := txt.file.ReadLine()
	for line == "" && n > 0 {
		line, n = txt.file.ReadLine()
	}
	if line == "" {
		return nil
	}
	if separator == "" {
		res = strings.Fields(line)
	}

	return strings.Split(line, separator)
}

type csvFile struct {
	file *csv.Reader
}

func (csv *csvFile) read(separator string) (res []string) {
	res, e := csv.file.Read()
	if e != nil {
		return nil
	}
	return
}

type xlsxFile struct {
	sheet []*xlsx.Sheet
	row   int
}

func (xlsx *xlsxFile) read(separator string) (res []string) {
	if len(xlsx.sheet) == 0 {
		return nil
	}
	rows := xlsx.sheet[0].Rows
	if xlsx.row >= len(rows) {
		return nil
	}
	cells := rows[xlsx.row].Cells
	res = make([]string, len(cells))
	for i, cell := range cells {
		res[i] = cell.String()
	}
	xlsx.row++
	return
}
