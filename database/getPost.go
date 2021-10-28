package database

import (
	"log"
	"strconv"
)

// Récupère TOUS les post appartennant à la catégorie dont l'ID est passé en argument :
func GetpostByCategoryID(id int) ([]Post, error) {
	var posts []Post
	rows, err := Db.Query("SELECT * FROM post WHERE category_id = ? ORDER BY id_post DESC", id) // id, title, author_id, content, category_id, date, image, state

	defer rows.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible de récupérer les post de la catégorie dont l'ID est ", id)
		return posts, err
	}

	for rows.Next() {
		var post Post
		rows.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Content, &post.CategoryID, &post.Date, &post.State, &post.Promoted)

		author, _ := GetUserByID(post.AuthorID)
		post.Author = author

		comment, _ := GetcommentByPostID(post.ID, 0)
		post.Comment = comment
		post.Like, post.Dislike, post.Liked, post.Disliked = GetlikeByPostID(post.ID, 0)
		posts = append(posts, post)
	}
	return posts, nil
}

// Récupère un post depuis son ID :
func GetPostByID(ID int, currentUserID int) (Post, error) {
	var post Post

	row := Db.QueryRow("SELECT * FROM post WHERE id_post = ?", ID) // id, title, author_id, content, category_id, date, image, state
	row.Scan(&post.ID, &post.Title, &post.AuthorID, &post.Content, &post.CategoryID, &post.Date, &post.State, &post.Promoted)
	author, _ := GetUserByID(post.AuthorID)
	post.Author = author
	post.Like, post.Dislike, post.Liked, post.Disliked = GetlikeByPostID(post.ID, currentUserID)
	return post, nil
}

// Récupère tous les post likés par un utilisateur dont l'ID est passé en paramètre :
func GetpostLikedByUser(userID int) ([]Post, error) {
	var posts []Post

	rows, err := Db.Query("SELECT post_id FROM post_like WHERE user_id = ?", userID)
	defer rows.Close()
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post Post
		rows.Scan(&post.ID)
		post, _ = GetPostByID(post.ID, userID)
		posts = append(posts, post)
	}
	return posts, nil
}

// VIRGIL :
func GetpostFromUserByID(identifier int) ([]Post, error) {
	var posts []Post
	inject := "SELECT * FROM post WHERE author_id = " + strconv.Itoa(identifier)
	rows, _ := Db.Query(inject)
	defer rows.Close()
	for rows.Next() {
		var newPost Post
		rows.Scan(&newPost.ID, &newPost.Title, &newPost.AuthorID, &newPost.Content, &newPost.CategoryID, &newPost.Date, &newPost.State, &newPost.Promoted)
		author, _ := GetUserByID(newPost.AuthorID)
		newPost.Author = author
		posts = append(posts, newPost)
	}

	return posts, nil

	// Si le username ou l'email n'existe pas, user.ID == 0 (car par défaut, variable de type int = 0)
}
