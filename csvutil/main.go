package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	viewCommand   = "view"
	searchCommand = "search"
	selectCommand = "select"
	filterCommand = "filter"
	countCommand  = "count"
	sortCommand   = "sort"
	addCommand    = "add"
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
	case sortCommand:
		sortOnCols()
	case addCommand:
		add()
	default:
		fmt.Println("Unknown command")
	}
}

// ./main.exe view -f data.csv
func view() {
	cmdFlags := flag.NewFlagSet(viewCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
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

// ./main.exe search -f data.csv -query "Ivan"
func search() {
	cmdFlags := flag.NewFlagSet(searchCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	query := cmdFlags.String("query", "", "Input a query to search")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
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

// ./main.exe select -f data.csv -cols "Name,Email"
func sel() {
	cmdFlags := flag.NewFlagSet(selectCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	colsNames := cmdFlags.String("cols", "", "Input names of cols with \",\"")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
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

// ./main.exe filter -f data.csv -col "Age" -val "30"
func filter() {
	cmdFlags := flag.NewFlagSet(filterCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	colName := cmdFlags.String("col", "", "Input names of cols with \",\"")
	valueInCol := cmdFlags.String("val", "", "Input value of col")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
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

// ./main.exe count -f data.csv
func count() {
	cmdFlags := flag.NewFlagSet(countCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}

	fmt.Println(len(input))
}

// ./main.exe sort -f data.csv -col "Salary" -order desc
func sortOnCols() {
	cmdFlags := flag.NewFlagSet(sortCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	colName := cmdFlags.String("col", "", "Input names of col")
	order := cmdFlags.String("order", "desc", "Input order(desc/inc)")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}
	colNum, err := findColNumFromName(input, colName)
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}

	file, err := os.Create("tmp.csv")
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
		if err := os.Remove(*filename); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
		if err := os.Rename("tmp.csv", *filename); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	if err := sortAndWrite(writer, input, colNum, order); err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}
}

// ./main.exe add -f data.csv -data "John Doe,john@example.com,25000"
func add() {
	cmdFlags := flag.NewFlagSet(addCommand, flag.ExitOnError)
	filename := cmdFlags.String("f", "", "Input CSV file")
	data := cmdFlags.String("data", "", "Input data with \",\"")

	input, err := parseFile(cmdFlags, filename)
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}

	file, err := os.Create("tmp.csv")
	if err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
		if err := os.Remove(*filename); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
		if err := os.Rename("tmp.csv", *filename); err != nil {
			cmdFlags.Usage()
			log.Fatal(err)
		}
	}()

	dataSplit := strings.Split(*data, ",")
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := addLineToTMP(writer, input, dataSplit); err != nil {
		cmdFlags.Usage()
		log.Fatal(err)
	}
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

func sortAndWrite(writer *csv.Writer, input [][]string, colNum int, order *string) error {
	slices.SortFunc(input[1:], func(i []string, j []string) int {
		numI, errI := strconv.Atoi(i[colNum])
		numJ, errJ := strconv.Atoi(j[colNum])

		if errI == nil && errJ == nil {
			if numI < numJ {
				return -1
			}
			if numI > numJ {
				return 1
			}
			if numI == numJ {
				return 0
			}
		}
		timeI, errI := time.Parse("02.01.2006", i[colNum])
		timeJ, errJ := time.Parse("02.01.2006", j[colNum])
		if errI == nil && errJ == nil {
			if timeI.After(timeJ) {
				return 1
			}
			if timeI.Before(timeJ) {
				return -1
			}
			if timeI == timeJ {
				return 0
			}
		}
		if i[colNum] < j[colNum] {
			return -1
		}
		if i[colNum] > j[colNum] {
			return 1
		}
		return 0
	})

	if err := writer.Write(input[0]); err != nil {
		return err
	}
	switch *order {
	case "desc":
		for i := len(input) - 1; i >= 1; i-- {
			if err := writer.Write(input[i]); err != nil {
				return err
			}
		}
	case "inc":
		for i := 1; i < len(input); i++ {
			if err := writer.Write(input[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func findColNumFromName(input [][]string, colName *string) (int, error) {
	colNum := -1
	if len(input) <= 2 {
		return colNum, errors.New("Input CSV file is empty ")
	}
	for i, cell := range input[0] {
		if cell == *colName {
			colNum = i
			break
		}
	}
	if colNum == -1 {
		return colNum, errors.New("Column Name not found ")
	}
	return colNum, nil
}

func addLineToTMP(writer *csv.Writer, input [][]string, dataSplit []string) error {
	for _, row := range input {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	if len(input) == 0 {
		if err := writer.Write(dataSplit); err != nil {
			return err
		}
	} else if len(dataSplit) != len(input[0]) {
		return errors.New("Input CSV file does not have the same number of columns ")
	} else if !validateTypes(input[0], dataSplit) {
		return errors.New("Input CSV file does not have the correct type ")
	} else if err := writer.Write(dataSplit); err != nil {
		return err
	}
	return nil
}

func validateTypes(topics []string, data []string) bool {
	typesTopics := lineToTypes(topics)
	typesData := lineToTypes(data)

	for i := range typesTopics {
		if typesTopics[i] != typesData[i] {
			return false
		}
	}
	return true
}

func lineToTypes(line []string) []string {
	types := make([]string, len(line))
	for i, cell := range line {
		if _, err := strconv.Atoi(cell); err == nil {
			types[i] = "int"
			continue
		}
		if _, err := time.Parse("02.01.2006", cell); err == nil {
			types[i] = "time"
			continue
		}
		types[i] = "string"
	}
	return types
}
