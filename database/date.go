package database

import (
	"database/sql"
	"strconv"
	"time"
)

//Remplie des ints correspondant a 4 périodes: a vie, du mois,de la semaine,des dernières 24heures. Nombre de post au cours de ces dates.
func GetNumberOfPostByDateAndPostCategory(cat int, life *int, month *int, week *int, day *int) {
	monthly := time.Now().AddDate(0, -1, 0)
	weekly := time.Now().AddDate(0, 0, -7)
	daily := time.Now().AddDate(0, 0, -1)
	rows, err := Db.Query("SELECT date FROM post p WHERE p.category_id = ?", cat)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		rows.Scan(&post.Date)
		*life++
		if post.Date.After(daily) {
			*month++
			*week++
			*day++
		} else if post.Date.After(weekly) {
			*month++
			*week++
		} else if post.Date.After(monthly) {
			*month++
		}
	}
}

//Remplie des ints correspondant a 4 périodes: a vie, du mois,de la semaine,des dernières 24heures. Nombre de commentaires au cours de ces dates.
func GetNumberOfCommentByDateAndPostCategory(cat int, life *int, month *int, week *int, day *int) {
	monthly := time.Now().AddDate(0, -1, 0)
	weekly := time.Now().AddDate(0, 0, -7)
	daily := time.Now().AddDate(0, 0, -1)
	rows, err := Db.Query("SELECT c.date FROM comment c INNER JOIN post p ON p.id = c.post_id WHERE p.category_id = ?", cat)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.Date)
		*life++
		if comment.Date.After(daily) {
			*month++
			*week++
			*day++
		} else if comment.Date.After(weekly) {
			*month++
			*week++
		} else if comment.Date.After(monthly) {
			*month++
		}
	}
}

//Remplie des ints correspondant a 4 périodes: a vie, du mois,de la semaine,des dernières 24heures. Nombre de like OU dislike au cours de ces dates. Reaction correspond a like ou dislike.
func GetNumberOfReactionByDate(cat int, reaction string, life *int, month *int, week *int, day *int) {
	monthly := time.Now().AddDate(0, -1, 0)
	weekly := time.Now().AddDate(0, 0, -7)
	daily := time.Now().AddDate(0, 0, -1)
	rows, err := Db.Query(`SELECT cl.date FROM comment_like cl INNER JOIN comment c ON c.id = cl.comment_id INNER JOIN post p ON p.id = c.post_id and p.category_id = ? WHERE cl.type = ? 
						   UNION ALL 
						   SELECT pl.date FROM post_like pl INNER JOIN post p ON p.id and p.id = pl.post_id and p.category_id = ? WHERE pl.type = ?`, cat, reaction, cat, reaction)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		//sa prend aussi les post like
		var comment CommentLike
		rows.Scan(&comment.Date)
		*life++
		if comment.Date.After(daily) {
			*month++
			*week++
			*day++
		} else if comment.Date.After(weekly) {
			*month++
			*week++
		} else if comment.Date.After(monthly) {
			*month++
		}
	}

}

//Renvoie le post le plus like (ne prend pas en compte les dislike) de la semaine et une erreur si une erreur parviens lors de l'appel de db
func GetMostLikedPostOfTheWeek() (Post, error) {
	var res Post
	var rows *sql.Rows
	var err error

	for day := 7; res.Title == ""; day = day + 7 {
		rows, err = Db.Query(`SELECT *,count(case post_id WHEN l.type = "like" then 1 else 0 end) AS amount FROM post p 
							INNER JOIN post_like l ON l.post_id = p.id_post
							WHERE p.date_post > datetime('now', '-` + strconv.Itoa(day) + ` day') AND p.state_post = 0
							GROUP BY l.post_id
							ORDER BY amount DESC
							LIMIT 1`)
		defer rows.Close()
		for rows.Next() {
			var postID int
			var userID int
			var myType string
			var date time.Time
			var amount int
			rows.Scan(&res.ID, &res.Title, &res.AuthorID, &res.Content, &res.CategoryID, &res.Date, &res.State, &res.Promoted, &postID, &userID, &myType, &date, &amount)
		}
	}

	return res, err
}

//Renvoie le post le plus commenter de la semaine et une erreur si une erreur parviens lors de l'appel de db
func GetMostCommentedPostOfTheWeek() (Post, error) {
	var res Post
	var rows *sql.Rows
	var err error

	for day := 7; res.Title == ""; day = day + 7 {
		rows, err = Db.Query(`SELECT *,count(post_id) AS amount FROM post p 
								INNER JOIN comment c ON c.post_id = p.id_post
								WHERE p.date_post > datetime('now', '-` + strconv.Itoa(day) + ` day') AND p.state_post = 0
								GROUP BY c.post_id
								ORDER BY amount DESC
								LIMIT 1`)
		defer rows.Close()
		for rows.Next() {
			var id int
			var authorID int
			var postID int
			var content string
			var gif string
			var date time.Time
			var state int
			var reason string
			var amount int
			rows.Scan(&res.ID, &res.Title, &res.AuthorID, &res.Content, &res.CategoryID, &res.Date, &res.State, &res.Promoted, &id, &authorID, &postID, &content, &gif, &date, &state, &reason, &amount)
		}
	}
	return res, err
}

//Renvoie le post le plus récent, une erreur est fournis si l'appel de Db plante
func GetMostRecentPost() (Post, error) {
	var res Post
	rows, err := Db.Query(`SELECT * FROM post p WHERE p.state_post = 0
	ORDER BY id_post DESC
	LIMIT 1 `)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Title, &res.AuthorID, &res.Content, &res.CategoryID, &res.Date, &res.State, &res.Promoted)
	}
	return res, err
}

//Renvoie le post promus par un modérateur ou administrateur
func GetPromotedPost() (Post, error) {
	var res Post
	rows, err := Db.Query(`SELECT * FROM post WHERE promoted=1 LIMIT 1`)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Title, &res.AuthorID, &res.Content, &res.CategoryID, &res.Date, &res.State, &res.Promoted)
	}
	return res, err
}
