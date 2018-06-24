package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dinedal/textql/storage"
)

var (
	storageOpts = &storage.SQLite3Options{}
)

func IsThereDataOnStdin() bool {
	stat, statErr := os.Stdin.Stat()

	if statErr != nil {
		return false
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	} else {
		return false
	}
}

func main() {
	readFlag()
	if len(cmd.input) == 0 && !IsThereDataOnStdin() {
		flag.PrintDefaults()
		return
	}

	var inputSources []string
	if IsThereDataOnStdin() {
		inputSources = append(inputSources, "stdin")
	}
	inputSources = append(inputSources, cmd.input...)

	storage := storage.NewSQLite3StorageWithDefaults()

	for _, file := range inputSources {
		input, inputErr := NewInput(*cmd.hed, *cmd.col, *cmd.sep, file)
		if inputErr != nil {
			log.Printf("Unable to load file %v\n", file)
			continue
		}

		storage.LoadInput(input)
	}

	sqlStrings := strings.Split(*cmd.sql, ";")

	for _, sqlQuery := range sqlStrings {
		queryResults, queryErr := storage.ExecuteSQLString(sqlQuery)

		if queryErr != nil {
			log.Fatalln(queryErr)
		}

		if queryResults != nil {
			cols, colsErr := queryResults.Columns()
			if colsErr != nil {
				log.Fatalln(colsErr)
			}

			rawResult := make([][]byte, len(cols))
			result := make([]string, len(cols))
			dest := make([]interface{}, len(cols))
			for i := range cols {
				dest[i] = &rawResult[i]
			}

			for queryResults.Next() {
				queryResults.Scan(dest...)

				for i, raw := range rawResult {
					result[i] = string(raw)
					fmt.Print(result[i])
					if i != len(result)-1 {
						fmt.Print("\t")
					}
				}
				fmt.Print("\n")
			}
			queryResults.Close()
		}
	}
}

type cmdParam struct {
	sql    *string
	sep    *string
	hed    *bool
	col    *int
	input  []string
	output *string
}

var cmd cmdParam

func readFlag() {
	cmd.sql = flag.String("s", "", "SQL Statement(s) to run on the data")
	cmd.sep = flag.String("f", " ", "Input delimiter character between fields")
	cmd.hed = flag.Bool("h", false, "Treat input files as having the first row as a header row")
	cmd.col = flag.Int("c", 0, "specail the count of column")
	flag.Parse()
	cmd.input = flag.Args()
}
