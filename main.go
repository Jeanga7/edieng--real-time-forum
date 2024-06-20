package main

import (
	"fmt"
	"forum/database"
	"forum/helpers"
	"forum/migrations"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	Routes(mux)

	port := ":8989"
	// Ouvrir la base de données
	database.OpenDB()

	args := os.Args
	if len(args) > 1 && args[1] == "--reset" {
		resetApp()
		// fmt.Println(database.FetchDataJoin([]string{"posts", "users"}, []string{"*"}, "posts.authorId=users.id")[0])
	} else if len(args) == 1 {
		// Réinitialisation des sessions
		migrations.MigrateSession()
	}

	// fmt.Println("Hello world !")
	defer database.DB.Close()

	fmt.Println("Server Online... : http://localhost" + port)
	err := http.ListenAndServe(port, mux)
	helpers.ErrorSimple(err)

}
