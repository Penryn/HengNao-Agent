package utils

import (
	"github.com/go-resty/resty/v2"
	"time"
)

var RestyClient *resty.Client

func InitRestyClient() {
	RestyClient = resty.New().
		SetRetryCount(3).            // 设置重试次数
		SetRetryWaitTime(5).         // 设置重试等待时间
		SetRetryMaxWaitTime(10).     // 设置最大重试等待时间
		SetTimeout(30 * time.Second) // 设置请求超时时间
}
