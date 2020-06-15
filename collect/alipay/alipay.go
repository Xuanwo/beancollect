package alipay

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/Xuanwo/beancollect/types"
)

// Type is the type of wechat.
const Type = "alipay"

var (
	startComment = []byte("---------------------------------交易记录明细列表------------------------------------")
	endComment   = []byte("------------------------------------------------------------------------------------")
)

// AliPay is the struct for alipay type.
type AliPay struct {
}

type record struct {
	ID             string    // 交易号
	PayeeID        string    // 商家订单号
	CreatedAt      time.Time // 交易创建时间
	PayedAt        time.Time // 付款时间
	UpdatedAt      time.Time // 最近修改时间
	Source         string    // 交易来源地
	Type           string    // 类型
	Payee          string    // 交易对方
	Commodity      string    // 商品名称
	Amount         float64   // 金额（元）
	Flow           string    // 收/支
	Status         string    // 交易状态
	Fees           float64   // 服务费（元）
	Refund         float64   // 成功退款（元）
	Comment        string    // 备注
	CurrencyStatus string    // 资金状态
}

// NewAliPay will create a new alipay.
func NewAliPay() *AliPay {
	return &AliPay{}
}

// Parse implement Collector.Parse.
func (ali *AliPay) Parse(c *types.Config, r io.Reader) (t types.Transactions, err error) {
	t = make(types.Transactions, 0)

	rb, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("ioutil read failed for %s", err)
		return
	}

	b := make([]byte, 2*len(rb))

	_, _, err = simplifiedchinese.GB18030.NewDecoder().Transform(b, rb, true)
	if err != nil {
		log.Errorf("GB18030 read failed for %s", err)
		return
	}

	startIdx := bytes.Index(b, startComment)
	if startIdx != -1 {
		b = b[startIdx+len(startComment):]
	}
	endIdx := bytes.Index(b, endComment)
	if endIdx != -1 {
		b = b[:endIdx]
	}
	cr := csv.NewReader(bytes.NewReader(b))
	records, err := cr.ReadAll()
	if err != nil {
		log.Errorf("csv read failed for %s", err)
		return
	}

	for _, s := range records[2:] {
		r := &record{}

		idx := 0

		s[idx] = strings.TrimSpace(s[idx])
		r.ID = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.PayeeID = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		if s[idx] != "" {
			r.CreatedAt, err = time.Parse("2006-01-02 15:04:05", s[idx])
			if err != nil {
				log.Errorf("time <%s> parse failed for [%v]", s[idx], err)
				return
			}
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		if s[idx] != "" {
			r.PayedAt, err = time.Parse("2006-01-02 15:04:05", s[idx])
			if err != nil {
				log.Errorf("time <%s> parse failed for [%v]", s[idx], err)
				return
			}
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		if s[idx] != "" {
			r.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", s[idx])
			if err != nil {
				log.Errorf("time <%s> parse failed for [%v]", s[idx], err)
				return
			}
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Source = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Type = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Payee = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Commodity = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Amount, err = strconv.ParseFloat(s[idx], 64)
		if err != nil {
			log.Errorf("amount <%s> parse failed for [%v]", s[idx], err)
			return
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Flow = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Status = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Fees, err = strconv.ParseFloat(s[idx], 64)
		if err != nil {
			log.Errorf("fees <%s> parse failed for [%v]", s[idx], err)
			return
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Refund, err = strconv.ParseFloat(s[idx], 64)
		if err != nil {
			log.Errorf("refund <%s> parse failed for [%v]", s[idx], err)
			return
		}
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.Comment = s[idx]
		idx++

		s[idx] = strings.TrimSpace(s[idx])
		r.CurrencyStatus = s[idx]
		idx++

		t = append(t, formatTransaction(r, c))
	}

	return t, nil
}

func formatTransaction(r *record, c *types.Config) types.Transaction {
	t := types.Transaction{}

	t.Time = r.PayedAt
	t.Flag = "!"
	t.Narration = r.Commodity
	t.Accounts = make([]string, 1)
	t.Payee = r.Payee
	t.Amount = r.Amount
	if r.Flow == "支出" {
		t.Amount = -r.Amount
	}
	t.Currency = "CNY"

	return t
}
