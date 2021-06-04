package main

import (
	"encoding/csv"
	"io"
	"log"
)

type QuizProblem struct {
	Question string
	Answer   string
}

func ReadCSVProblems(csvProblems io.Reader) []QuizProblem {
	csvReader := csv.NewReader(csvProblems)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse input", err)
	}
	var problems []QuizProblem
	for _, line := range records {
		problem := QuizProblem{
			Question: line[0],
			Answer:   line[1],
		}
		problems = append(problems, problem)
	}
	return problems
}
