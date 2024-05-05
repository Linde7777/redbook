package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"main/internal/service"
	"main/pkg/ratelimiter"

	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSMSService(redisCmd redis.Cmdable) service.SMSService {
	// SecretId、SecretKey 查询: https://console.cloud.tencent.com/cam/capi */
	credential := common.NewCredential(
		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"SecretId",
		"SecretKey",
	)
	/* 非必要步骤:

	- 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK默认使用POST方法。

	- 如果您一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK有默认的超时时间，非必要请不要进行调整

	- 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5

	/* 指定接入地域域名，默认就近地域接入域名为 sms.tencentcloudapi.com ，也支持指定地域域名访问，例如广州地域的域名为 sms.ap-guangzhou.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化要请求产品(以sms为例)的client对象

	- 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考 https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := tencentSMS.NewClient(credential, "ap-guangzhou", cpf)

	rateLimiter := ratelimiter.NewRedisSlidingWinLimiter(
		redisCmd, 10, 60)
	return service.NewTencentSMSServiceWithLimiter(client, "appID",
		"signature", rateLimiter, "tencent-sms-ratelimiter")
}
