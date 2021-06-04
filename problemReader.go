package main

import (
	"encoding/csv"
	"log"
	"os"
)

func readCsvFile(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file: "+filePath, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse CSV file: "+filePath, err)
	}

	return records
}

type QuizProblem struct {
	Question string
	Answer   string
}

func ReadProblems(filePath string) []QuizProblem {
	lines := readCsvFile(filePath)

	var problems []QuizProblem
	for _, line := range lines {
		problem := QuizProblem{
			Question: line[0],
			Answer:   line[1],
		}
		problems = append(problems, problem)
	}
	return problems
}
