package migrations

import (
	"forum/models"
)

func MigrateSession() {
	// Réinitialisation de la table sessions
	models.Session.Reset()
}
