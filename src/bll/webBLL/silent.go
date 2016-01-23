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
	"time"
)

func silentCallback(w http.ResponseWriter, r *http.Request) {
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
	writeRequestLog("silent", r)

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
		responseObj.SetResultStatus(responseDataObject.DataError)
		return
	}
	if duration, err = strconv.Atoi(duration_str); err != nil {
		responseObj.SetResultStatus(responseDataObject.DataError)
		return
	}

	// 验证类型是否正确(0:查看禁言状态 1:禁言 2:解禁)
	if _type != 0 && _type != 1 && _type != 2 {
		responseObj.SetResultStatus(responseDataObject.DataError)
		return
	}

	// 验证签名
	if verifySilentSign(playerId, _type, duration, sign) == false {
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
		data := make(map[string]interface{}, 2)
		leftMinutes := (playerObj.SilentEndTime.Unix() - time.Now().Unix()) / 60
		data["Status"] = leftMinutes > 0
		if leftMinutes > 0 {
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

		playerBLL.UpdateSilentStatus(playerObj, silentEndTime)
	}
}

func verifySilentSign(playerId string, _type int, duration int, sign string) bool {
	rawstring := fmt.Sprintf("%s-%d-%d-%s-%s", playerId, _type, duration, configBLL.AppId(), configBLL.AppKey())
	if sign == securityUtil.Md5String(rawstring, false) {
		return true
	}

	return false
}
