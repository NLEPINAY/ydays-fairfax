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

// Récupère un post depuis son ID :
func GetPostOnlyByID(ID int, Data Data) Data {
	var post database.Post
	id := ID
	row := database.Db.QueryRow("SELECT * FROM posts INNER JOIN users ON posts.author_id = users.id WHERE posts.id = ?", id) // id, title, author_id, content, category_id, date, image, state
	row.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Content, &post.CategoryID, &post.Date, &post.Image, &post.State, &post.Reason, &post.Author.ID, &post.Author.Username, &post.Author.Password, &post.Author.Email, &post.Author.Role, &post.Author.Avatar, &post.Author.Date, &post.Author.State, &post.Author.SecretQuestion, &post.Author.SecretAnswer, &post.Author.House.ID)
	Data.Post = append(Data.Post, post)
	return Data
}

func GetCommentByID(ID int, Data Data) Data {
	var comment database.Comment
	id := ID
	row := database.Db.QueryRow("SELECT * FROM comments INNER JOIN users ON comments.author_id = users.id WHERE comments.id = ?", id)
	row.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &comment.Gif, &comment.Date, &comment.State, &comment.Reason, &comment.Author.ID, &comment.Author.Username, &comment.Author.Password, &comment.Author.Email, &comment.Author.Role, &comment.Author.Avatar, &comment.Author.Date, &comment.Author.State, &comment.Author.SecretQuestion, &comment.Author.SecretAnswer, &comment.Author.House.ID)
	Data.Comment = append(Data.Comment, comment)
	return Data
}

func GetCategoriesList(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM categories CROSS JOIN (SELECT COUNT(*) AS Count FROM categories)")
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
	rows, _ := database.Db.Query("SELECT * FROM posts CROSS JOIN (SELECT COUNT(*) AS Count FROM posts) INNER JOIN users ON posts.author_id = users.id INNER JOIN categories ON posts.category_id = categories.id")
	defer rows.Close()
	for rows.Next() {
		var newPost database.Post
		err := rows.Scan(&newPost.ID, &newPost.Title, &newPost.AuthorID, &newPost.Content, &newPost.CategoryID, &newPost.Date, &newPost.Image, &newPost.State, &newPost.Reason, &newPost.Count, &newPost.Author.ID, &newPost.Author.Username, &newPost.Author.Password, &newPost.Author.Email, &newPost.Author.Role, &newPost.Author.Avatar, &newPost.Author.Date, &newPost.Author.State, &newPost.Author.SecretQuestion, &newPost.Author.SecretAnswer, &newPost.Author.House.ID, &newPost.Category.ID, &newPost.Category.Name, &newPost.Category.Theme, &newPost.Category.Description)
		if err != nil {
			panic(err)
		}
		Data.Post = append(Data.Post, newPost)
	}
	return Data
}

func GetLikes(Data Data) Data {
	rows, _ := database.Db.Query("SELECT post_id, sum(case when type = 'like' then 1 else 0 end) as likes, sum(case when type = 'dislike' then 1 else 0 end) as dislikes FROM post_likes GROUP BY post_id")
	defer rows.Close()
	for rows.Next() {
		var newCount database.CountLike
		err := rows.Scan(&newCount.PostId, &newCount.CountLikes, &newCount.CountDislikes)
		if err != nil {
			panic(err)
		}
		Data.CountLike = append(Data.CountLike, newCount)
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
