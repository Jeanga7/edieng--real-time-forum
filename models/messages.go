package models

// Déclaration des attributs de la classe
var Discussions = Model{
	Table: "discussions",
	Columns: []string{
		"id",
		"idDiscussingUsers", // [idEmma, idJean]
		"messages", // [{},{},]
	},
	Types: []string{
		"id int PRIMARY KEY AUTOINCREMENT",
		"idDiscussingUsers TEXT",
		"messages TEXT",
	},
}
