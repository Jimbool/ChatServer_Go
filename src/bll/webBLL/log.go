package webBLL

import (
	"encoding/json"
	"github.com/Jordanzuo/ChatServer_Go/src/dal/requestLogDAL"
	"net/http"
)

func writeRequestLog(apiName string, r *http.Request) error {
	log, _ := json.Marshal(r.Form)
	return requestLogDAL.Insert(apiName, string(log))
}
