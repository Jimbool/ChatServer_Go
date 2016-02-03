package playerDAL

import (
	"github.com/Jordanzuo/ChatServer_Go/src/dal"
	"github.com/Jordanzuo/goutil/logUtil"
	"sync/atomic"
	"time"	
	"fmt"
)

var(
	ChatPlayerTotalCount int32 = 0
	ChatPlayerSucceedCount int32= 0
	GamePlayerTotalCount int32 = 0
	GamePlayerSucceedCount int32= 0
)

func init() {
	go display()
}

func display() {
	for {
		time.Sleep(time.Minute)
		logUtil.Log(fmt.Sprintf("ChatPlayerTotalCount:%d，ChatPlayerSucceedCount:%d,GamePlayerTotalCount:%d,GamePlayerSucceedCount:%d", ChatPlayerTotalCount, ChatPlayerSucceedCount,GamePlayerTotalCount,GamePlayerSucceedCount), logUtil.Error, true)
	}	
}

func GetGamePlayer(id string) (string, string, bool) {	
	atomic.AddInt32(&ChatPlayerTotalCount, 1)
	if id == "7cd63ef5-d4b1-4204-a894-36b31731afc9"{
		logUtil.Log("301", logUtil.Fatal, true)		
	}	

	sql := "SELECT p.Name, g.GuildId FROM p_player p LEFT JOIN p_guild_info g ON p.Id = g.PlayerId WHERE p.Id = ?;"
	rows, err := dal.GameDB().Query(sql, id)
	if err != nil {
		panic(err)
	}

	atomic.AddInt32(&ChatPlayerSucceedCount, 1)

	if id == "7cd63ef5-d4b1-4204-a894-36b31731afc9"{
		logUtil.Log("302", logUtil.Fatal, true)
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
				if id == "7cd63ef5-d4b1-4204-a894-36b31731afc9"{
					logUtil.Log("303", logUtil.Fatal, true)
					logUtil.Log(unionId, logUtil.Fatal, true)
				}
			}

			if id == "7cd63ef5-d4b1-4204-a894-36b31731afc9"{
				logUtil.Log("304", logUtil.Fatal, true)
			}
		}
	}

	if id == "7cd63ef5-d4b1-4204-a894-36b31731afc9"{
		logUtil.Log("305", logUtil.Fatal, true)
	}

	return name, unionId, name != ""
}
