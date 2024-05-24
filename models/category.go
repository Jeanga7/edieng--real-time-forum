package models

// Déclaration des attributs de la classe
var Category = Model{
	Table: "categories",
	Columns: []string{
		"id",
		"name",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"name TEXT",
		"creationDate TEXT",
	},
}
