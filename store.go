package velocitylimit

import "sync"

// Store holds the current state of all active accounts and transactions
// This would be replaced by a database in production env
type Store struct {
	sync.RWMutex
	accounts map[string]*Account
	txns     map[string]bool
}

// NewStore initializes the internal in-mem storage
// This would be replaced with a request to connect to DB
func NewStore() *Store {
	return &Store{
		accounts: make(map[string]*Account),
		txns:     make(map[string]bool),
	}
}

// GetAccount returns current a copy of Account with current state
func (s *Store) GetAccount(custID string) *Account {
	acc := s.getAccountFromStore(custID)
	if acc != nil {
		return acc
	}
	return s.addAccountToStore(custID)
}

// getAccountFromStore returns an existing copy of account from cache (if available)
func (s *Store) getAccountFromStore(custID string) *Account {
	s.RLock()
	defer s.RUnlock()
	if acc, ok := s.accounts[custID]; ok {
		return acc
	}
	return nil
}

// addAccountToStore adds a new account for given customer ID in cache
func (s *Store) addAccountToStore(custID string) *Account {
	s.Lock()
	defer s.Unlock()
	acc := NewAccount(custID)
	s.accounts[custID] = acc
	return acc
}

// AddTxn adds a record for a given transaction ID against a customer ID
func (s *Store) AddTxn(id, custID string) {
	s.Lock()
	defer s.Unlock()
	s.txns[id+custID] = true
}

// IsDupTxn checks if the given transaction ID has already been processed for a customer ID
func (s *Store) IsDupTxn(id, custID string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.txns[id+custID]
	return ok
}
