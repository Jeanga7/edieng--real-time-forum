package migrations

import (
	"forum/models"
)

func MigrateComment() {
	// Réinitialisation de la table comments
	models.Comment.Reset()
}
