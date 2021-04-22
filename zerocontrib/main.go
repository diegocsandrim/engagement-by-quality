package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const inputFileName = "../data/raw.csv"
const outputFileName = "../data/all.csv"

// This program fills the gaps in inactive months.
// The gasps exists because the quality analyser only runs for months with commits.
// If in a certain month, no contributor made commit, tha analysis don't run.
// For this cases, the code have not changed.
// So the code metrics are equals to the previous month, and the contributors is equal to zero.
func main() {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		log.Panicf("could not open the file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Panicf("could not create the output file: %s", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	//Copy the header
	scanner.Scan()
	_, err = outputFile.WriteString(scanner.Text() + "\n")
	if err != nil {
		log.Panicf("could not write header to output file: %s", err)
	}

	var previousFields []string

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		err := fillGaps(outputFile, previousFields, fields)
		if err != nil {
			log.Panicf("could not write line to output file: %s", err)
		}

		previousFields = fields

		_, err = outputFile.WriteString(strings.Join(fields, ";") + "\n")
		if err != nil {
			log.Panicf("could not write line to output file: %s", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
}

func fillGaps(file *os.File, lineA, lineB []string) error {
	//Skips if its the first line
	if len(lineA) == 0 {
		return nil
	}

	//Skips if changed the project
	projectA := lineA[0]
	projectB := lineB[0]
	if projectA != projectB {
		return nil
	}
	dateLayout := "2006-01-02"

	dateA, err := time.Parse(dateLayout, lineA[1])
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	dateB, err := time.Parse(dateLayout, lineB[1])
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	for dateA = dateA.AddDate(0, 1, 0); dateA.Before(dateB); dateA = dateA.AddDate(0, 1, 0) {
		lineA[1] = dateA.Format(dateLayout)

		_, err = file.WriteString(strings.Join(lineA, ";") + "\n")
		if err != nil {
			log.Panicf("could not write line to output file: %s", err)
		}
	}

	return nil
}
