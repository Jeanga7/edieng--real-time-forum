package migrations

import (
	"forum/models"
)

func MigrateCategory() {
	// Réinitialisation de la table categories
	models.Category.Reset()

	// Ajout de catégories fictives
	var categoriesValues = make([]map[string]interface{}, 0)

	categoriesValues = append(categoriesValues, map[string]interface{}{
		"Name": "Politics",
	})

	categoriesValues = append(categoriesValues, map[string]interface{}{
		"Name": "Mangas",
	})

	categoriesValues = append(categoriesValues, map[string]interface{}{
		"Name": "Movies",
	})

	categoriesValues = append(categoriesValues, map[string]interface{}{
		"Name": "Sport",
	})

	for _, val := range categoriesValues {
		var category = models.Category
		category.Values = val
		category.Create()
	}
}
