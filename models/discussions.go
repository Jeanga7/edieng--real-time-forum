package models

// Déclaration des attributs de la classe
var Discussions = Model{
	Table: "discussions",
	Columns: []string{
		"id",
		"idUsers",  // [idEmma, idJean]
		"messages", // [{},{},]
		"creationDate",
	},
	Types: []string{
		"id TEXT",
		"idUsers TEXT",
		"messages TEXT",
		"creationDate",
	},
}

/*

	var Discussions = Model {
	Table: "discussions",
	Columns: []string{
		"id",
		"idSender",
		"idReceiver",
		"messages",          // [{},{},]
		"creationDate",
	},
	Types: []string{
		"id Autoincrement",
		"idSender" TEXT,
		"idReceiver" TEXT,
		"messages TEXT",
		"creationDate",
	},
}

John Envoi à Alioune

id => 2

idSender = John

idReceiver = Alioune

messages = **********

creationDate

Pour fetcher les messages:

--> if idsender == celui qui est connecté && idreceiver == alioune  || idsender == alioune && ==celui qui est connecté et idreceiver

--> if idsender == celui qui est connecté mettre les messages à droite else mettre les messages à gauche??





*/
