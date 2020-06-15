package cmbchina

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Xuanwo/beancollect/types"
	log "github.com/sirupsen/logrus"
)

const Type = "cmbchina"

type CMBChina struct {
}

type record struct {
	TransDate           string  // 交易日
	PostDate            string  // 记账日
	Description         string  // 交易摘要
	RMBAmount           float64 // 人民币金额
	CardNumber          string  // 卡号后四位
	Area                string  // 交易地点
	OriginalTransAmount float64 // 交易地金额
}

func NewCMBChina() *CMBChina {
	return &CMBChina{}
}

func (cmb *CMBChina) Parse(c *types.Config, r io.Reader) (t types.Transactions, err error) {
	t = make(types.Transactions, 0)

	rb, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("ioutil read failed for %s", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rb))
	if err != nil {
		log.Errorf("parse html failed for %s", err)
		return
	}

	var rs []*record
	var tr *record
	doc.Find("span[id$=\"fixBand15\"]").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			if _, ok := s.Attr("valign"); !ok {
				return
			}

			value := s.Text()
			switch (i - 2) % 7 {
			case 0:
				tr = &record{}
				rs = append(rs, tr)

				tr.TransDate = value
			case 1:
				tr.PostDate = value
			case 2:
				tr.Description = value
			case 3:
				value = strings.TrimPrefix(value, "￥")
				value = strings.TrimSpace(value)
				// CMBChina will insert "," in amount
				value = strings.ReplaceAll(value, ",", "")
				// CMBChina will have "\u00a0" between "-" and amount
				value = strings.ReplaceAll(value, "\u00a0", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Errorf("parse rmb amount failed: %s", err)
					return
				}
				tr.RMBAmount = amount
			case 4:
				tr.CardNumber = value
			case 5:
				tr.Area = value
			case 6:
				// CMBChina will insert "," in amount
				value = strings.ReplaceAll(value, ",", "")
				amount, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Errorf("parse original trans amount failed: %s", err)
					return
				}
				tr.OriginalTransAmount = amount
			}
		})
	})

	for _, v := range rs {
		t = append(t, formatTransaction(v, c))
	}
	return
}

func formatTransaction(r *record, c *types.Config) types.Transaction {
	t := types.Transaction{}

	var err error
	if len(r.PostDate) != 0 {
		t.Time, err = time.Parse("0102", r.PostDate)
		if err != nil {
			log.Errorf("parse time failed for %s", err)
		}
	}
	t.Flag = "!"
	t.Accounts = append(t.Accounts, c.Account[r.CardNumber])
	t.Payee = r.Description
	t.Amount = r.RMBAmount
	t.Currency = "CNY"

	return t
}
