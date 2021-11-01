package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type QuizContent struct {
	Question string
	Answer   string
}

func main() {

	// Receving flag from CLI
	CSVFileFlag := flag.String("csv", "problems.csv", "A csv file in the format of 'question,answer'")
	LimitFlag := flag.Int("limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	// Read the CSV file
	records, err := readData(*CSVFileFlag)
	if err != nil {
		log.Fatal(err)
	}

	QtyQuestion := 0
	QtyRightAnswer := 0
	// Loop for each CSV line

	Timer := time.NewTimer(time.Duration(*LimitFlag) * time.Second)

QuizLoop:
	// Parsing CSV lines
	for i, record := range records {
		QtyQuestion++
		quiz := QuizContent{
			Question: record[0],
			Answer:   record[1],
		}
		fmt.Printf("Problem #%v: %s = ", i+1, quiz.Question)
		answerCh := make(chan string)

		go func() {
			// Show question to user
			var response string
			fmt.Scan(&response)
			answerCh <- response
		}()

		select {
		case <-Timer.C:
			fmt.Println()
			break QuizLoop
		case response := <-answerCh:
			// Comparing user answer with Right answer (2nd column CSV)
			if response == quiz.Answer {
				QtyRightAnswer++
			}
		}

	}
	// Show the score
	fmt.Println("You scored", QtyRightAnswer, "out of", len(records), ".")
}

// Function to get the CSV file
func readData(fileName string) ([][]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}
