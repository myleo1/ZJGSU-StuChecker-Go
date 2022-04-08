package model

type CheckerList []*Checker

type Checker struct {
	Name          string `json:"name" description:"姓名"`
	Username      string `json:"username" description:"学号"`
	Password      string `json:"password" description:"教务系统密码"`
	Campus        string `json:"campus" description:"校区 [金沙港/钱江湾/教工路]"`
	WechatPushKey string `json:"wechatPushKey" description:"微信推送账号"`
	DeviceId      string `json:"deviceId" description:"识别码/每个手机都不一样"`
	Skip          bool   `json:"skip" description:"是否跳过打卡"`
}
