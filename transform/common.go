package transform

import (
	"github.com/Xuanwo/beancollect/types"
)

// Available transform type.
const (
	TypeAddAccounts = "add_accounts"
)

// Transformer is the interface to do transform.
type Transformer interface {
	Transform(t *types.Transactions)
}

// Execute will execute transformers from rule.
func Execute(rule types.Rule, t *types.Transactions) {
	switch rule.Type {
	case TypeAddAccounts:
		tr := &AddAccounts{
			Condition: rule.Condition,
			Value:     rule.Value,
		}
		tr.Transform(t)
	default:
		return
	}
}
