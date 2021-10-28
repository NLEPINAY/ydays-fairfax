package database

import "log"

// Fonction récupératrice de like/dislike d'un post et du booléen 'liked/disliked par le user' en fonction de l'ID du post (et de l'ID du user) :
func GetlikeByPostID(ID int, currentUserID int) ([]PostLike, []PostLike, bool, bool) {
	var like, Dislike []PostLike
	var Liked, Disliked bool

	rows, err := Db.Query("SELECT * FROM post_like WHERE post_id = ?", ID)
	defer rows.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible de sélectionner les colonnes de la table post_like avec post_id = ", ID)
		panic(err)
	}

	for rows.Next() {
		var postLike PostLike
		rows.Scan(&postLike.PostID, &postLike.UserID, &postLike.Type, &postLike.Date)

		switch postLike.Type {
		case "like":
			like = append(like, postLike)
			if postLike.UserID == currentUserID {
				Liked = true
			}

		case "dislike":
			Dislike = append(Dislike, postLike)
			if postLike.UserID == currentUserID {
				Disliked = true
			}
		}
	}

	return like, Dislike, Liked, Disliked // Nombre de like/dislike, et booléen 'liké/disliké par l'utilisateur connecté'
}

// Fonction récupératrice de like/dislike d'un commentaire et du booléen 'liked/disliked par le user' en fonction de l'ID du commentaire (et de l'ID du user) :
func GetlikeByCommentID(ID int, currentUserID int) ([]CommentLike, []CommentLike, bool, bool) {
	var like, Dislike []CommentLike
	var Liked, Disliked bool

	rows, err := Db.Query("SELECT * FROM comment_like WHERE comment_id = ?", ID)
	defer rows.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible de sélectionner les colonnes de la table comment_like avec comment_id = ", ID)
		panic(err)
	}

	for rows.Next() {
		var commentLike CommentLike
		rows.Scan(&commentLike.CommentID, &commentLike.UserID, &commentLike.Type, &commentLike.Date)
		switch commentLike.Type {
		case "like":
			like = append(like, commentLike)
			if commentLike.UserID == currentUserID {
				Liked = true
				Disliked = false
			}

		case "dislike":
			Dislike = append(Dislike, commentLike)
			if commentLike.UserID == currentUserID {
				Disliked = true
				Liked = false
			}
		}
	}

	return like, Dislike, Liked, Disliked
}
