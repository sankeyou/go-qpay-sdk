[![Go Report Card](https://goreportcard.com/badge/github.com/ldeng7/go-qpay-sdk)](https://goreportcard.com/report/github.com/ldeng7/go-qpay-sdk)

Synopsis
========

```go
package main

import (
	"fmt"
	"time"

	"github.com/ldeng7/go-qpay-sdk/qpay"
)

// Should put it in a context struct
var ctx *qpay.Context

func init() {
	// Do it in process initialization
	// NOT in each time posting an order
	ctx = qpay.Init(&qpay.InitArgs{
		MerchantId:  "{merchantId}",
		ApiKey:      "{apiKey}",
		HttpTimeout: 5 * time.Second,
	})
}

func postOrder() {
	// https://qpay.qq.com/qpaywiki/showdocument.php?pid=38&docid=58
	now := time.Now()
	expire := now.Add(time.Hour)
	args := &qpay.PostOrderArgs{
		Body:              "test", // production description
		Attach:            "",     // attached data, optional
		FeeType:           "CNY",  // default "CNT", optional
		Fee:               8,      // in Chinese Fen
		Ip:                "171.221.200.59",
		TimeStart:         &now,                         // optional
		TimeExpire:        &expire,                      // optional
		PayLimitation:     qpay.PAY_LIMITATION_NO_DEBIT, // optional
		ContractCode:      "",                           // optional
		PromotionSaleTag:  "",                           // optional
		PromotionLevelTag: "",                           // optional
		TradeType:         qpay.TRADE_TYPE_NATIVE,       // default qpay.TRADE_TYPE_NATIVE, optional
		NotifyUrl:         "http://112.124.50.175:30778/tt",
		DeviceInfo:        "", // optional
	}

	res, _ := ctx.PostOrder(args)
	fmt.Printf("%v\n", res)
}

func parseNotification() {
	args, _ := qpay.NotifyOrderParseRequest([]byte(`<?xml version="1.0" encoding="UTF-8" ?><xml>...</xml>`))
	fmt.Printf("%v\n", args)
	println(ctx.NotifyOrderCheckSign(args))

	respBody, _ := qpay.NotifyOrdeGenerateResponse(true, "")
	println(string(respBody))
}

func queryOrder() {
	res, _ := ctx.QueryOrder("{transactionId}", "")
	fmt.Printf("%v\n", res)
}

func closeOrder() {
	res, _ := ctx.CloseOrder("{tradeNumber}", 8)
	fmt.Printf("%v\n", res)
}
```
