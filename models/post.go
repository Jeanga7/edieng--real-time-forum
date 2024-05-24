package models

// Déclaration des attributs de la classe
var Post = Model{
	Table: "posts",
	Columns: []string{
		"id",
		"authorId",
		"title",
		"content",
		"categoriesNames",
		"picture",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"authorId TEXT",
		"title TEXT",
		"content TEXT",
		"categoriesNames TEXT",
		"picture TEXT",
		"creationDate TEXT",
	},
}
