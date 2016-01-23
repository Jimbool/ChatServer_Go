package webBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"net/http"
)

func pushCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	responseObj := responseDataObject.NewWebResponseObject()

	defer func() {
		// 捕获异常
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
			responseObj.SetResultStatus(responseDataObject.DataError)
		}

		// 输出结果给客户端
		responseBytes, _ := json.Marshal(responseObj)
		fmt.Fprintf(w, string(responseBytes))
	}()

	// 添加日志
	writeRequestLog("push", r)

	// 解析数据
	message := r.Form["Message"][0]
	sign := r.Form["Sign"][0]

	// 验证签名
	if verifyPushSign(message, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return
	}

	// 推送数据
	playerBLL.PushMessage(message)
}

func verifyPushSign(message string, sign string) bool {
	rawstring := fmt.Sprintf("%s-%s-%s", message, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
