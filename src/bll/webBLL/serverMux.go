package webBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"net/http"
)

var (
	// 处理的方法映射表
	funcMap = make(map[string]func(http.ResponseWriter, *http.Request))
)

func init() {
	funcMap["/API/forbid"] = forbidCallback
	funcMap["/API/silent"] = silentCallback
	funcMap["/API/push"] = pushCallback
	funcMap["/API/sensitive"] = sensitiveCallback
}

// 定义自定义的Mux对象
type SelfDefineMux struct {
}

func (mux *SelfDefineMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 定义返回值
	responseObj := responseDataObject.NewWebResponseObject()

	// 最终返回数据
	defer func() {
		// 捕获异常
		if r := recover(); r != nil {
			logUtil.LogUnknownError(r)
			responseObj.SetResultStatus(responseDataObject.DataError)
		}

		// 只处理失败的情况；正确地情况已经转到各个具体的方法里面去处理了
		if responseObj.Code != responseDataObject.Success {
			responseBytes, _ := json.Marshal(responseObj)
			fmt.Fprintln(w, string(responseBytes))
		}
	}()

	// 判断是否是POST方法
	if r.Method != "POST" {
		responseObj.SetResultStatus(responseDataObject.OnlySupportPOST)
		return
	}

	// 根据路径选择不同的处理方法
	if f, ok := funcMap[r.RequestURI]; ok {
		f(w, r)
	} else {
		responseObj.SetResultStatus(responseDataObject.APINotDefined)
	}
}
