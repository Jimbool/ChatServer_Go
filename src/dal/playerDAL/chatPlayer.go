package playerDAL

import (
	"database/sql"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func GetPlayer(id string) (playerObj *player.Player, exists bool, err error) {
	command := "SELECT Name, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime FROM player WHERE Id = ?;"

	var name string
	var unionId string
	var extraMsg string
	var registerTime time.Time
	var loginTime time.Time
	var isForbidden bool
	var silentEndTime time.Time
	if err = dal.ChatDB().QueryRow(command, id).Scan(&name, &unionId, &extraMsg, &registerTime, &loginTime, &isForbidden, &silentEndTime); err != nil {
		if err == sql.ErrNoRows {
			// 重置err，使其为nil；因为这代表的是没有查找到数据，而不是真正的错误
			err = nil
			return
		} else {
			logUtil.Log(fmt.Sprintf("Scan失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
			return
		}
	}

	playerObj = player.NewPlayer(id, name, unionId, extraMsg, registerTime, loginTime, isForbidden, silentEndTime)
	exists = true

	return
}

func Insert(player *player.Player) error {
	command := `INSERT INTO 
                player(Id, Name, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime)
            VALUES
                (?, ?, ?, ?, ?, ?, ?, ?);
    `
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.Id, player.Name, player.UnionId, player.ExtraMsg, player.RegisterTime, player.LoginTime, player.IsForbidden, player.SilentEndTime); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}

func UpdateInfo(player *player.Player) error {
	command := "UPDATE player SET Name = ?, UnionId = ?, ExtraMsg = ? WHERE Id = ?"
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.Name, player.UnionId, player.ExtraMsg, player.Id); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}

func UpdateLoginTime(player *player.Player) error {
	command := "UPDATE player SET LoginTime = ? WHERE Id = ?"
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.LoginTime, player.Id); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}

func UpdateForbiddenStatus(player *player.Player) error {
	command := "UPDATE player SET IsForbidden = ? WHERE Id = ?"
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.IsForbidden, player.Id); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}

func UpdateSilentEndTime(player *player.Player) error {
	command := "UPDATE player SET SilentEndTime = ? WHERE Id = ?"
	stmt, err := dal.ChatDB().Prepare(command)
	if err != nil {
		logUtil.Log(fmt.Sprintf("Prepare失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	// 最后关闭
	defer stmt.Close()

	if _, err = stmt.Exec(player.SilentEndTime, player.Id); err != nil {
		logUtil.Log(fmt.Sprintf("Exec失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
		return err
	}

	return nil
}
