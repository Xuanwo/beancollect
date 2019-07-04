package wechat

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Xuanwo/beancollect/types"
)

// Type is the type of wechat.
const Type = "wechat"

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

	headerParsed := false
	buf := bufio.NewScanner(r)
	for buf.Scan() {
		line := buf.Text()

		// WeChat will use ",,,,,,,," suffix for comment.
		if strings.HasSuffix(line, ",,,,,,,,") {
			continue
		}

		// The first line is csv header, we should also ignore them.
		if !headerParsed {
			headerParsed = true
			continue
		}

		s := strings.Split(line, ",")

		r := &record{}
		r.Time, err = time.Parse("2006-01-02 15:04:05", s[0])
		if err != nil {
			log.Errorf("time parse failed for %v", err)
			return
		}
		r.Type = strings.TrimSpace(s[1])
		r.Payee = strings.TrimSpace(s[2])
		r.Commodity = strings.TrimSpace(s[3])
		r.Flow = strings.TrimSpace(s[4])
		// WeChat will use "¥298.00" for amount, we should trim it.
		r.Amount, err = strconv.ParseFloat(s[5][2:], 64)
		if err != nil {
			log.Errorf("amount parse failed for %v", err)
			return
		}
		r.Payment = strings.TrimSpace(s[6])
		r.Status = strings.TrimSpace(s[7])
		r.ID = strings.TrimSpace(s[8])
		r.PayeeID = strings.TrimSpace(s[9])
		r.Comment = strings.TrimSpace(s[10])

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
	t.Payee = r.Payee
	t.Accounts = append(t.Accounts, c.Account[r.Payment])
	t.Amount = r.Amount
	if r.Flow == "支出" {
		t.Amount = -r.Amount
	}
	t.Currency = "CNY"
	return t
}
