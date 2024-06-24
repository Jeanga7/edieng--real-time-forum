package models

// Déclaration des attributs de la classe
var Notification = Model{
	Table: "notification",
	Columns: []string{
		"id",
		"idSender",
		"idReceiver",
		"notifications",
		"creationDate",
	},
	Types: []string{
		"id",
		"idSender TEXT",
		"idReceiver TEXT",
		"notifications INTEGER",
		"creationDate",
	},
}
