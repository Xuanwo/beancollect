package types

import "time"

// Transaction means a real transaction in beancount.
type Transaction struct {
	Time      time.Time
	Flag      string
	Narration string
	Payee     string
	Tags      []string
	Accounts  []string
	Amount    float64
	Currency  string

	// Metadata map[string]string
}

// IsMatch will check whether this transaction match condition.
// TODO: support regex
func (t Transaction) IsMatch(cond map[string]string) bool {
	// Currently, we only support payee.
	if cond["payee"] == t.Payee {
		return true
	}
	return false
}

// Transactions is the array for transactions
type Transactions []Transaction

// Len implement Sorter.Len
func (t Transactions) Len() int {
	return len(t)
}

// Less implement Sorter.Less
func (t Transactions) Less(i, j int) bool {
	return t[i].Time.Before(t[j].Time)
}

// Swap implement Sorter.Swap
func (t Transactions) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
