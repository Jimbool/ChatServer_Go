package playerDAL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func GetPlayer(id string) (*player.Player, bool) {
	sql := "SELECT Id, Name, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime FROM player WHERE Id = ?;"
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	rows, err := dal.ChatDB().Query(sql, id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id string
		var name string
		var unionId string
		var extraMsg string
		var registerTime time.Time
		var loginTime time.Time
		var isForbidden bool
		var silentEndTime time.Time
		err := rows.Scan(&id, &name, &unionId, &extraMsg, &registerTime, &loginTime, &isForbidden, &silentEndTime)
		if err != nil {
			panic(err)
		}

		return player.NewPlayer(id, name, unionId, extraMsg, registerTime, loginTime, isForbidden, silentEndTime), true
	}

	return nil, false
}

func Insert(player *player.Player) {
	sql := `INSERT INTO 
                player(Id, Name, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime)
            VALUES
                (?, ?, ?, ?, ?, ?, ?, ?);
    `
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(player.Id, player.Name, player.UnionId, player.ExtraMsg, player.RegisterTime, player.LoginTime, player.IsForbidden, player.SilentEndTime)
	if err != nil {
		panic(err)
	}
}

func UpdateInfo(player *player.Player) {
	sql := "UPDATE player SET Name = ?, UnionId = ?, ExtraMsg = ? WHERE Id = ?"
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.UnionId, player.ExtraMsg, player.Id)
	if err != nil {
		panic(err)
	}
}

func UpdateLoginTime(player *player.Player) {
	sql := "UPDATE player SET LoginTime = ? WHERE Id = ?"
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(player.LoginTime, player.Id)
	if err != nil {
		panic(err)
	}
}

func UpdateForbiddenStatus(player *player.Player) {
	sql := "UPDATE player SET IsForbidden = ? WHERE Id = ?"
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(player.IsForbidden, player.Id)
	if err != nil {
		panic(err)
	}
}

func UpdateSilentEndTime(player *player.Player) {
	sql := "UPDATE player SET SilentEndTime = ? WHERE Id = ?"
	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end
	stmt, err := dal.ChatDB().Prepare(sql)
	if err != nil {
		panic(err)
	}

	// 最后关闭
	defer stmt.Close()

	_, err = stmt.Exec(player.SilentEndTime, player.Id)
	if err != nil {
		panic(err)
	}
}
