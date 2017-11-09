package qpay

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"
)

func GeneratorTradeNumber() string {
	return fmt.Sprintf("%10d%05d", time.Now().UnixNano(), rand.Int31n(100000))
}

func GeneratorNonce() string {
	return fmt.Sprintf("%10d%05d", time.Now().UnixNano(), rand.Int31n(100000))
}

type StringMap map[string]string

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "xml"}
	if err := e.EncodeToken(start); nil != err {
		return err
	}
	for k, v := range m {
		if err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}}); nil != err {
			return err
		}
	}
	return e.EncodeToken(start.End())
}
