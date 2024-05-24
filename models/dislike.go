package models

// Déclaration des attributs de la classe
var Dislike = Model{
	Table: "dislikes",
	Columns: []string{
		"id",
		"authorId",
		"dislikedId",
		"typeDisliked",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"authorId TEXT",
		"dislikedId TEXT",
		"typeDisliked TEXT",
		"creationDate TEXT",
	},
}
