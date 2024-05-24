package database

import (
	"database/sql"
	"fmt"
	"forum/helpers"
	"log"
	"strings"
)

var DB *sql.DB

func OpenDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/forum.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTable(table string, columns []string) {
	// Créer des entrées dans la base de données
	req, err := DB.Prepare("CREATE TABLE IF NOT EXISTS " + table + " (" + strings.Join(columns, ", ") + ")")
	helpers.ErrorSimple(err)

	defer req.Close()

	_, err = req.Exec()
	helpers.ErrorSimple(err)

}

func DropTable(table string) { // Supprimer la table et la créer à nouveau
	req, err := DB.Prepare("DROP TABLE IF EXISTS " + table)
	helpers.ErrorSimple(err)

	defer req.Close()

	_, err = req.Exec()
	helpers.ErrorSimple(err)

}

func InsertData(table string, data map[string]interface{}) {
	// Préparer les listes des colonnes et des valeurs à insérer
	var columns []string
	var values []interface{}

	// Parcourir les données et construire les listes des colonnes et des valeurs
	for key, value := range data {
		columns = append(columns, key)
		values = append(values, value)
	}

	// Préparer la chaîne de placeholder avec des points d'interrogation
	placeholders := strings.Repeat("?,", len(columns))
	placeholders = strings.TrimRight(placeholders, ",") // Supprimer la virgule finale

	// Préparer la requête SQL
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), placeholders)

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	// Exécuter la requête avec les valeurs
	_, err = req.Exec(values...)
	helpers.ErrorSimple(err)

}

func FetchData(table string, attrs ...[]interface{}) []map[string]interface{} {
	// Construire la requête SQL en fonction des colonnes et des valeurs fournies
	query := "SELECT * FROM " + table
	var args []interface{}

	// Vérifier si des colonnes et des valeurs ont été fournies
	if len(attrs) > 0 {
		columns := attrs[0]
		values := attrs[1]

		if len(columns) > 0 && len(columns) == len(values) {
			query += " WHERE "
			for i, column := range columns {
				query += fmt.Sprintf("%s=?", column)
				args = append(args, values[i])

				if i < len(columns)-1 {
					query += " AND "
				}
			}
		}
	}

	query += " ORDER BY creationDate DESC"

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	var rows *sql.Rows
	if len(args) > 0 {
		rows, err = req.Query(args...)
	} else {
		rows, err = req.Query()
	}
	helpers.ErrorSimple(err)

	defer rows.Close()

	// Extraire les lignes SQL et les retourner
	return helpers.ExtractSQLRows(rows)
}

func FetchLike(table string, attrs ...[]interface{}) []map[string]interface{} {
	// Construire la requête SQL en fonction des colonnes et des valeurs fournies
	query := "SELECT * FROM " + table
	var args []interface{}

	// Vérifier si des colonnes et des valeurs ont été fournies
	if len(attrs) > 0 {
		columns := attrs[0]
		values := attrs[1]

		if len(columns) > 0 && len(columns) == len(values) {
			query += " WHERE "
			for i, column := range columns {
				// Utiliser LOWER() pour rendre la recherche insensible à la casse
				query += fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", column)
				args = append(args, values[i])

				if i < len(columns)-1 {
					query += " AND "
				}
			}
		}
	}

	query += " ORDER BY creationDate DESC"

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	var rows *sql.Rows
	if len(args) > 0 {
		rows, err = req.Query(args...)
	} else {
		rows, err = req.Query()
	}
	helpers.ErrorSimple(err)

	defer rows.Close()

	// Extraire les lignes SQL et les retourner
	return helpers.ExtractSQLRows(rows)
}

func FetchDataJoin(tables []string, columns []string, condition string, values ...interface{}) []map[string]interface{} {
	// Exemple : FetchDataJoin(["posts", "users"], ["*"], "posts.authorId=users.id")
	// Construire la requête SQL
	query := "SELECT " + strings.Join(columns, ", ") + " FROM " + strings.Join(tables, ", ")

	// Ajouter une clause WHERE si une condition est fournie
	if condition != "" {
		query += " WHERE " + condition
	}

	query += " ORDER BY creationDate DESC"

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	var rows *sql.Rows
	if len(values) > 0 {
		rows, err = req.Query(values...)
	} else {
		rows, err = req.Query()
	}
	helpers.ErrorSimple(err)

	defer rows.Close()

	// Extraire les lignes SQL et les retourner
	return helpers.ExtractSQLRows(rows)
}

func DeleteData(table string, attrs ...string) {
	// Cette fonction peut être appelée avec 3 arguments les éléments
	// d'une condition (column, operator, value)
	// Supprimer des entrées de la base de données
	var query string

	// Construire la requête SQL en fonction de la présence des attributs et de leurs valeurs
	if len(attrs) == 3 {
		query = "DELETE FROM " + table + " WHERE " + attrs[0] + attrs[1] + "?"
	} else {
		query = "DELETE FROM " + table
	}

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	// Exécuter la requête
	var result sql.Result
	if len(attrs) == 3 {
		result, err = req.Exec(attrs[2])
	} else {
		result, err = req.Exec()
	}
	helpers.ErrorSimple(err)

	// Récupérer le nombre de lignes affectées, si nécessaire
	_, err = result.RowsAffected()
	helpers.ErrorSimple(err)

}

func UpdateData(table string, attributes []string, data []interface{}, conditionAttribute string, conditionValue interface{}) {
	// Préparer la chaîne de placeholders pour les attributs à mettre à jour
	setClause := ""
	for i, attr := range attributes {
		if i > 0 {
			setClause += ", "
		}
		setClause += attr + "=?"
	}

	// Préparer la requête SQL de mise à jour
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", table, setClause, conditionAttribute)

	// Préparer et exécuter la requête
	req, err := DB.Prepare(query)
	helpers.ErrorSimple(err)

	defer req.Close()

	// Concaténer les données de condition avec les autres données
	newData := append(data, conditionValue)

	// Exécuter la requête
	_, err = req.Exec(newData...)
	helpers.ErrorSimple(err)

}
