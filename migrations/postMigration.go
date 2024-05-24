package migrations

import (
	"fmt"
	"forum/helpers"
	"forum/models"
)

func MigratePost() {
	// Réinitialisation de la table posts
	models.Post.Reset()

	users := models.User.Fetch()
	categories := models.Category.Fetch()

	var randCat = func() string {
		r1 := helpers.RandN(1, len(categories))
		r2 := helpers.RandN(1, len(categories))
		return fmt.Sprintf("%s,%s",categories[r1].Values["name"],categories[r2].Values["name"])
	}

	for i, user := range users {
		newPost := models.Post
		newPost.Values = map[string]interface{}{
			"authorId":        user.Values["id"],
			"title":           fmt.Sprintf("Post %d", i+1),
			"content":         helpers.DummyText,
			"creationDate":    helpers.Now(),
			"categoriesNames": randCat(),
		}
		newPost.Create()
	}
}
