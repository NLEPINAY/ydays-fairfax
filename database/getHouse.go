package database

func GetHouseByID(houseID int) House {
	var house House

	row := Db.QueryRow("SELECT * FROM house WHERE id_house = ?", houseID) // id, name, image
	row.Scan(&house.ID, &house.Name, &house.Image)

	return house
}
