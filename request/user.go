package request

import (
	"fmt"
	"forum/database"
	"log"
	"net/http"
)

func getUsers(w http.ResponseWriter, r *http.Request, user database.User)  {
	var users []database.User
	var resUsers []database.User

	for index := 0; len(users); index++ {
		/*if (Strstr(getAjaxValue,users[index]) {
			resUsers.append(resUsers, users[index])
		}*/
	}
	return resUsers
}

func Strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}