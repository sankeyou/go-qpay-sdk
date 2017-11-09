package qpay

import (
	"encoding/xml"
	"strconv"
	"time"
)

type NotifyOrderArgs struct {
	MerchantId    string
	Nonce         string
	DeviceInfo    string
	TradeType     string
	BankType      string
	FeeType       string
	TotalFee      uint64
	CashFee       uint64
	CouponFee     uint64
	TransactionId string
	TradeNumber   string
	TimeEnd       time.Time
	OpenId        string

	raw NotifyOrderArgsRaw
}

type NotifyOrderArgsRaw struct {
	AppId         string `xml:"appid"`
	MerchantId    string `xml:"mch_id"`
	Nonce         string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	DeviceInfo    string `xml:"device_info"`
	TradeType     string `xml:"trade_type"`
	TradeState    string `xml:"trade_state"`
	BankType      string `xml:"bank_type"`
	FeeType       string `xml:"fee_type"`
	TotalFee      string `xml:"total_fee"`
	CashFee       string `xml:"cash_fee"`
	CouponFee     string `xml:"coupon_fee"`
	TransactionId string `xml:"transaction_id"`
	TradeNumber   string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEnd       string `xml:"time_end"`
	OpenId        string `xml:"openid"`
}

func NotifyOrderParseRequest(body []byte) (*NotifyOrderArgs, error) {
	args := &NotifyOrderArgs{}
	raw := &args.raw
	if err := xml.Unmarshal(body, raw); nil != err {
		return nil, err
	}

	args.MerchantId = raw.MerchantId
	args.Nonce = raw.Nonce
	args.DeviceInfo = raw.DeviceInfo
	args.TradeType = raw.TradeType
	args.BankType = raw.BankType
	args.FeeType = raw.FeeType
	args.TotalFee, _ = strconv.ParseUint(raw.TotalFee, 10, 64)
	args.CashFee, _ = strconv.ParseUint(raw.CashFee, 10, 64)
	args.CouponFee, _ = strconv.ParseUint(raw.CouponFee, 10, 64)
	args.TransactionId = raw.TransactionId
	args.TradeNumber = raw.TradeNumber
	args.TimeEnd, _ = time.Parse("20060102150405", raw.TimeEnd)
	args.OpenId = raw.OpenId

	return args, nil
}

func (ctx *Context) NotifyOrderCheckSign(args *NotifyOrderArgs) bool {
	raw := &args.raw
	m := StringMap{
		"mch_id":         raw.MerchantId,
		"nonce_str":      raw.Nonce,
		"trade_type":     raw.TradeType,
		"trade_state":    raw.TradeState,
		"bank_type":      raw.BankType,
		"fee_type":       raw.FeeType,
		"total_fee":      raw.TotalFee,
		"cash_fee":       raw.CashFee,
		"transaction_id": raw.TransactionId,
		"out_trade_no":   raw.TradeNumber,
		"time_end":       raw.TimeEnd,
	}
	if len(raw.AppId) > 0 {
		m["appid"] = raw.AppId
	}
	if len(raw.DeviceInfo) > 0 {
		m["device_info"] = raw.DeviceInfo
	}
	if len(raw.CouponFee) > 0 {
		m["coupon_fee"] = raw.CouponFee
	}
	if len(raw.Attach) > 0 {
		m["attach"] = raw.Attach
	}
	if len(raw.OpenId) > 0 {
		m["openid"] = raw.OpenId
	}
	return raw.Sign == ctx.sign(m)
}

func NotifyOrdeGenerateResponse(success bool, message string) ([]byte, error) {
	m := StringMap{}
	if success {
		m["return_code"] = "SUCCESS"
	} else {
		m["return_code"] = "FAIL"
		m["return_msg"] = message
	}
	return xml.MarshalIndent(m, "", "  ")
}
