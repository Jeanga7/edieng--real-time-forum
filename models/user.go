package models

// Déclaration des attributs de la classe
var User = Model{
	Table: "users",
	Columns: []string{
		"id",
		"firstname",
		"lastname",
		"nickname",
		"email",
		"password",
		"age",
		"gender",
		"picture",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"firstname TEXT",
		"lastname TEXT",
		"nickname TEXT",
		"email TEXT",
		"password TEXT",
		"age TEXT",
		"gender TEXT",
		"picture TEXT",
		"creationDate TEXT",
	},
}
