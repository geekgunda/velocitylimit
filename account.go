package velocitylimit

import (
	"log"
	"time"
)

// Account holds details about a single customer
type Account struct {
	CustID      string
	Balance     float64
	DailyLimit  *DailyLimit
	WeeklyLimit *WeeklyLimit
}

// DailyLimit holds velocity limits applied to an account
// with in a window of a day
type DailyLimit struct {
	Date        time.Time
	DailyBudget float64
	DailyTxn    int
}

// WeeklyLimit holds velocity limits applied to an account
// within a window of a week
type WeeklyLimit struct {
	Date         time.Time
	WeeklyBudget float64
}

// NewDailyLimit initializes a new copy of DailyLimit with default limits
func NewDailyLimit(d time.Time) *DailyLimit {
	return &DailyLimit{
		Date:        getStartOfDay(d),
		DailyBudget: 5000,
		DailyTxn:    3,
	}
}

// NewWeeklyLimit initializes a copy of WeeklyLimit with default values
func NewWeeklyLimit(d time.Time) *WeeklyLimit {
	return &WeeklyLimit{
		Date:         getStartOfWeek(d),
		WeeklyBudget: 20000,
	}
}

// NewAccount initializes a copy of Account for a given customer_id
func NewAccount(custID string) *Account {
	return &Account{
		CustID: custID,
	}
}

// Validate checks whether given transaction violates existing DailyLimits
func (dl *DailyLimit) Validate(amount float64, t time.Time) bool {
	transDay := getStartOfDay(t)
	if transDay.After(dl.Date) {
		log.Println("Transaction date: ", transDay, " is newer than current daily date: ", dl.Date)
		newdl := NewDailyLimit(t)
		*dl = *newdl
	}
	if dl.DailyBudget-amount < 0 {
		return false
	}
	if dl.DailyTxn-1 < 0 {
		return false
	}
	return true
}

// Apply updates DailyLimit quota as a result of accepting a transaction
func (dl *DailyLimit) Apply(amount float64) {
	dl.DailyBudget -= amount
	dl.DailyTxn--
}

// Validate verifies whether a given transaction violates WeeklyLimits
func (wl *WeeklyLimit) Validate(amount float64, t time.Time) bool {
	transWeek := getStartOfWeek(t)
	if transWeek.After(wl.Date) {
		log.Println("Transaction date: ", transWeek, " is newer than current week start: ", wl.Date)
		newwl := NewWeeklyLimit(t)
		*wl = *newwl
	}
	if wl.WeeklyBudget-amount < 0 {
		return false
	}
	return true
}

// Apply updates weekly limit after a successful transaction
func (wl *WeeklyLimit) Apply(amount float64) {
	wl.WeeklyBudget -= amount
}

// LoadFunds handles request to add funcs to an account
// It first verifies whether velocity limits allow this transaction
// If they do, the fields are updated accordingly
func (a *Account) LoadFunds(r *Request) bool {
	// Validate if daily limits are satisfied
	if a.DailyLimit == nil {
		log.Println("Daily limit is missing. Adding now")
		a.DailyLimit = NewDailyLimit(r.ParsedTime)
	}
	if a.DailyLimit.Validate(r.ParsedAmount, r.ParsedTime) == false {
		log.Println("Daily limit breached. Transaction rejected: ", r.ID)
		return false
	}
	// Validate if weekly limits are satisfied
	if a.WeeklyLimit == nil {
		log.Println("Weekly limit is missing. Adding now")
		a.WeeklyLimit = NewWeeklyLimit(r.ParsedTime)
	}
	if a.WeeklyLimit.Validate(r.ParsedAmount, r.ParsedTime) == false {
		log.Println("Weekly limit breached. Transaction rejected: ", r.ID)
		return false
	}
	// Load the funds
	a.Balance += r.ParsedAmount
	// Update the limits after acting on this transaction
	a.DailyLimit.Apply(r.ParsedAmount)
	a.WeeklyLimit.Apply(r.ParsedAmount)
	log.Println("Transaction approved: ", r.ID)
	return true
}

// getStartOfDay returns the time value at the beginning of a day with UTC TZ
func getStartOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

// getStartOfWeek returns the time value at the beginning of the Monday within the given week with UTC TZ
func getStartOfWeek(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+int(time.Monday-d.Weekday()), 0, 0, 0, 0, time.UTC)
}
