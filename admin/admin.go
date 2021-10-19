package admin

import "forum/database"

//Récupère tout les posts

//Récupère tout les posts
func GetPost(Data Data) Data {
	rows, _ := database.Db.Query("SELECT * FROM posts")
	defer rows.Close()
	for rows.Next() {
		var newPost database.Post
		err := rows.Scan(&newPost.ID, &newPost.Title, &newPost.AuthorID, &newPost.Content, &newPost.CategoryID, &newPost.Date, &newPost.Image, &newPost.State, &newPost.Reason)
		if err != nil {
			panic(err)
		}
		Data.Post = append(Data.Post, newPost)
	}
	return Data
}
