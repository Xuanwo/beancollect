package wechat

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Xuanwo/beancollect/types"
)

// Type is the type of wechat.
const Type = "wechat"

// Available status for wechat.
const (
	StatusPaymentSuccess  = "支付成功"
	StatusWithdrawSuccess = "提现已到账"
	StatusDepositSuccess  = "已存入零钱"
	StatusRefundSuccess   = "已全额退款"
)

var comments = []byte(",,,,,,,,")

// WeChat is the struct for wechat type.
type WeChat struct {
}

// NewWeChat will create a new wechat.
func NewWeChat() *WeChat {
	return &WeChat{}
}

type record struct {
	Time      time.Time // 交易时间
	Type      string    // 交易类型
	Payee     string    // 交易对方
	Commodity string    // 商品
	Flow      string    // 收/支
	Amount    float64   // 金额(元)
	Payment   string    // 支付方式
	Status    string    // 当前状态
	ID        string    // 交易单号
	PayeeID   string    // 商户单号
	Comment   string    // 备注
}

// Parse implement Collector.Parse.
func (w *WeChat) Parse(c *types.Config, r io.Reader) (t types.Transactions, err error) {
	t = make(types.Transactions, 0)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("ioutil read failed for %s", err)
		return
	}

	idx := bytes.LastIndex(b, comments)
	if idx != -1 {
		b = b[idx+len(comments):]
	}

	cr := csv.NewReader(bytes.NewReader(b))
	records, err := cr.ReadAll()
	if err != nil {
		log.Errorf("csv read failed for %s", err)
		return
	}
	for _, s := range records[1:] {
		r := &record{}
		r.Time, err = time.Parse("2006-01-02 15:04:05", s[0])
		if err != nil {
			log.Errorf("time <%s> parse failed for [%v]", s[0], err)
			return
		}
		r.Type = strings.TrimSpace(s[1])
		r.Payee = strings.TrimSpace(s[2])
		r.Commodity = strings.TrimSpace(s[3])
		r.Flow = strings.TrimSpace(s[4])
		// WeChat will use "¥298.00" for amount, we should trim it.
		r.Amount, err = strconv.ParseFloat(s[5][2:], 64)
		if err != nil {
			log.Errorf("amount <%s> parse failed for [%v]", s[5], err)
			return
		}
		r.Payment = strings.TrimSpace(s[6])
		r.Status = strings.TrimSpace(s[7])
		r.ID = strings.TrimSpace(s[8])
		r.PayeeID = strings.TrimSpace(s[9])
		r.Comment = strings.TrimSpace(s[10])

		// Ignore all refund payment.
		if r.Status == StatusRefundSuccess {
			continue
		}

		t = append(t, formatTransaction(r, c))
	}

	return t, nil
}

// formatTransaction will format record into transaction.
func formatTransaction(r *record, c *types.Config) types.Transaction {
	t := types.Transaction{}

	t.Time = r.Time
	t.Flag = "!"
	// WeChat may have " around narration, let's trim them.
	t.Narration = strings.Trim(r.Commodity, "\"")
	if t.Narration == "/" {
		t.Narration = strings.Trim(r.Type, "\"")
	}
	t.Payee = r.Payee

	if _, ok := c.Account[r.Payment]; !ok {
		log.Infof("payment %s doesn't have related account", r.Payment)
	}
	t.Accounts = append(t.Accounts, c.Account[r.Payment])
	t.Amount = r.Amount
	if r.Flow == "支出" {
		t.Amount = -r.Amount
	}
	t.Currency = "CNY"

	if r.Status == StatusWithdrawSuccess {
		t.Accounts = append(t.Accounts, c.Account["零钱"])
	}

	return t
}
