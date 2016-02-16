package webBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/sensitiveWordsBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/securityUtil"
	"net/http"
)

var (
	sensitiveAPIName = "sensitive"
)

func init() {
	registerAPI(sensitiveAPIName, sensitiveCallback)
}

func sensitiveCallback(w http.ResponseWriter, r *http.Request) *responseDataObject.WebResponseObject {
	r.ParseForm()
	responseObj := responseDataObject.NewWebResponseObject()

	// 记录日志
	if err := writeRequestLog(sensitiveAPIName, r); err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	// 解析数据
	message := r.Form["Message"][0]
	sign := r.Form["Sign"][0]

	// 验证签名
	if verifySensitiveSign(message, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 重新加载敏感词
	if err := sensitiveWordsBLL.Reload(); err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	return responseObj
}

func verifySensitiveSign(message string, sign string) bool {
	rawstring := fmt.Sprintf("%s-%s-%s", message, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
