package migrations

import (
	"forum/models"
)

func MigrateDiscussion() {
	models.Discussions.Reset()
}
