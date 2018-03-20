package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func startQuiz(c []problem, t int64) {
	fmt.Printf("GO! You have %v seconds!\n", t)
	timer := time.NewTimer(time.Second * time.Duration(t))
	correct := 0
	for i, val := range c {
		fmt.Printf("#%v: %v\n", i+1, val.q)
		answerChan := make(chan string)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			answerChan <- text
		}()

		select {
		case <-timer.C:
			fmt.Printf("Times Up! The %v sec. timer has expired.", t)
			fmt.Printf("\nYou got %v correct out of %v\n", correct, len(c))
			return
		case answer := <-answerChan:
			if strings.TrimSpace(answer) == val.a {
				correct++
			} else {
				fmt.Print("WRONG!!\n")
			}
		}
	}

	// This line will only print if the user makes it through all questions before the time limit.
	fmt.Printf("\nYou got %v correct out of %v\n", correct, len(c))
}

func main() {
	// Check for and parse command line flags for the problems csv file, defaulting to problems.csv in the current dir
	csvFile := flag.String("csv", "./problems.csv", "input csv file for problems and answers")
	timeLimit := flag.Int64("timer", 30, "Set a time limit on the quiz.")
	flag.Parse()

	// Read in the file and put into a csv reader type. Then read all into a new var (csvData)
	dat, err := os.Open(*csvFile)
	check(err)
	csvReader := csv.NewReader(dat)
	csvData, err := csvReader.ReadAll()
	check(err)

	lines := parseLines(csvData)

	fmt.Print("Press 'Enter' to start the quiz...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	startQuiz(lines, *timeLimit)
}
