package models

// Déclaration des attributs de la classe
var User = Model{
	Table: "users",
	Columns: []string{
		"id",
		"name",
		"username",
		"email",
		"password",
		"picture",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"name TEXT",
		"username TEXT",
		"email TEXT",
		"password TEXT",
		"picture TEXT",
		"creationDate TEXT",
	},
}
