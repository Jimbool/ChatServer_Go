package playerDAL

import (
	"errors"
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
)

func GetGamePlayer(id string) (string, string, bool, error) {
	sql := "SELECT p.Name, g.GuildId FROM p_player p LEFT JOIN p_guild_info g ON p.Id = g.PlayerId WHERE p.Id = ?;"
	rows, err := dal.GameDB().Query(sql, id)
	if err != nil {
		return "", "", false, errors.New(fmt.Sprintf("Query失败，错误信息：%s，sql:%s", err, sql))
	}

	var name string
	var unionId string
	for rows.Next() {
		var guildId interface{}
		err := rows.Scan(&name, &guildId)
		if err != nil {
			rows.Close()
			return "", "", false, errors.New(fmt.Sprintf("Scan失败，错误信息：%s，sql:%s", err, sql))
		}

		if guildId != nil {
			if unionIdArr, ok := guildId.([]byte); ok {
				unionId = string(unionIdArr)
			}
		}
	}

	return name, unionId, name != "", nil
}
