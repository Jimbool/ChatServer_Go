package webBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"net/http"
)

var (
	// 处理的方法映射表
	funcMap = make(map[string]func(http.ResponseWriter, *http.Request) *responseDataObject.WebResponseObject)
)

// 注册API
// apiName：API名称
// callback：回调方法
func registerAPI(apiName string, callback func(http.ResponseWriter, *http.Request) *responseDataObject.WebResponseObject) {
	funcMap[fmt.Sprintf("/API/%s", apiName)] = callback
}

// 定义自定义的Mux对象
type SelfDefineMux struct {
}

func (mux *SelfDefineMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseObj := responseDataObject.NewWebResponseObject()

	// 判断是否是POST方法
	if r.Method != "POST" {
		responseObj.SetResultStatus(responseDataObject.OnlySupportPOST)
		responseResult(w, responseObj)
		return
	}

	// 根据路径选择不同的处理方法
	if f, ok := funcMap[r.RequestURI]; !ok {
		responseObj.SetResultStatus(responseDataObject.APINotDefined)
		responseResult(w, responseObj)
	} else {
		responseObj = f(w, r)
		responseResult(w, responseObj)
	}
}

func responseResult(w http.ResponseWriter, responseObj *responseDataObject.WebResponseObject) {
	responseBytes, _ := json.Marshal(responseObj)
	fmt.Fprintln(w, string(responseBytes))
}
