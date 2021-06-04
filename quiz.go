package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultProblemsCSV       = "data/problems.csv"
	defaultTimerOutInSeconds = 30
	csvUsage                 = "a csv file in the format of 'question,answer' (default \"problems.csv\")"
	limitUsage               = "the time limit for the quiz in seconds (default 30)"
)

type QuizInputs struct {
	csv   string
	limit int
}

func getQuizInputs() QuizInputs {
	problemsCSV := flag.String("csv", defaultProblemsCSV, csvUsage)
	limit := flag.Int("limit", defaultTimerOutInSeconds, limitUsage)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	return QuizInputs{
		csv:   *problemsCSV,
		limit: *limit,
	}
}

func askQuestion(inputReader *bufio.Reader, problem QuizProblem) string {
	fmt.Print(problem.Question, ",")
	line, err := inputReader.ReadString('\n')
	if err != nil {
		log.Fatal("Unable to read user input", err)
	}
	return strings.TrimSuffix(line, "\n")
}

func StartQuiz(inputReader *bufio.Reader, problems []QuizProblem) <-chan int {
	answers := make(chan int)
	go func() {
		for _, problem := range problems {
			if answer := askQuestion(inputReader, problem); answer == problem.Answer {
				answers <- 1
			} else {
				answers <- 0
			}
		}
	}()
	return answers
}

func CalculateUserScore(userAnswers <-chan int, limit <-chan time.Time) int {
	var score int
	for {
		select {
		case answer := <-userAnswers:
			score += answer
		case <-limit:
			return score
		}
	}
}

func main() {
	quizInputs := getQuizInputs()

	userInputReader := bufio.NewReader(os.Stdin)

	fmt.Print("Press enter to start quiz")
	_, err := userInputReader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read user input", err)
	}

	file, err := os.Open(quizInputs.csv)
	if err != nil {
		log.Fatal("Unable to read input file: "+quizInputs.csv, err)
	}
	defer file.Close()
	problems := ReadCSVProblems(file)

	userAnswers := StartQuiz(userInputReader, problems)

	timeOut := time.After(time.Duration(quizInputs.limit) * time.Second)
	score := CalculateUserScore(userAnswers, timeOut)

	fmt.Printf("\n%d/%d correct", score, len(problems))
}
