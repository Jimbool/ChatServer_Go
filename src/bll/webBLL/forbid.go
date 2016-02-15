package webBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/logUtil"
	"github.com/Jordanzuo/goutil/securityUtil"
	"net/http"
	"strconv"
)

var (
	forbidAPIName = "forbid"
)

func init() {
	registerAPI(forbidAPIName, forbidCallback)
}

func forbidCallback(w http.ResponseWriter, r *http.Request) *responseDataObject.WebResponseObject {
	r.ParseForm()
	responseObj := responseDataObject.NewWebResponseObject()

	// 记录日志
	err := writeRequestLog(forbidAPIName, r)
	if err != nil {
		logUtil.Log(err.Error(), logUtil.Error, true)
		responseObj.SetResultStatus(responseDataObject.DataError)
		return responseObj
	}

	// 解析数据
	playerId := r.Form["PlayerId"][0]
	_type_str := r.Form["Type"][0]
	sign := r.Form["Sign"][0]

	// 类型转换
	var _type int
	if _type, err = strconv.Atoi(_type_str); err != nil {
		responseObj.SetResultStatus(responseDataObject.APIDataError)
		return responseObj
	}

	// 验证类型是否正确(0:查看封号状态 1:封号 2:解封)
	if _type != 0 && _type != 1 && _type != 2 {
		responseObj.SetResultStatus(responseDataObject.APIDataError)
		return responseObj
	}

	// 验证签名
	if verifyForbidSign(playerId, _type, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 判断玩家是否存在
	playerObj, ok, err := playerBLL.GetPlayer(playerId, true)
	if err != nil {
		responseObj.SetResultStatus(responseDataObject.DataError)
		return responseObj
	}

	if !ok {
		responseObj.SetResultStatus(responseDataObject.PlayerNotExist)
		return responseObj
	}

	// 判断是否为查询状态
	if _type == 0 {
		responseObj.SetData(playerObj.IsForbidden)
	} else {
		// 修改封号状态
		playerBLL.UpdateForbidStatus(playerObj, _type == 1)
	}

	return responseObj
}

func verifyForbidSign(playerId string, _type int, sign string) bool {
	rawstring := fmt.Sprintf("%s-%d-%s-%s", playerId, _type, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
