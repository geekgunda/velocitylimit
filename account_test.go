package velocitylimit

import (
	"testing"
	"time"
)

func TestSingleTxnDailyLimit(t *testing.T) {
	var cases = []struct {
		name    string
		amount  float64
		date    time.Time
		expResp bool
	}{
		{"tc1", 5000, time.Now(), true},
		{"tc2", 5001, time.Now(), false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			dl := NewDailyLimit(time.Now())
			resp := dl.Validate(tt.amount, tt.date)
			if resp != tt.expResp {
				t.Errorf("Test failed. Resp: %v | Case: %v", resp, tt)
			}
		})
	}
}

func TestMultipleTxnDailyLimit(t *testing.T) {
	var cases = []struct {
		name    string
		amount  float64
		date    time.Time
		expResp bool
	}{
		{"tc1", 5000, time.Now(), true},
		{"tc2", 10, time.Now(), false},
	}
	dl := NewDailyLimit(time.Now().AddDate(0, 0, -1))
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := dl.Validate(tt.amount, tt.date)
			if resp != tt.expResp {
				t.Errorf("Test failed. Resp: %v | Case: %v", resp, tt.name)
			}
			if resp {
				dl.Apply(tt.amount)
			}
		})
	}
}

func TestDailyTxnLimit(t *testing.T) {
	var cases = []struct {
		name    string
		amount  float64
		date    time.Time
		expResp bool
	}{
		{"tc1", 2000, time.Now(), true},
		{"tc2", 1000, time.Now(), true},
		{"tc3", 500, time.Now(), true},
		{"tc4", 500, time.Now(), false},
	}
	dl := NewDailyLimit(time.Now().AddDate(0, 0, -1))
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := dl.Validate(tt.amount, tt.date)
			if resp != tt.expResp {
				t.Errorf("Test failed. Resp: %v | Case: %v", resp, tt.name)
			}
			if resp {
				dl.Apply(tt.amount)
			}
		})
	}
}

func TestSingleTxnWeeklyLimit(t *testing.T) {
	var cases = []struct {
		name    string
		amount  float64
		date    time.Time
		expResp bool
	}{
		{"tc1", 20000, time.Now(), true},
		{"tc2", 20001, time.Now(), false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			wl := NewWeeklyLimit(time.Now())
			resp := wl.Validate(tt.amount, tt.date)
			if resp != tt.expResp {
				t.Errorf("Test failed. Resp: %v | Case: %v", resp, tt)
			}
		})
	}
}

func TestMultipleTxnWeeklyLimit(t *testing.T) {
	var cases = []struct {
		name    string
		amount  float64
		date    time.Time
		expResp bool
	}{
		{"tc1", 5000, time.Now(), true},
		{"tc2", 1000, time.Now(), true},
		{"tc3", 5000, time.Now(), true},
		{"tc4", 4000, time.Now(), true},
		{"tc5", 5000, time.Now(), true},
		{"tc6", 100, time.Now(), false},
	}
	wl := NewWeeklyLimit(time.Now().AddDate(0, 0, -1))
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp := wl.Validate(tt.amount, tt.date)
			if resp != tt.expResp {
				t.Errorf("Test failed. Resp: %v | Case: %v", resp, tt.name)
			}
			if resp {
				wl.Apply(tt.amount)
			}
		})
	}
}
