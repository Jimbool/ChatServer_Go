package playerDAL

import (
	"database/sql"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
)

func GetGamePlayer(id string) (name string, unionId string, exists bool, err error) {
	command := "SELECT p.Name, g.GuildId FROM p_player p LEFT JOIN p_guild_info g ON p.Id = g.PlayerId WHERE p.Id = ?;"

	var guildId interface{}
	if err = dal.GameDB().QueryRow(command, id).Scan(&name, &guildId); err != nil {
		if err == sql.ErrNoRows {
			// 重置err，使其为nil；因为这代表的是没有查找到数据，而不是真正的错误
			err = nil
			return
		} else {
			logUtil.Log(fmt.Sprintf("Scan失败，错误信息：%s，command:%s", err, command), logUtil.Error, true)
			return
		}
	}
	exists = true

	// 处理公会Id
	if guildId != nil {
		if unionIdArr, ok := guildId.([]byte); ok {
			unionId = string(unionIdArr)
		}
	}

	return
}
