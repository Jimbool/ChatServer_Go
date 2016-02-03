package playerDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
)

func GetGamePlayer(id string) (string, string, bool) {
	sql := "SELECT p.Name, g.GuildId FROM p_player p LEFT JOIN p_guild_info g ON p.Id = g.PlayerId WHERE p.Id = ?;"
	rows, err := dal.GameDB().Query(sql, id)
	if err != nil {
		panic(err)
	}

	var name string
	var unionId string
	for rows.Next() {
		var guildId interface{}
		err := rows.Scan(&name, &guildId)
		if err != nil {
			panic(err)
		}

		if guildId != nil {
			if unionIdArr, ok := guildId.([]byte); ok {
				unionId = string(unionIdArr)
			}
		}
	}

	return name, unionId, name != ""
}
