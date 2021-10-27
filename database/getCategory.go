package database

import "log"

func GetCategoryByID(id int) (Category, error) {
	var myCategory Category

	row := Db.QueryRow("SELECT * FROM category WHERE id_category = ?", id)
	row.Scan(&myCategory.ID, &myCategory.Name, &myCategory.Theme, &myCategory.Description)

	return myCategory, nil
}

func GetcategoryList() []Category {
	rows, err := Db.Query("SELECT * FROM category CROSS JOIN (SELECT COUNT(*) AS Count FROM category)")
	defer rows.Close()
	if err != nil {
		log.Println("❌ DATABASE | ERREUR : Impossible de récupérer la liste des catégories.")
		panic(err)
	}

	var categorys []Category
	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name, &category.Theme, &category.Description, &category.Count)
		categorys = append(categorys, category)
	}
	return categorys
}
