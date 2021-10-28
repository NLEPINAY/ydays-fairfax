package database

func GetBadgeByUserID(id int) []Badge {
	var result []Badge
	rows, err := Db.Query("SELECT id_badge,type,image_badge FROM user_badge u INNER JOIN badge ON badge.id_badge = u.badge_id WHERE user_id = $1", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var temp Badge
		rows.Scan(&temp.ID, &temp.Type, &temp.Image)
		result = append(result, temp)
	}
	return result
}
