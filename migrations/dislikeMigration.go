package migrations

import (
	"forum/models"
)

func MigrateDislike() {
	// Réinitialisation de la table dislikes
	models.Dislike.Reset()
}
