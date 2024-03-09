package main

/* This is utilizing the flag package
in order to add arguments to the command line when calling the Go application.
The first argument is the file type, second argument is the name of the file and third argument describes what code does.
*/

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	// create an argument to enter the csv filename and don't accept any file that isn't named problems or ending with '.csv'
	csvFileName := flag.String("csv", "problems.csv", "a csv file that stores problems in the format: 'problem:answer' ")
	timelimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	// must be called after all flags to parse the value
	flag.Parse()

	// using the os package, open the file by its filename and return the error based on the following conditions:
	file, err := os.Open(*csvFileName)

	// if there's an error, print failed to open the CSV file.
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}

	// create a reader that reads the file recieved from the os.Open()
	r := csv.NewReader(file)

	// a slice of slices, usually in other programming languages representing a 2D array.
	lines, err := r.ReadAll()

	// if there's an error exit with the following string saying there was a failure to parse the CSV file
	if err != nil {
		exit("Failed to parse the provided CSV File")
	}

	// parse the lines into a variable named problems which is a slice of our struct type problem
	// then store the main logic for the program in a function named check Answers which takes the parsed problems and
	// returns the correctly answered ones in the form of an int named correct_answers
	// then print how many we scored out of the length of the problems.

	problems := parseLines(lines)
	check_answers(problems, timelimit)
}

func check_answers(problems []problem, timelimit *int) {

	timer := time.NewTimer(time.Duration(*timelimit) * time.Second)

	// initalize a variable named correct
	correct := 0

	// for each problem in the range of problems, iterate by one and ask the user for input on what the answer may be
	// if the answer is correct, increase the amount of correct answers, if not, do nothing.

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)

		// goroutine, similar to asynchronous method in C#, java or C++ allows concurrent methods
		// (execution simultaneously while other methods are executed)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d. \n", correct, len(problems))
			return

		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}

		}
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

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
