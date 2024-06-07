package migrations

import (
	"forum/helpers"
	"forum/models"
)

func MigrateUser() {
	// Réinitialisation de la table users
	models.User.Reset()

	// Ajout d'utilisateurs fictifs
	var usersValues = make([]map[string]interface{}, 0)

	usersValues = append(usersValues, map[string]interface{}{
		"firstname": "Tony",
		"lastname":  "Stark",
		"nickname":  "ironman",
		"email":     "stark@avenger.com",
		"age":       47,
		"gender":    "M",
		"password":  helpers.HashPassword("azerty1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"firstname": "Steven",
		"lastname":  "Strange",
		"nickname":  "drstrange",
		"email":     "strange@avenger.com",
		"age":       42,
		"gender":    "M",
		"password":  helpers.HashPassword("azerty1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"firstname": "Bruce",
		"lastname":  "Banner",
		"nickname":  "hulk",
		"email":     "banner@avenger.com",
		"age":       45,
		"gender":    "M",
		"password":  helpers.HashPassword("azerty1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"firstname": "Natacha",
		"lastname":  "Romanov",
		"nickname":  "blackwidow",
		"email":     "romanov@avenger.com",
		"age":       32,
		"gender":    "F",
		"password":  helpers.HashPassword("azerty1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"firstname": "Wanda",
		"lastname":  "Maximoff",
		"nickname":  "scarletwitch",
		"email":     "wanda@avenger.com",
		"age":       28,
		"gender":    "F",
		"password":  helpers.HashPassword("azerty1234"),
	})

	for _, val := range usersValues {
		var user = models.User
		user.Values = val
		user.Create()
	}
}
