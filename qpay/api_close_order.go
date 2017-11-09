package qpay

import (
	"encoding/xml"
	"strconv"
)

const CLOSE_ORDER_URL = "https://qpay.qq.com/cgi-bin/pay/qpay_close_order.cgi"

type CloseOrderResult struct {
	BaseResult
}

func (ctx *Context) CloseOrder(tradeNumber string, fee uint64) (*CloseOrderResult, error) {
	m := StringMap{
		"mch_id":       ctx.merchantId,
		"nonce_str":    GeneratorNonce(),
		"out_trade_no": tradeNumber,
		"total_fee":    strconv.FormatUint(fee, 10),
	}
	if len(ctx.appId) > 0 {
		m["appid"] = ctx.appId
	}

	m["sign"] = ctx.sign(m)
	bs, err := xml.MarshalIndent(m, "", "  ")
	if nil != err {
		return nil, err
	}
	resp, err := ctx.HttpClient.R().SetBody(bs).Post(CLOSE_ORDER_URL)
	if nil != err {
		return nil, err
	}

	out := &CloseOrderResult{}
	if err = xml.Unmarshal(resp.Body(), out); nil != err {
		return nil, err
	}
	out.BaseResult.process()
	return out, nil
}
