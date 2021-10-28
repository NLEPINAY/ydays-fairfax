package database

func getChat(user User, secondUser string) Messagerie {
	var msg Messagerie
	var chat Chat
	rows, err := Db.Query("SELECT name FROM sqlite_master WHERE type ='table' AND name NOT LIKE 'sqlite_%';")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&chat)

	}
	return msg
}
