package captcha

import (
	"encoding/json"
	"loginhub/pkg/convert"
	"time"
)

const (
	sendRate = time.Minute
)

type Captcha struct {
	Code        string `json:"code,omitempty"`
	LastReqTime int64  `json:"last_req_time,omitempty"`
}

func (c *Captcha) IsCorrect(code string) bool {
	return c.Code == code
}

func (c *Captcha) Allow() bool {
	return time.Now().Unix()-c.LastReqTime > int64(sendRate)
}

func (c *Captcha) ToJSONString() string {
	b, _ := json.Marshal(c)
	return convert.BytesToString(b)
}

func (c *Captcha) FromJSONString(s string) *Captcha {
	_ = json.Unmarshal(convert.StringToBytes(s), c)
	return c
}

func NewCaptcha(code string) *Captcha {
	return &Captcha{
		Code:        code,
		LastReqTime: time.Now().Unix(),
	}
}
