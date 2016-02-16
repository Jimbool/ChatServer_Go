package webBLL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/configBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/bll/playerBLL"
	"github.com/Jordanzuo/ChatServer_Go/src/model/responseDataObject"
	"github.com/Jordanzuo/goutil/securityUtil"
	"net/http"
	"strconv"
	"time"
)

var (
	silentAPIName = "silent"
)

func init() {
	registerAPI(silentAPIName, silentCallback)
}

func silentCallback(w http.ResponseWriter, r *http.Request) *responseDataObject.WebResponseObject {
	r.ParseForm()
	responseObj := responseDataObject.NewWebResponseObject()

	// 记录日志
	if err := writeRequestLog(silentAPIName, r); err != nil {
		responseObj.SetDataError()
		return responseObj
	}

	// 解析数据
	playerId := r.Form["PlayerId"][0]
	_type_str := r.Form["Type"][0]
	duration_str := r.Form["Duration"][0] // 单位：分钟
	sign := r.Form["Sign"][0]

	// 类型转换
	var _type int
	var duration int
	var err error

	if _type, err = strconv.Atoi(_type_str); err != nil {
		responseObj.SetAPIDataError()
		return responseObj
	}
	if duration, err = strconv.Atoi(duration_str); err != nil {
		responseObj.SetAPIDataError()
		return responseObj
	}

	// 验证类型是否正确(0:查看禁言状态 1:禁言 2:解禁)
	if _type != 0 && _type != 1 && _type != 2 {
		responseObj.SetAPIDataError()
		return responseObj
	}

	// 验证签名
	if verifySilentSign(playerId, _type, duration, sign) == false {
		responseObj.SetResultStatus(responseDataObject.SignError)
		return responseObj
	}

	// 判断玩家是否存在
	playerObj, exists, err := playerBLL.GetPlayer(playerId, true)
	if err != nil {
		responseObj.SetDataError()
		return responseObj
	}
	if !exists {
		responseObj.SetResultStatus(responseDataObject.PlayerNotExist)
		return responseObj
	}

	// 判断是否为查询状态
	if _type == 0 {
		data := make(map[string]interface{}, 2)
		isInSilent, leftMinutes := playerObj.IsInSilent()
		data["Status"] = isInSilent
		if isInSilent {
			data["LeftMinutes"] = leftMinutes
		}
		responseObj.SetData(data)
	} else {
		// 修改禁言状态
		silentEndTime := time.Now()
		if _type == 1 {
			if duration == 0 {
				silentEndTime = silentEndTime.AddDate(10, 0, 0)
			} else {
				silentEndTime = silentEndTime.Add(time.Duration(duration) * time.Minute)
			}
		}

		if err := playerBLL.UpdateSilentStatus(playerObj, silentEndTime); err != nil {
			responseObj.SetDataError()
			return responseObj
		}
	}

	return responseObj
}

func verifySilentSign(playerId string, _type int, duration int, sign string) bool {
	rawstring := fmt.Sprintf("%s-%d-%d-%s-%s", playerId, _type, duration, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
