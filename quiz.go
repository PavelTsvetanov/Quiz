package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

func readProblems(filePath string) []QuizProblem {
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

func askQuestion(inputReader *bufio.Reader, problem QuizProblem) string {
	fmt.Print(problem.Question, ",")
	line, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatal("Unable to read user input", err)
	}
	return strings.TrimSuffix(line, "\n")
}

const defaultProblemsCSV = "data/problems.csv"

func main() {
	csvUsage := "a csv file in the format of 'question,answer' (default \"problems.csv\")"
	problemsCSV := flag.String("csv", defaultProblemsCSV, csvUsage)
	//limit := flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	inputReader := bufio.NewReader(os.Stdin)
	problems := readProblems(*problemsCSV)
	var correctProblems int
	for _, problem := range problems {
		if answer := askQuestion(inputReader, problem); answer == problem.Answer {
			correctProblems += 1
		}
	}
	fmt.Println(correctProblems, "/", len(problems), " correct")
}
