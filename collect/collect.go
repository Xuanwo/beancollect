package collect

import (
	"io"
	"log"

	"github.com/Xuanwo/beancollect/collect/alipay"
	"github.com/Xuanwo/beancollect/collect/wechat"
	"github.com/Xuanwo/beancollect/types"
)

// Collector is the interface to parse transactions from io.Reader.
type Collector interface {
	Parse(c *types.Config, r io.Reader) (types.Transactions, error)
}

// NewCollector will create a new collector.
func NewCollector(t string) Collector {
	switch t {
	case wechat.Type:
		return wechat.NewWeChat()
	case alipay.Type:
		return alipay.NewAliPay()
	default:
		log.Fatalf("not supported type %s", t)
	}

	return nil
}
