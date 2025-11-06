package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	viewCommand   = "view"
	searchCommand = "search"
	selectCommand = "select"
	filterCommand = "filter"
	countCommand  = "count"
)

func main() {
	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case viewCommand:
		view()
	case searchCommand:
		search()
	case selectCommand:
		sel()
	case filterCommand:
		filter()
	case countCommand:
		count()
	default:
		fmt.Println("Unknown command")
	}
}

func view() {
	cmdFlags := flag.NewFlagSet(viewCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		log.Fatal(err)
	}

	colCapacities := getCapCols(input)

	for _, row := range input {
		for col, cell := range row {
			fmt.Printf("%s", cell)
			for space := 0; space < colCapacities[col]-len(cell); space++ {
				fmt.Printf(" ")
			}
			fmt.Printf("| ")
		}
		fmt.Println()
	}
}

func search() {
	cmdFlags := flag.NewFlagSet(searchCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	query := cmdFlags.String("query", "", "Input a query to search")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range input {
		for _, cell := range row {
			if strings.Contains(cell, *query) {
				fmt.Printf("%v", row)
				break
			}
		}
	}
}

func sel() {
	cmdFlags := flag.NewFlagSet(selectCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	colsNames := cmdFlags.String("cols", "", "Input names of cols with \",\"")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		log.Fatal(err)
	}

	isColInChoose := getChosenCols(colsNames, input)
	colCapacities := getCapCols(input)

	for _, row := range input {
		for col, cell := range row {
			if isColInChoose[col] {
				fmt.Printf("%s", cell)
				for space := 0; space < colCapacities[col]-len(cell); space++ {
					fmt.Printf(" ")
				}
				fmt.Printf("| ")
			}
		}
		fmt.Println()
	}
}

func filter() {
	cmdFlags := flag.NewFlagSet(filterCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	colName := cmdFlags.String("col", "", "Input names of cols with \",\"")
	valueInCol := cmdFlags.String("val", "", "Input value of col")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		log.Fatal(err)
	}

	var colNum int
	for i, cell := range input[0] {
		if cell == *colName {
			colNum = i
			break
		}
	}

	for _, row := range input {
		if row[colNum] == *valueInCol {
			fmt.Printf("%v\n", row)
		}
	}
}

func count() {
	cmdFlags := flag.NewFlagSet(countCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		log.Fatal(err)
	}
	var count int
	for range input {
		count++
	}
	fmt.Println(count)
}

func getCapCols(input [][]string) []int {
	colCapacities := make([]int, len(input[0]))

	for _, row := range input {
		for i, cell := range row {
			colCapacities[i] = max(colCapacities[i], len(cell))
		}
	}

	return colCapacities
}

func getChosenCols(colsNames *string, input [][]string) []bool {
	cols := strings.Split(*colsNames, ",")

	isColInChoose := make([]bool, len(input[0]))
	for i, colName := range input[0] {
		for _, col := range cols {
			if col == colName {
				isColInChoose[i] = true
			}
		}
	}

	return isColInChoose
}

func parseFile(cmdFlags *flag.FlagSet, filename *string) ([][]string, error) {
	if err := cmdFlags.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	reader := csv.NewReader(file)
	input, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return input, nil
}
