package models

// Déclaration des attributs de la classe
var Session = Model{
	Table: "sessions",
	Columns: []string{
		"id",
		"idUser",
		"creationDate",
		"expiryDate",
	},
	Types: []string{
		"id int PRIMARY KEY AUTOINCREMENT",
		"idUser TEXT",
		"creationDate TEXT",
		"expiryDate TEXT",
	},
}
