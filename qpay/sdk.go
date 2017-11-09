package qpay

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty"
)

type InitArgs struct {
	MerchantId string
	ApiKey     string
	AppId      string

	TradeNumberGenerator func() string
	Logger               *log.Logger
	HttpTimeout          time.Duration
}

type Context struct {
	merchantId string
	apiKey     string
	appId      string

	tradeNumberGenerator func() string
	logger               *log.Logger
	HttpClient           *resty.Client
}

func (ctx *Context) sign(args StringMap) string {
	nArgs := len(args)
	argKeys := make([]string, nArgs)
	argCons := make([]string, nArgs+1)
	i := 0
	for k, _ := range args {
		argKeys[i] = k
		i++
	}
	sort.Strings(argKeys)

	for i, k := range argKeys {
		argCons[i] = fmt.Sprintf("%s=%s", k, args[k])
	}
	argCons[nArgs] = "key=" + ctx.apiKey
	s := strings.Join(argCons, "&")

	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%X", string(h.Sum(nil)))
}

func Init(args *InitArgs) *Context {
	ctx := &Context{}

	ctx.merchantId = args.MerchantId
	ctx.apiKey = args.ApiKey
	ctx.appId = args.AppId

	ctx.tradeNumberGenerator = args.TradeNumberGenerator
	if nil == ctx.tradeNumberGenerator {
		ctx.tradeNumberGenerator = GeneratorTradeNumber
	}

	ctx.HttpClient = resty.New()
	if args.HttpTimeout > 0 {
		ctx.HttpClient.SetTimeout(args.HttpTimeout)
	}

	rand.Seed(time.Now().Unix())
	return ctx
}

type BaseResult struct {
	CgiCode    int    `xml:"retcode"`
	CgiMessage string `xml:"retmsg"`
	BizSuccess bool
	BizResult  string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrDesc    string `xml:"err_code_des"`
}

func (br *BaseResult) process() {
	br.BizSuccess = br.BizResult == "SUCCESS"
}
