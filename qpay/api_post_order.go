package qpay

import (
	"encoding/xml"
	"strconv"
	"time"
)

const POST_ORDER_URL = "https://qpay.qq.com/cgi-bin/pay/qpay_unified_order.cgi"

const (
	_ = iota
	PAY_LIMITATION_NO_BALANCE
	PAY_LIMITATION_NO_CREDIT
	PAY_LIMITATION_NO_DEBIT
	PAY_LIMITATION_BALANCE_ONLY
	PAY_LIMITATION_DEBIT_ONLY
	PAY_LIMITATION_NO_BIND_NO_BALANCE
)

var payLimitations map[int]string = map[int]string{
	PAY_LIMITATION_NO_BALANCE:         "no_balance",
	PAY_LIMITATION_NO_CREDIT:          "no_credit",
	PAY_LIMITATION_NO_DEBIT:           "no_debit",
	PAY_LIMITATION_BALANCE_ONLY:       "balance_only",
	PAY_LIMITATION_DEBIT_ONLY:         "debit_only",
	PAY_LIMITATION_NO_BIND_NO_BALANCE: "NoBindNoBalan",
}

const (
	_ = iota
	TRADE_TYPE_MICROAPP
	TRADE_TYPE_APP
	TRADE_TYPE_JSAPI
	TRADE_TYPE_NATIVE
)

var tradeTypes map[int]string = map[int]string{
	TRADE_TYPE_MICROAPP: "MICROAPP",
	TRADE_TYPE_APP:      "APP",
	TRADE_TYPE_JSAPI:    "JSAPI",
	TRADE_TYPE_NATIVE:   "NATIVE",
}

type PostOrderArgs struct {
	Body              string
	Attach            string
	FeeType           string
	Fee               uint64
	Ip                string
	TimeStart         *time.Time
	TimeExpire        *time.Time
	PayLimitation     int
	ContractCode      string
	PromotionSaleTag  string
	PromotionLevelTag string
	TradeType         int
	NotifyUrl         string
	DeviceInfo        string
}

type PostOrderResult struct {
	BaseResult
	Nonce       string
	TradeNumber string
	PrePayId    string `xml:"prepay_id"`
	QrcodeUrl   string `xml:"code_url"`
}

func (ctx *Context) PostOrder(args *PostOrderArgs) (*PostOrderResult, error) {
	nonce, tradeNumber := GeneratorNonce(), ctx.tradeNumberGenerator()
	m := StringMap{
		"mch_id":           ctx.merchantId,
		"nonce_str":        nonce,
		"body":             args.Body,
		"out_trade_no":     tradeNumber,
		"fee_type":         "CNY",
		"total_fee":        strconv.FormatUint(args.Fee, 10),
		"spbill_create_ip": args.Ip,
		"trade_type":       tradeTypes[TRADE_TYPE_NATIVE],
		"notify_url":       args.NotifyUrl,
	}

	if len(ctx.appId) > 0 {
		m["appid"] = ctx.appId
	}
	if len(args.Attach) > 0 {
		m["attach"] = args.Attach
	}
	if len(args.FeeType) > 0 {
		m["fee_type"] = args.FeeType
	}
	if nil != args.TimeStart {
		m["time_start"] = args.TimeStart.Format("20060102150405")
	}
	if nil != args.TimeExpire {
		m["time_expire"] = args.TimeExpire.Format("20060102150405")
	}
	if payLimitationStr, ok := payLimitations[args.PayLimitation]; ok {
		m["limit_pay"] = payLimitationStr
	}
	if len(args.ContractCode) > 0 {
		m["contract_code"] = args.ContractCode
	}
	if len(args.PromotionSaleTag) > 0 {
		s := "promotion_tag=" + args.PromotionSaleTag
		if len(args.PromotionLevelTag) > 0 {
			s += "&sale_tag=" + args.PromotionLevelTag
		}
		m["promotion_tag"] = s
	}
	if tradeTypeStr, ok := tradeTypes[args.TradeType]; ok {
		m["trade_type"] = tradeTypeStr
	}
	if len(args.DeviceInfo) > 0 {
		m["device_info"] = args.DeviceInfo
	}

	m["sign"] = ctx.sign(m)
	bs, err := xml.MarshalIndent(m, "", "  ")
	if nil != err {
		return nil, err
	}
	resp, err := ctx.HttpClient.R().SetBody(bs).Post(POST_ORDER_URL)
	if nil != err {
		return nil, err
	}

	out := &PostOrderResult{}
	if err = xml.Unmarshal(resp.Body(), out); nil != err {
		return nil, err
	}
	out.BaseResult.process()
	out.Nonce = nonce
	out.TradeNumber = tradeNumber
	return out, nil
}
