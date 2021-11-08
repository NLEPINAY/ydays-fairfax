package request

import (
	"forum/database"
	"net/http"
)

func Messages(w http.ResponseWriter, r *http.Request, user database.User) {
	if r.URL.Path != "/messages" {
		err := MyTemplates.ExecuteTemplate(w, "400", user)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
		// http.Error(w, "404 NOT FOUND", http.StatusNotFound)
	}

	err5 := MyTemplates.ExecuteTemplate(w, "messages", nil)
	if err5 != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

}
