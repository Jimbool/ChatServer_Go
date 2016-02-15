package webBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"net/http"
)

var (
	pushAPIName = "push"
)

func init() {
	registerAPI(pushAPIName, pushCallback)
}

func pushCallback(w http.ResponseWriter, r *http.Request) *responseDataObject.WebResponseObject {
	r.ParseForm()
	responseObj := responseDataObject.NewWebResponseObject()

	// 记录日志
	err := writeRequestLog(pushAPIName, r)
	if err != nil {
		logUtil.Log(err.Error(), logUtil.Error, true)
		responseObj.SetResultStatus(responseDataObject.DataError)
		return responseObj
	}

	// 解析数据
	message := r.Form["Message"][0]
	sign := r.Form["Sign"][0]

	// 验证签名
	if verifyPushSign(message, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 推送数据
	go playerBLL.PushMessage(message)

	return responseObj
}

func verifyPushSign(message string, sign string) bool {
	rawstring := fmt.Sprintf("%s-%s-%s", message, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
