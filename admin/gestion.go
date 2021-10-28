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
	User      []database.User
	Post      []database.Post
	Comment   []database.Comment
	Category  []database.Category
	Self      database.User
	CountLike []database.CountLike
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
		Data.Category = database.GetcategoryList()
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
	rows, _ := database.Db.Query("SELECT * FROM user CROSS JOIN (SELECT COUNT(*) AS Count FROM user) LEFT JOIN house ON house.id_house = user.house_id")

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

// Récupère un post depuis son ID :
func GetPostOnlyByID(ID int, Data Data) Data {
	var post database.Post
	id := ID
	row := database.Db.QueryRow("SELECT * FROM post INNER JOIN users ON post.author_id = users.id WHERE post.id = ?", id) // id, title, author_id, content, category_id, date, image, state
	//defer rows.Close()
	row.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Content, &post.CategoryID, &post.Date, &post.Image, &post.State, &post.Promoted, &post.Author.ID, &post.Author.Username, &post.Author.Password, &post.Author.Email, &post.Author.Role, &post.Author.Avatar, &post.Author.Date, &post.Author.State, &post.Author.SecretQuestion, &post.Author.SecretAnswer, &post.Author.House.ID)
	//author, _ := GetUserByID(post.AuthorID)
	//post.Author = author
	Data.Post = append(Data.Post, post)
	return Data
}

func GetcategoryList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM category CROSS JOIN (SELECT COUNT(*) AS Count FROM category)")
	defer rows.Close()
	for rows.Next() {
		var category database.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Theme, &category.Description, &category.Count)
		if err != nil {
			panic(err)
		}
		Data.Category = append(Data.Category, category)
	}
	return Data
}

//Récupère tout les commentaires
func GetCommentList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM comment CROSS JOIN (SELECT COUNT(*) AS Count FROM comment)")
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

//Récupère tout les post
func GetPostList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM post CROSS JOIN (SELECT COUNT(*) AS Count FROM post) INNER JOIN user ON post.author_id = user.id_user INNER JOIN category ON post.category_id = category.id_category")
	defer rows.Close()
	for rows.Next() {
		var newPost database.Post
		err := rows.Scan(&newPost.ID, &newPost.Title, &newPost.AuthorID, &newPost.Content, &newPost.CategoryID, &newPost.Date, &newPost.State, &newPost.Promoted, &newPost.Count, &newPost.Author.ID, &newPost.Author.Username, &newPost.Author.Password, &newPost.Author.Email, &newPost.Author.Role, &newPost.Author.Avatar, &newPost.Author.Date, &newPost.Author.State, &newPost.Author.SecretQuestion, &newPost.Author.SecretAnswer, &newPost.Author.House.ID, &newPost.Category.ID, &newPost.Category.Name, &newPost.Category.Theme, &newPost.Category.Description)
		if err != nil {
			panic(err)
		}
		Data.Post = append(Data.Post, newPost)
	}
	return Data
}

func Getlike(Data Data) Data {
	rows, _ := database.Db.Query("SELECT post_id, sum(case when type = 'like' then 1 else 0 end) as like, sum(case when type = 'dislike' then 1 else 0 end) as dislike FROM post_like GROUP BY post_id")
	defer rows.Close()
	for rows.Next() {
		var newCount database.CountLike
		err := rows.Scan(&newCount.PostId, &newCount.Countlike, &newCount.CountDislike)
		if err != nil {
			panic(err)
		}
		Data.CountLike = append(Data.CountLike, newCount)
	}
	return Data
}

//Compte les post par mois
func GetPostChart(Data DataForChart) DataForChart {
	rows, _ := database.Db.Query("SELECT COUNT(id_post) AS Count, strftime('%m', date_post) as Critere FROM post WHERE strftime('%Y', date_post) = '2021' GROUP BY strftime('%m', date_post)")
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

//Compte les post par category
func GetcategoryChart(Data DataForChart) DataForChart {
	rows, _ := database.Db.Query("SELECT COUNT(id_post) AS Count, category_id as Critere FROM post WHERE strftime('%Y', date_post) = '2021' GROUP BY category_id")
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
	rows, _ := database.Db.Query("SELECT COUNT(id_user) AS Count, strftime('%m', date_user) as Critere FROM user WHERE strftime('%Y', date_user) = '2021' GROUP BY strftime('%m', date_user)")
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
