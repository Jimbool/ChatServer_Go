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
	"strconv"
)

func forbidCallback(w http.ResponseWriter, r *http.Request) {
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
	writeRequestLog("forbid", r)

	// 解析数据
	playerId := r.Form["PlayerId"][0]
	_type_str := r.Form["Type"][0]
	sign := r.Form["Sign"][0]

	// 类型转换
	var _type int
	var err error

	if _type, err = strconv.Atoi(_type_str); err != nil {
		responseObj.SetResultStatus(responseDataObject.DataError)
		return
	}

	// 验证类型是否正确(0:查看封号状态 1:封号 2:解封)
	if _type != 0 && _type != 1 && _type != 2 {
		responseObj.SetResultStatus(responseDataObject.DataError)
		return
	}

	// 验证签名
	if verifyForbidSign(playerId, _type, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return
	}

	// 判断玩家是否存在
	playerObj, ok := playerBLL.GetPlayer(playerId, true)
	if !ok {
		responseObj.SetResultStatus(responseDataObject.PlayerNotExist)
		return
	}

	// 判断是否为查询状态
	if _type == 0 {
		responseObj.SetData(playerObj.IsForbidden)
	} else {
		// 修改封号状态
		go playerBLL.UpdateForbidStatus(playerObj, _type == 1)
	}
}

func verifyForbidSign(playerId string, _type int, sign string) bool {
	rawstring := fmt.Sprintf("%s-%d-%s-%s", playerId, _type, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
