package migrations

import (
	"forum/models"
)

func MigrateNotification() {
	// Réinitialisation de la table likes
	models.Notification.Reset()
}
