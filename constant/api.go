package constant

var ServerAPIDomains = []string{"pcas.local", "pcas-test.cloudtrust.com.cn", "pcas.cloudtrust.com.cn"}

// 后端服务地址
const ServerAddressDev = "http://pcas.local:30080/api/org-am/"

const ServerAddressTest = "https://pcas-test.cloudtrust.com.cn/api/org-am/"

const ServerAddressProd = "https://pcas.cloudtrust.com.cn/api/org-am/"

// 上报上网行为记录接口
const UploadUrl = "behavior/record/report"
