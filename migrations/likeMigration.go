package migrations

import (
	"forum/models"
)

func MigrateLike() {
	// Réinitialisation de la table likes
	models.Like.Reset()
}
