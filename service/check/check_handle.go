package check

import (
	"ZJGSU-StuChecker-Go/model"
	"ZJGSU-StuChecker-Go/service/httpkit"
	"ZJGSU-StuChecker-Go/service/wechat"
	"fmt"
	"github.com/myleo1/go-core-kit/class/exception"
	"github.com/myleo1/go-core-kit/library/commonkit"
	"github.com/myleo1/go-core-kit/library/mathkit"
	"github.com/myleo1/go-core-kit/library/timekit"
	"github.com/myleo1/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/http"
	"sync"
	"time"
)

// BeginCheck 开始打卡
func BeginCheck(checkerList []*model.Checker) {
	if len(checkerList) == 0 {
		logkit.Fatal("未读取到用户")
	}
	var wg sync.WaitGroup
	wg.Add(len(checkerList))
	for _, v := range checkerList {
		checker := v
		if _preHandle(checker) {
			//参数校验通过 开始打卡流程
			//每个用户一个协程 因为有不同的阻塞时间（每个人打卡前等待时间不同）
			commonkit.RecoverGoFuncWrapper(func() {
				defer wg.Done()
				if checker.Skip {
					//跳过打卡
					logkit.Info(fmt.Sprintf("【%s】跳过打卡...", checker.Name))
					return
				}
				now := time.Now().UTC()
				wait := mathkit.RandInt32(waitMin, waitMax)
				startCheckDt := now.Add(time.Duration(wait) * time.Second).Format(timekit.TimeLayout)
				logkit.Info(fmt.Sprintf("【%s】将在 UTC <%s> 打卡...", checker.Name, startCheckDt))
				time.Sleep(time.Duration(wait) * time.Second)
				_check(checker)
				logkit.Info(fmt.Sprintf("【%s】打卡完毕...", checker.Name))
			})
		} else {
			wg.Done()
		}
	}
	wg.Wait()
	logkit.Info(fmt.Sprintf("---------全部打卡完毕---------"))
}

//校验config.json中的参数
func _preHandle(checker *model.Checker) bool {
	if checker.Name != "" && checker.Username != "" && checker.Password != "" && checker.Campus != "" {
		return true
	}
	return false
}

func _check(checker *model.Checker) {
	//按校区生成经纬度
	var long, lat float64
	var place string
	switch checker.Campus {
	case CampusJSG:
		long = mathkit.FloatRound(mathkit.RandFloat64(JSGLongitudeStart, JSGLongitudeEnd), 6)
		lat = mathkit.FloatRound(mathkit.RandFloat64(JSGLatitudeStart, JSGLatitudeEnd), 6)
		place = JSGPlaceName
	case CampusQJW:
		long = mathkit.FloatRound(mathkit.RandFloat64(QJWLongitudeStart, QJWLongitudeEnd), 6)
		lat = mathkit.FloatRound(mathkit.RandFloat64(QJWLatitudeStart, QJWLatitudeEnd), 6)
		place = QJWPlaceName
	case CampusJGL:
		long = mathkit.FloatRound(mathkit.RandFloat64(JGLLongitudeStart, JGLLongitudeEnd), 6)
		lat = mathkit.FloatRound(mathkit.RandFloat64(JGLLatitudeStart, JGLLatitudeEnd), 6)
		place = JGLPlaceName
	default:
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "校区填写错误"))
	}
	//获取access_token
	resp, err := httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    "https://uia.zjgsu.edu.cn/cas/mobile/getAccessToken",
		FormData: map[string]string{
			"clientId":    clientId,
			"mobileBT":    mobileBT,
			"redirectUrl": redirectUrl,
			"username":    checker.Username,
			"password":    checker.Password,
		},
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "获取access_token失败"))
	}
	res, err := resp.RespBody2Str()
	if err != nil {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + err.Error()))
	}
	accessToken := gjson.Get(res, "access_token")
	//cas认证
	resp, err = httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
		Url:    fmt.Sprintf("https://uia.zjgsu.edu.cn/cas/login?service=https://myapp.zjgsu.edu.cn/home/index&access_token=%s&mobileBT=%s", accessToken, mobileBT),
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "cas认证失败"))
	}
	resp, err = httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
		Url:    "https://ticket.zjgsu.edu.cn/stucheckservice/auth/login/stuCheck",
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "获取打卡系统token失败"))
	}
	token := resp.GetQueryParam("token")
	//打卡
	resp, err = httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    "https://ticket.zjgsu.edu.cn/stucheckservice/service/stuclockin",
		Header: map[string]string{
			"token":    token,
			"app-name": "stuCheck",
		},
		JsonData: map[string]string{
			"coordinate": fmt.Sprintf("%s,%s", cast.ToString(long), cast.ToString(lat)),
			"place":      place,
		},
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "POST打卡接口失败"))
	}
	res, err = resp.RespBody2Str()
	if err != nil {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + err.Error()))
	}
	if ok := gjson.Get(res, "success").Bool(); ok {
		wechat.Push2Wechat(checker, model.SuccessCheck)
	} else {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "打卡失败"))
	}
}

func errHandle(checker *model.Checker) {
	wechat.Push2Wechat(checker, model.FailCheck)
}
