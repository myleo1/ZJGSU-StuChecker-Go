package wechat

import (
	"ZJGSU-StuChecker-Go/model"
	"ZJGSU-StuChecker-Go/service/httpkit"
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"net/http"
	"time"
)

func Push2Wechat(checker *model.Checker, title string) {
	if checker.WechatPushKey == "" {
		return
	}
	description := ""
	if title == model.SuccessCheck {
		description = fmt.Sprintf("<div class=\"gray\">%s\n</div> <div class=\"normal\">恭喜您%s~,学号[%s]打卡成功！\n</div><div class=\"highlight\">点击卡片查看打卡源码~</div>", time.Now().UTC().Add(time.Hour*8).Format(timekit.TimeLayoutYMD), checker.Name, checker.Username)
	} else {
		description = fmt.Sprintf("<div class=\"gray\">%s\n</div> <div class=\"normal\">很遗憾%s~,学号[%s]打卡失败！\n</div><div class=\"highlight\">点击卡片查看打卡源码~</div>", time.Now().UTC().Add(time.Hour*8).Format(timekit.TimeLayoutYMD), checker.Name, checker.Username)
	}
	resp, err := httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    configkit.GetStringD(ApiWechat),
		Header: map[string]string{
			"Cookie": fmt.Sprintf("session=%s", configkit.GetStringD(TokenWechat)),
		},
		FormData: map[string]string{
			"to":          checker.WechatPushKey,
			"title":       title,
			"description": description,
			"url":         "https://github.com/myleo1/ZJGSU-StuChecker-Go",
		},
		Timeout: 5,
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		panic(exception.New("微信推送失败"))
	}
}
