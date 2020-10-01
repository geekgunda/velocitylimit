package main

import (
	"fmt"
	"log"
	"time"

	"github.com/geekgunda/velocitylimit"
)

func main() {
	db := velocitylimit.NewStore()
	// inputFile from where we read transactions
	inputFile := "../../input.txt"
	// outputFile where we record the responses
	outputFile := fmt.Sprintf("../../output-%d.txt", time.Now().Unix())
	// driver fn actually triggering the whole process
	if err := velocitylimit.Process(inputFile, outputFile, db); err != nil {
		log.Println("Error processing: ", err)
	}
}
