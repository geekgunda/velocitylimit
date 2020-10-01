package velocitylimit

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// Process is the driver function for the given task
// It takes an inputFile and outputFile as arguments along with a reference to DB
// It reads data from inputFile, processes it and records response in outputFile.
// DB instance is used to fetch and record details about the account and transaction
func Process(inputFile string, outputFile string, db *Store) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	for scanner.Scan() {
		// Parse the request
		req, err := NewRequest(scanner.Text())
		if err != nil {
			return err
		}
		// Check if it's already processed
		if db.IsDupTxn(req.ID, req.CustID) {
			log.Println("Ignoring duplicate txn: ", req.ID)
			continue
		}
		// Fetch the account from DB
		account := db.GetAccount(req.CustID)
		// Act on the request (if velocity limits agree)
		accepted := account.LoadFunds(req)
		response := NewResponse(req.ID, req.CustID, accepted)
		resBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		// Record the response
		if _, err = writer.WriteString(string(resBytes) + "\n"); err != nil {
			return err
		}
		// Record the transaction in DB
		db.AddTxn(req.ID, req.CustID)
	}
	// Check if there were any errors while reading the input file
	if err := scanner.Err(); err != nil {
		return err
	}
	// Flush any pending writes
	writer.Flush()
	return nil
}
