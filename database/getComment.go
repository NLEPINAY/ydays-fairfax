package database

import (
	"log"
)

// Fonction récupérant TOUS les commentaires du post dont l'ID est passé en argument :

func GetcommentByPostID(ID int, currentUserID int) ([]Comment, error) {
	var comments []Comment

	rows, err := Db.Query("SELECT * FROM comment WHERE post_id = ?", ID) // id, author_id, post_id, content, date, state
	defer rows.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible de récupérer les commentaires du post dont l'ID est ", ID)
		return comments, err
	}

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &comment.Gif, &comment.Date, &comment.State)
		post, _ := GetPostByID(comment.PostID, 0)

		comment.PostTitle = post.Title
		comment.posttate = post.State
		comment.Author, _ = GetUserByID(comment.AuthorID)
		comment.like, comment.Dislike, comment.Liked, comment.Disliked = GetlikeByCommentID(comment.ID, currentUserID)
		comments = append(comments, comment)
	}

	return comments, nil
}

func GetCommentByID(ID int, userID int) (Comment, error) {
	var comment Comment

	row := Db.QueryRow("SELECT * FROM comment WHERE id = ?", ID)
	row.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.Content, &comment.Gif, &comment.Date, &comment.State, &comment.Reason)

	author, _ := GetUserByID(comment.AuthorID)
	comment.Author = author

	post, _ := GetPostByID(comment.PostID, 0)
	comment.PostTitle = post.Title
	comment.posttate = post.State
	comment.like, comment.Dislike, comment.Liked, comment.Disliked = GetlikeByCommentID(comment.ID, userID)

	return comment, nil
}

// Récupère tous les commentaires likés par un utilisateur dont l'ID est passé en paramètre :
func GetcommentLikedByUser(userID int) ([]Comment, error) {
	var comments []Comment

	rows, err := Db.Query("SELECT comment_id FROM comment_like WHERE user_id = ?", userID)
	defer rows.Close()
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.ID)
		comment, _ = GetCommentByID(comment.ID, userID)
		comments = append(comments, comment)
	}
	return comments, nil
}

// Récupère tous les commentaires d'un utilisateur dont l'ID est passé en paramètre :
func GetCommentFromUserByID(userID int) ([]Comment, error) {
	var comments []Comment

	rows, err := Db.Query("SELECT id FROM comment WHERE author_id = ?", userID)
	defer rows.Close()

	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.ID)
		comment, _ = GetCommentByID(comment.ID, 0)
		comments = append(comments, comment)
	}

	return comments, nil
}
