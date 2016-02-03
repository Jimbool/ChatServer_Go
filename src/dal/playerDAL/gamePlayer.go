package playerDAL

import (
	"fmt"
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
	"time"
)

func GetGamePlayer(id string) (string, string, bool) {

	sql := "SELECT p.Name, g.GuildId FROM p_player p LEFT JOIN p_guild_info g ON p.Id = g.PlayerId WHERE p.Id = ?;"

	// log begin
	start := time.Now().Unix()
	defer func() {
		end := time.Now().Unix()
		logUtil.Log(fmt.Sprintf("%s执行耗时%d", sql, (end-start)), logUtil.Fatal, true)
	}()
	// log end

	rows, err := dal.GameDB().Query(sql, id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var name string
		var unionId string
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

		return name, unionId, true
	}

	return "", "", false
}
