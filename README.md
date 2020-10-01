
## Problem

In finance, it's common for accounts to have so-called "velocity limits". In this task, you'll write a program that accepts or declines attempts to load funds into customers' accounts in real-time.

Each attempt to load funds will come as a single-line JSON payload, structured as follows:

```json
{
  "id": "1234",
  "customer_id": "1234",
  "load_amount": "$123.45",
  "time": "2018-01-01T00:00:00Z"
}
```

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

As such, a user attempting to load $3,000 twice in one day would be declined on the second attempt, as would a user attempting to load $400 four times in a day.

For each load attempt, you should return a JSON response indicating whether the fund load was accepted based on the user's activity, with the structure:

```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```

You can assume that the input arrives in ascending chronological order and that if a load ID is observed more than once for a particular user, all but the first instance can be ignored. Each day is considered to end at midnight UTC, and weeks start on Monday (i.e. one second after 23:59:59 on Sunday).

Your program should process lines from `input.txt` and return output in the format specified above, either to standard output or a file. Expected output given our input data can be found in `output.txt`.

You're welcome to write your program in a general-purpose language of your choosing, but as we use Go on the back-end and TypeScript on the front-end, we do have a preference towards solutions written in Go (back-end) and TypeScript (front-end).

We value well-structured, self-documenting code with sensible test coverage. Descriptive function and variable names are appreciated, as is isolating your business logic from the rest of your code.

## Solution
This app solves the problem statement above.  

Usage: `make all`

Assumption: 
- Go is installed and setup already
- Code is checked out, inside `$GOPATH`, with proper directory structure (`$GOPATH/src/github.com/geekgunda`)

Configuration:  
By default, the app takes `input.txt` file to read the transactions.  
Output is recorded in a time-stamp appended `output-<timestamp>.txt` file in the same folder.  
These values can be changed inside `cmd/app/main.go`


### Design:  

`velocitylimit` package: Holds the business logic.  
`cmd/app/` package: Holds the main.go file, that triggers the app.

`Process()` function acts as a driver.  
It reads the input file line by line and records the response in output file in the same order.  
Responses for duplicate transactions are omitted (as they aren't processed)

The state is maintained in memory via a map with mutex lock.  
For the input size of 1000 transactions, an external DB wasn't necessary.  
But the app is written with a model (`store.go`), such that adding a DB should have minimal changes.  

Due to the restriction on output formatting, it was difficult to come up with a concurrent solution.  
Possible ways to do that (assuming the output restrictions aren't in place):  
1. Change this into a API service, that can be called concurrently.
2. Update Process function to spawn a new go channel per account, and a go routine listening on this channel. This way, requests for an account can be serialized, while concurrently processing requests at an app or service level


