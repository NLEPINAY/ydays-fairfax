package request

import (
	"encoding/json"
	"forum/database"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Messagerie(w http.ResponseWriter, r *http.Request, user database.User) {
	switch r.Method {
	case "POST":
		var ERROR []byte
		var received database.ReceivedData
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			ERROR, _ = json.Marshal("ERROR WHILE DECODING JSON")
			w.Write(ERROR)
			panic(err)
		}

		break
	case "GET":

	}
}
