package models

// Déclaration des attributs de la classe
var Discussions = Model {
	Table: "discussions",
	Columns: []string{
		"id",
		"idUsers", // [idEmma, idJean]
		"messages",          // [{},{},]
		"creationDate",
	},
	Types: []string{
		"id TEXT",
		"idUsers TEXT",
		"messages TEXT",
		"creationDate",
	},
}
