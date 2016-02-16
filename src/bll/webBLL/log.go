package webBLL

import (
	"encoding/json"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal/requestLogDAL"
	"github.com/Jordanzuo/goutil/logUtil"
	"net/http"
)

func writeRequestLog(apiName string, r *http.Request) error {
	log, err := json.Marshal(r.Form)
	if err != nil {
		logUtil.Log(fmt.Sprintf("序列化数据错误，原始数据位：%v，错误信息为：%s", r.Form, err), logUtil.Error, true)
		return err
	}

	return requestLogDAL.Insert(apiName, string(log))
}
