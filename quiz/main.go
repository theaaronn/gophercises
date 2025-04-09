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
	q []string
	a string
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			q: line[:len(line)-1],
			a: strings.TrimSpace(line[len(line)-1]),
		}
	}
	return problems
}

func readFile(name string) ([][]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	in := csv.NewReader(file)
	lines, err := in.ReadAll()
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	// Flag initialization
	fileNameFlag := flag.String("file", "problems.csv", "problems file name")
	timeFlag := flag.Int("time", 30, "timer for each question")
	randoFlag := flag.Bool("rando", false, "randomize question order")
	_ = randoFlag
	flag.Parse()

	lines, err := readFile(*fileNameFlag)
	if err != nil {
		fmt.Println("Error reading csv: ", err.Error())
	}
	problems := parseLines(lines)

	// Waiting input to start
	fmt.Println("All ready, press any key to start")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	score := 0
	timer := time.NewTimer(time.Duration(*timeFlag) * time.Second)

	for i, problem := range problems {
		fmt.Printf("%d.- %s: ", i+1, problem.q[:])
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer

		}()
		select {
		case <-timer.C:
			println()
			return
		case answer := <-answerCh:
			if answer == problem.a {
				score++
			}
		}
	}
	fmt.Printf("Final score: %d/%d", score, len(lines))
}
