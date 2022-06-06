package check

import (
	"ZJGSU-StuChecker-Go/model"
	"ZJGSU-StuChecker-Go/service/httpkit"
	"ZJGSU-StuChecker-Go/service/wechat"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/bytekit"
	"github.com/mizuki1412/go-core-kit/library/commonkit"
	"github.com/mizuki1412/go-core-kit/library/cryptokit"
	"github.com/mizuki1412/go-core-kit/library/mathkit"
	"github.com/mizuki1412/go-core-kit/library/timekit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex

// BeginCheck 开始打卡
// Deprecated
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

// BeginYzy 开始云战役
func BeginYzy(checkerList []*model.Checker) {
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
				wait := mathkit.RandInt32(waitMinYzy, waitMaxYzy)
				startCheckDt := now.Add(time.Duration(wait) * time.Second).Format(timekit.TimeLayout)
				logkit.Info(fmt.Sprintf("【%s】将在 UTC <%s> 打卡...", checker.Name, startCheckDt))
				//time.Sleep(time.Duration(wait) * time.Second)
				_checkYzy(checker)
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

func _checkCommon(checker *model.Checker) *model.PostInfo {
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
	info := &model.PostInfo{
		Long:  long,
		Lat:   lat,
		Place: place,
	}
	return info
}

// Deprecated
func _check(checker *model.Checker) {
	mutex.Lock()
	defer mutex.Unlock()
	info := _checkCommon(checker)
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
			"coordinate": fmt.Sprintf("%s,%s", cast.ToString(info.Long), cast.ToString(info.Lat)),
			"place":      info.Place,
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

func _checkYzy(checker *model.Checker) {
	mutex.Lock()
	defer mutex.Unlock()
	info := _checkCommon(checker)
	//22.06.06 update 重新获取了web端获取token的方式
	resp, err := httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    "https://yzy.zjgsu.edu.cn/cloudbattleservice/service/cloudLogin",
		Header: map[string]string{
			"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.18(0x1800123f) NetType/WIFI Language/zh_CN",
		},
		JsonData: map[string]string{
			"gh":    checker.Username,
			"psswd": checker.Password,
		},
		Timeout: 5,
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "获取打卡系统token失败"))
	}
	res, err := resp.RespBody2Str()
	if err != nil {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "获取打卡系统token失败"))
	}
	if ok := gjson.Get(res, "success").Bool(); !ok {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "获取打卡系统token失败"))
	}
	token := gjson.Get(res, "data.token").String()

	//22.05.28 update 以上方案获取token失效,需要绑定微信,操作起来太麻烦,最终还是走我的商大的API获取token
	//获取access_token
	//resp, err := httpkit.Request(httpkit.Req{
	//	Method: http.MethodPost,
	//	Url:    "https://uia.zjgsu.edu.cn/cas/mobile/getAccessToken",
	//	FormData: map[string]string{
	//		"clientId":    clientId,
	//		"mobileBT":    mobileBT,
	//		"redirectUrl": redirectUrl,
	//		"username":    checker.Username,
	//		"password":    checker.Password,
	//	},
	//})
	//if err != nil || resp.StatusCode() != http.StatusOK {
	//	errHandle(checker)
	//	panic(exception.New("【" + checker.Name + "】" + "获取access_token失败"))
	//}
	//res, err := resp.RespBody2Str()
	//if err != nil {
	//	errHandle(checker)
	//	panic(exception.New("【" + checker.Name + "】" + err.Error()))
	//}
	//accessToken := gjson.Get(res, "access_token")
	////cas认证
	//resp, err = httpkit.Request(httpkit.Req{
	//	Method: http.MethodGet,
	//	Url:    fmt.Sprintf("https://uia.zjgsu.edu.cn/cas/login?service=https://myapp.zjgsu.edu.cn/home/index&access_token=%s&mobileBT=%s", accessToken, mobileBT),
	//})
	//if err != nil || resp.StatusCode() != http.StatusOK {
	//	errHandle(checker)
	//	panic(exception.New("【" + checker.Name + "】" + "cas认证失败"))
	//}
	//resp, err = httpkit.Request(httpkit.Req{
	//	Method: http.MethodGet,
	//	Url:    "https://ticket.zjgsu.edu.cn/stucheckservice/auth/login/stuCheck",
	//})
	//if err != nil || resp.StatusCode() != http.StatusOK {
	//	errHandle(checker)
	//	panic(exception.New("【" + checker.Name + "】" + "获取打卡系统token失败"))
	//}
	//token := resp.GetQueryParam("token")
	t := cast.ToString(time.Now().UnixMilli()) + "26B0"
	plaintext := []byte("882D")
	plaintext = append(plaintext, []byte(t)...)
	zjgsuCheck := _AESCBCEncrypt(plaintext)
	keyCode := checker.Username + "*" + t + "^25A622DCE625882D8085CC9F00BF8C12"
	zjgsuAuth := cryptokit.MD5(keyCode)
	//云战役
	resp, err = httpkit.Request(httpkit.Req{
		Method: http.MethodPost,
		Url:    "https://yzy.zjgsu.edu.cn/cloudbattleservice/service/add",
		Header: map[string]string{
			"token":      token,
			"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;zfsoft",
			"zjgsuCheck": zjgsuCheck,
			"zjgsuAuth":  zjgsuAuth,
		},
		JsonData: map[string]string{
			"currentResd": info.Place,
			"fromHbToZj":  "C",
			"fromWtToHz":  "B",
			"meetCase":    "C",
			"travelCase":  "D",
			"medObsv":     "B",
			"belowCase":   "D",
			"hzQrCode":    "A",
			"specialDesc": "无",
			"deviceId":    checker.DeviceId,
			"fromDevice":  "App",
			"isNewEpid":   "否",
			"location":    info.Place,
			"coordinate":  fmt.Sprintf("%s,%s", cast.ToString(info.Long), cast.ToString(info.Lat)),
		},
		Timeout: 5,
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
	ok := gjson.Get(res, "message").String()
	if ok == "成功" {
		wechat.Push2Wechat(checker, model.SuccessCheck)
	} else if ok == "今日已打卡！" {
		wechat.Push2Wechat(checker, model.AlreadyCheck)
	} else {
		errHandle(checker)
		panic(exception.New("【" + checker.Name + "】" + "打卡失败"))
	}
}

func errHandle(checker *model.Checker) {
	wechat.Push2Wechat(checker, model.FailCheck)
}

func _padding(plaintext []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plaintext)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plaintext = append(plaintext, temp...)
	return plaintext
}

func _AESCBCEncrypt(plaintext []byte) string {
	block, err := aes.NewCipher([]byte("ED7925CF8acd26B0"))
	if err != nil {
		panic(exception.New(err.Error()))
	}
	//进行填充
	plaintext = _padding(plaintext, block.BlockSize())
	//指定初始向量，长度和block的块尺寸一致
	iv := []byte("3670759D768a359f")
	//指定分组模式，返回一个blockMode对象接口
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	//返回密文
	return bytekit.Bytes2HexString1(ciphertext)
}
