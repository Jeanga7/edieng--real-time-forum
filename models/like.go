package models

// Déclaration des attributs de la classe
var Like = Model{
	Table: "likes",
	Columns: []string{
		"id",
		"authorId",
		"likedId",
		"typeLiked",
		"creationDate",
	},
	Types: []string{
		"id TEXT PRIMARY KEY",
		"authorId TEXT",
		"likedId TEXT",
		"typeLiked TEXT",
		"creationDate TEXT",
	},
}
