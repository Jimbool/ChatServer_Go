/*
玩家数据处理包
*/
package playerDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/ChatServer_Go/src/model/player"
	"time"
)

func GetPlayer(id string) (*player.Player, bool) {
	sql := "SELECT Id, Name, UnionId, ExtraMsg, RegisterTime, LoginTime, IsForbidden, SilentEndTime FROM player WHERE Id = ?;"
	rows, err := dal.DB().Query(sql, id)
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
	stmt, err := dal.DB().Prepare(sql)
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
	stmt, err := dal.DB().Prepare(sql)
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
	stmt, err := dal.DB().Prepare(sql)
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
	stmt, err := dal.DB().Prepare(sql)
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
	stmt, err := dal.DB().Prepare(sql)
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
