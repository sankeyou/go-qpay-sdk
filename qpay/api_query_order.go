package qpay

import (
	"encoding/xml"
	"time"
)

const QUERY_ORDER_URL = "https://qpay.qq.com/cgi-bin/pay/qpay_order_query.cgi"

type QueryOrderResult struct {
	BaseResult
	DeviceInfo    string `xml:"device_info"`
	TradeType     string `xml:"trade_type"`
	TradeState    string `xml:"trade_state"`
	BankType      string `xml:"bank_type"`
	FeeType       string `xml:"fee_type"`
	TotalFee      uint64 `xml:"total_fee"`
	CashFee       uint64 `xml:"cash_fee"`
	CouponFee     uint64 `xml:"coupon_fee"`
	TransactionId string `xml:"transaction_id"`
	TradeNumber   string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEndStr    string `xml:"time_end"`
	TimeEnd       time.Time
	OpenId        string `xml:"openid"`
}

func (ctx *Context) QueryOrder(transactionId string, tradeNumber string) (*QueryOrderResult, error) {
	m := StringMap{
		"mch_id":    ctx.merchantId,
		"nonce_str": GeneratorNonce(),
	}
	if len(ctx.appId) > 0 {
		m["appid"] = ctx.appId
	}
	if len(transactionId) > 0 {
		m["transaction_id"] = transactionId
	}
	if len(tradeNumber) > 0 {
		m["out_trade_no"] = tradeNumber
	}

	m["sign"] = ctx.sign(m)
	bs, err := xml.MarshalIndent(m, "", "  ")
	if nil != err {
		return nil, err
	}
	resp, err := ctx.HttpClient.R().SetBody(bs).Post(QUERY_ORDER_URL)
	if nil != err {
		return nil, err
	}

	out := &QueryOrderResult{}
	if err = xml.Unmarshal(resp.Body(), out); nil != err {
		return nil, err
	}
	out.BaseResult.process()
	out.TimeEnd, _ = time.Parse("20060102150405", out.TimeEndStr)
	return out, nil
}
