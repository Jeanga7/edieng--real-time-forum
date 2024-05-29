package migrations

import (
	"forum/models"
)

func MigrateMessage() {
	models.Discussions.Reset()
}
