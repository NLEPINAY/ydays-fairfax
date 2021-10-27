package request

import (
	"encoding/json"
	"fmt"
	"forum/admin"
	"forum/database"
	"forum/toolbox"
	"net/http"
	"strconv"
)

var userCanDelete = []string{"posts", "comments", "gif", "images", "users", "tickets", "ticket_answers"}
var userCanAlterate = []string{"posts", "comments", "users", "tickets", "ticket_answers"}

//requete de modification de la bdd
func Fetching(w http.ResponseWriter, r *http.Request, user database.User) {
	fmt.Println(r.Body, "ttttttt")
	switch r.Method {
	//savoir quelle action est demander
	case "POST":
		var ERROR []byte
		var received database.ReceivedData
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			ERROR, _ = json.Marshal("ERROR WHILE DECODING JSON")
			w.Write(ERROR)
			panic(err)
		}
		switch received.Action {
		case "CREATE":

		case "UPDATE":
			if isLegal(received, user) {
				performAction(received)
				received.Message = "Success."
				ok, _ := json.Marshal(received)
				w.Write(ok)
			} else {
				received.Message = "DOESNT HAVE PERMISSION."
				ok, _ := json.Marshal(received)
				w.Write(ok)
			}
		case "DELETE":
			if isLegal(received, user) {
				performAction(received)
				received.Message = "Success."
				ok, _ := json.Marshal(received)
				w.Write(ok)
			} else {
				received.Message = "DOESNT HAVE PERMISSION."
				ok, _ := json.Marshal(received)
				w.Write(ok)
			}
		case "get":
			var Data admin.Data

			switch received.Table {
			case "Charts":
				Data = admin.GetCommentList(Data)
				Data = admin.GetClientList(Data)
				Data = admin.GetPostList(Data)
				Data.Category = database.GetCategoriesList()
			case "Category":
				Data.Category = database.GetCategoriesList()
			case "Post":
				Data = admin.GetPostList(Data)
				Data = admin.GetLikes(Data)
			case "User":
				Data = admin.GetClientList(Data)
			case "Comment":
				Data = admin.GetCommentList(Data)
			}
			ok, _ := json.Marshal(Data)
			w.Write(ok)
		case "getForUpdate": // Récupere un(e) article / catégorie / utilisateur / commentaire pour remplir automatiquement la modal d'update
			var Data admin.Data
			ID, err := strconv.Atoi(received.ID)
			if err != nil {
				panic(err)

			}
			switch received.Table {
			case "Post":

				Data = admin.GetPostOnlyByID(ID, Data)
				Data = admin.GetCategoriesList(Data)

			}
			ok, _ := json.Marshal(Data)
			w.Write(ok)

		case "getStats":
			var Data admin.DataForChart
			switch received.What {
			case "Category":
				Data = admin.GetCategoriesChart(Data)
			case "Post":
				Data = admin.GetPostChart(Data)
			case "User":
				Data = admin.GetUserChart(Data)
			}

			/*Data = admin.GetClientChart(Data)
			Data = admin.GetCommentChart(Data)
			Data = admin.GetPostChart(Data)
			Data.Category = database.GetCategoriesChart()*/
			ok, _ := json.Marshal(Data)
			w.Write(ok)
		}

	}
}

//Tout est vérifier je fait ce que la requête demande
func performAction(r database.ReceivedData) {
	var query string
	switch r.Action {
	case "UPDATE":
		query = "UPDATE " + r.Table + " SET " + r.What + " = \"" + r.NewValue + "\" WHERE id" + " = \"" + r.ID + "\""
		break
	case "DELETE":
		if r.Is == "cell" { //si c'est une cellule c'est un update sur une valeur null
			query = "UPDATE " + r.Table + " SET " + r.What + " = \"\" WHERE id" + " = \"" + r.ID + "\""
		} else { //une table je la supprime
			query = "DELETE FROM " + r.Table + " WHERE id = \"" + r.ID + "\""
		}
		break
	case "CREATE":

	}
	_, err := database.Db.Exec(query)
	if err != nil {
		panic(err)
	}
}

//a les droits administrateur ou demande de toucher quelque chose qui lui appartient
func isLegal(received database.ReceivedData, user database.User) bool {
	var answer bool
	if user.Role > 2 || user.IsAuthor(received.ID, received.Table) {
		switch received.Action {
		case "UPDATE":
			answer = canUpdate(user.Role, received, user)
			fmt.Println(answer)
		case "DELETE":
			answer = canDelete(user.Role, received, user)
		case "CREATE":
		}
	} else {
		answer = false
	}
	return answer
}

//Si Delete est quelque chose que les droits de l'user permet de modifier ou qu'il est administrateur
func canDelete(role int, r database.ReceivedData, user database.User) bool {
	var answer bool
	if (admin.TableExist(r.Table)) && (user.Role > 2 || toolbox.Contain(userCanDelete, r.Table)) {
		answer = true
	} else {
		answer = false
	}

	return answer
}

//Si update est quelque chose que les droits de l'user permet de modifier ou qu'il est administrateur
func canUpdate(role int, r database.ReceivedData, user database.User) bool {
	var answer bool
	if (admin.ColExist(r.Table, r.What)) && (user.Role > 2 || toolbox.Contain(userCanAlterate, r.Table)) {
		answer = true
	} else {
		answer = false
	}

	return answer
}
