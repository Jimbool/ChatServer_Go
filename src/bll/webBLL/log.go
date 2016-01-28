package webBLL

import (
	"encoding/json"
	"github.com/Jordanzuo/ChatServer_Go/src/dal/requestLogDAL"
	"net/http"
)

func writeRequestLog(apiName string, r *http.Request) {
	log, _ := json.Marshal(r.Form)
	go requestLogDAL.Insert(apiName, string(log))
}
