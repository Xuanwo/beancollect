package transform

import (
	"github.com/Xuanwo/beancollect/types"
)

// AddAccounts will do transform for add_accounts type.
type AddAccounts struct {
	Condition map[string]string
	Value     string
}

// Transform implements Transformer interface.
func (a *AddAccounts) Transform(t *types.Transactions) {
	for k, v := range *t {
		if !v.IsMatch(a.Condition) {
			continue
		}
		(*t)[k].Accounts = append((*t)[k].Accounts, a.Value)
	}
	return
}
