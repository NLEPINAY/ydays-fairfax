package admin

import (
	"encoding/json"
	"forum/database"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	// "golang.org/x/crypto/bcrypt"
)

var profileTmpl = template.Must(template.ParseGlob("./templates/*"))

type Data struct {
	User     []database.User
	Post     []database.Post
	Comment  []database.Comment
	Category []database.Category
	Self     database.User
}

type DataForChart struct {
	/*DataUser     []database.User*/
	DataChart []database.CritereChart
	/*DataComment  []database.Comment
	DataCategory []database.Category*/
}

type ReceivedData struct {
	ID       string `json:"id"`
	Category string `json:"cat"`
	Value    string `json:"val"`
	NewValue string `json:"newVal"`
	Table    string `json:"table"`
	Reason   string `json:"reason"`
}

var results []string

var tmpl *template.Template

func Moderation(w http.ResponseWriter, r *http.Request, user database.User) {
	switch r.Method {
	case "POST":
		var p ReceivedData
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			ERROR, _ := json.Marshal("ERROR")
			w.Write(ERROR)
			panic(err)
		}
		//Si la demande est de promote un post alors on le promote ici
		if p.Table == "promote" {
			database.Db.Exec("insert or replace into promoted_post (post_id) values (?)", p.ID)
		} else {
			//Sinon c'est pour update un élément de la base de donnée.
			if isDatabaseTable(p.Table) && ColExist(p.Table, p.Category) {
				querry := "UPDATE " + p.Table + " SET " + p.Category
				_, err = database.Db.Exec(querry+` = ? WHERE id = ?`, p.NewValue, p.ID)
				if err != nil {
					panic(err)
				}
			}
		}
		var msg = "{\"message\": \"" + p.NewValue + "\"}"
		w.Write([]byte(msg))
	case "GET":
		//preparation de la bdd pour afficher les utilisateurs
		var Data Data
		Data = GetClientList(Data)
		Data = GetCommentList(Data)
		Data = GetPostList(Data)
		Data.Self = user
		Data.Category = database.GetCategoriesList()
		err := profileTmpl.ExecuteTemplate(w, "moderation", Data)
		if err != nil {
			panic(err)
		}
	}
}
func TableExist(table string) bool {
	query := "SELECT * FROM " + table
	cols, err := database.Db.Query(query)
	rows, err := cols.Columns()
	if err != nil {
		panic(err)
	}
	return len(rows) > 0
}
func ColExist(table string, cat string) bool {
	query := "SELECT * FROM " + table
	cols, err := database.Db.Query(query)
	rows, err := cols.Columns()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(rows); i++ {
		if rows[i] == cat {
			return true
		}
	}
	return false
}
func isDatabaseTable(table string) bool {
	var newtable string
	rows, err := database.Db.Query("SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&newtable)
		if newtable == table {
			return true
		}
	}
	return false
}

//Récupère tout les users
func GetClientList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM users CROSS JOIN (SELECT COUNT(*) AS Count FROM users) LEFT JOIN houses ON houses.id = users.house_id")

	defer rows.Close()
	for rows.Next() {
		var newUser database.User

		err := rows.Scan(&newUser.ID, &newUser.Username, &newUser.Password, &newUser.Email, &newUser.Role, &newUser.Avatar, &newUser.Date, &newUser.State, &newUser.SecretQuestion, &newUser.SecretAnswer, &newUser.House.ID, &newUser.Count, &newUser.House.ID, &newUser.House.Name, &newUser.House.Image)
		if err != nil {
			panic(err)
		}
		//newUser.House = database.GetHouseByID(newUser.House.ID)
		Data.User = append(Data.User, newUser)

	}

	return Data
}

//Récupère tout les commentaires
func GetCommentList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM comments CROSS JOIN (SELECT COUNT(*) AS Count FROM comments)")
	defer rows.Close()
	for rows.Next() {
		var newComment database.Comment
		err := rows.Scan(&newComment.ID, &newComment.AuthorID, &newComment.PostID, &newComment.Content, &newComment.Gif, &newComment.Date, &newComment.State, &newComment.Reason, &newComment.Count)
		if err != nil {
			panic(err)
		}
		Data.Comment = append(Data.Comment, newComment)
	}
	return Data
}

//Récupère tout les posts
func GetPostList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM posts CROSS JOIN (SELECT COUNT(*) AS Count FROM posts)")
	defer rows.Close()
	for rows.Next() {
		var newPost database.Post
		err := rows.Scan(&newPost.ID, &newPost.Title, &newPost.AuthorID, &newPost.Content, &newPost.CategoryID, &newPost.Date, &newPost.Image, &newPost.State, &newPost.Reason, &newPost.Count)
		if err != nil {
			panic(err)
		}
		Data.Post = append(Data.Post, newPost)
	}
	return Data
}

//Compte les posts par mois
func GetPostChart(Data DataForChart) DataForChart {
	rows, _ := database.Db.Query("SELECT COUNT(ID) AS Count, strftime('%m', Date) as Critere FROM posts WHERE strftime('%Y', Date) = '2021' GROUP BY strftime('%m', Date)")
	defer rows.Close()
	for rows.Next() {
		var newDataPost database.CritereChart
		err := rows.Scan(&newDataPost.Count, &newDataPost.Critere)
		if err != nil {
			panic(err)
		}
		Data.DataChart = append(Data.DataChart, newDataPost)
	}
	return Data
}

//Compte les posts par categorie
func GetCategoriesChart(Data DataForChart) DataForChart {
	rows, _ := database.Db.Query("SELECT COUNT(ID) AS Count, category_id as Critere FROM posts WHERE strftime('%Y', Date) = '2021' GROUP BY category_id")
	defer rows.Close()
	for rows.Next() {
		var newDataPost database.CritereChart
		err := rows.Scan(&newDataPost.Count, &newDataPost.Critere)
		if err != nil {
			panic(err)
		}
		Data.DataChart = append(Data.DataChart, newDataPost)
	}
	return Data
}

//Evolution du nombre des utilisateurs
func GetUserChart(Data DataForChart) DataForChart {
	rows, _ := database.Db.Query("SELECT COUNT(ID) AS Count, strftime('%m', Date) as Critere FROM users WHERE strftime('%Y', Date) = '2021' GROUP BY strftime('%m', Date)")
	defer rows.Close()
	for rows.Next() {
		var newDataUser database.CritereChart
		err := rows.Scan(&newDataUser.Count, &newDataUser.Critere)
		if err != nil {
			panic(err)
		}
		Data.DataChart = append(Data.DataChart, newDataUser)
	}
	return Data
}
