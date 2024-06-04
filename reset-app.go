package main

import "forum/migrations"

func resetApp() {
	// Réinitialisation des users
	migrations.MigrateUser()

	// Réinitialisation des categories
	migrations.MigrateCategory()

	// Réinitialisation des posts
	migrations.MigratePost()

	// Réinitialisation des comments
	migrations.MigrateComment()

	// Réinitialisation des likes
	migrations.MigrateLike()

	// Réinitialisation des dislikes
	migrations.MigrateDislike()

	// Réinitialisation des sessions
	migrations.MigrateSession()

	// Réinitialisation des Messages
	migrations.MigrateDiscussion()

	// Réinitialisation des Notifications
	migrations.MigrateNotification()
}
