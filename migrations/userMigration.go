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
		"Name":     "Tony Stark",
		"Username": "ironman",
		"Email":    "stark@avenger.com",
		"Password": helpers.HashPassword("1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"Name":     "Steven Strange",
		"Username": "drstrange",
		"Email":    "strange@avenger.com",
		"Password": helpers.HashPassword("1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"Name":     "Bruce Banner",
		"Username": "hulk",
		"Email":    "banner@avenger.com",
		"Password": helpers.HashPassword("1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"Name":     "Natacha Romanov",
		"Username": "blackwidow",
		"Email":    "romanov@avenger.com",
		"Password": helpers.HashPassword("1234"),
	})

	usersValues = append(usersValues, map[string]interface{}{
		"Name":     "Wanda Maximoff",
		"Username": "scarletwitch",
		"Email":    "wanda@avenger.com",
		"Password": helpers.HashPassword("1234"),
	})

	for _, val := range usersValues {
		var user = models.User
		user.Values = val
		user.Create()
	}
}
