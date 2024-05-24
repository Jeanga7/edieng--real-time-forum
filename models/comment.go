package models

// Déclaration des attributs de la classe
var Comment = Model{
	Table: "comments",
	Columns: []string{
		"id",
		"authorId",
		"postId",
		"content",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"authorId TEXT",
		"postId TEXT",
		"content TEXT",
		"creationDate TEXT",
	},
}
