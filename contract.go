package velocitylimit

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

// Request holds a request to process addFunds request
type Request struct {
	ID           string    `json:"id"`
	CustID       string    `json:"customer_id"`
	Amount       string    `json:"load_amount"`
	Time         string    `json:"time"`
	ParsedAmount float64   `json:"-"`
	ParsedTime   time.Time `json:"-"`
}

// Response defines the response structure for addFunds request
type Response struct {
	ID       string `json:"id"`
	CustID   string `json:"customer_id"`
	Accepted bool   `json:"accepted"`
}

// NewRequest returns a parsed and populated copy of Request for a given string input
func NewRequest(reqStr string) (*Request, error) {
	var r Request
	if err := json.Unmarshal([]byte(reqStr), &r); err != nil {
		log.Println("Error parsing line: ", err)
		return nil, err
	}
	if a, err := strconv.ParseFloat(strings.Trim(r.Amount, "$"), 64); err != nil {
		log.Println("Error parsing amount: ", err)
		return nil, err
	} else {
		r.ParsedAmount = a
	}
	//if t, err := time.Parse("2006-01-02T22:04:05Z", r.Time); err != nil {
	if t, err := time.Parse(time.RFC3339, r.Time); err != nil {
		log.Println("Error parsing time: ", err)
		return nil, err
	} else {
		r.ParsedTime = t
	}
	return &r, nil
}

// NewResponse populates and returns a copy of Response with given inputs
func NewResponse(id string, custID string, accepted bool) *Response {
	return &Response{
		ID:       id,
		CustID:   custID,
		Accepted: accepted,
	}
}
