package check

//至少等待30s
const waitMin = 30

//todo 定时任务开始后7200s(2h)内随机开始打卡
const waitMax = 7200

const (
	clientId    = "qnFZATsB6D25EnZeII"
	mobileBT    = "CEB5F29A-39B1-4662-ABFE-E82475475A87"
	redirectUrl = "https://myapp.zjgsu.edu.cn/home/index"
)

const (
	CampusJSG = "金沙港"
	CampusQJW = "钱江湾"
	CampusJGL = "教工路"
)

//金沙港
const (
	JSGPlaceName      = "浙江省,杭州市,钱塘区,学林街,浙江省杭州市钱塘区学林街靠近浙江工商大学金沙港生活园区"
	JSGLongitudeStart = 120.379535
	JSGLongitudeEnd   = 120.381493
	JSGLatitudeStart  = 30.311940
	JSGLatitudeEnd    = 30.313440
)

//钱江湾
const (
	QJWPlaceName      = "浙江省,杭州市,钱塘区,学林街,浙江省杭州市钱塘区学林街靠近浙江工商大学钱江湾生活园区"
	QJWLongitudeStart = 120.388574
	QJWLongitudeEnd   = 120.396202
	QJWLatitudeStart  = 30.311300
	QJWLatitudeEnd    = 30.312162
)

//教工路
const (
	JGLPlaceName      = "浙江省,杭州市,西湖区,保俶北路,浙江省杭州市西湖区保俶北路靠近浙江工商大学教工路校区"
	JGLLongitudeStart = 120.135760
	JGLLongitudeEnd   = 120.139488
	JGLLatitudeStart  = 30.284914
	JGLLatitudeEnd    = 30.285812
)
