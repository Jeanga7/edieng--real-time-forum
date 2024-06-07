package main

import (
	"fmt"
	"forum/database"
	"forum/helpers"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	Routes(mux)

	port := ":8989"
	// Ouvrir la base de données
	database.OpenDB()
	if args := os.Args; len(args) > 1 && args[1] == "--reset" {
		resetApp()
	}
	defer database.DB.Close()

	fmt.Println("Server Online... : http://localhost" + port)
	err := http.ListenAndServe(port, mux)
	helpers.ErrorSimple(err)

}
