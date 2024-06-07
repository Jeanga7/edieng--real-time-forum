package models

import (
	"errors"
	"fmt"
	"forum/database"
	"forum/helpers"
)

type Model struct {
	Table   string
	Columns []string
	Types   []string
	Values  map[string]interface{}
}

func (model *Model) Fetch(attr ...[]interface{}) []Model {
	// Cette fonction peut être appelée avec 2 arguments les éléments
	// d'une condition ([]columns, []values)
	// Exemple : model.Fetch([]string{"typeLiked", "creationDate"}, []interface{}{"post","2024-01-03"})
	// Lorsqu'elle est appelée sans argument, toutes les entrées de la
	// table du model sont récupérées
	// L'opérateur pour la récupération est "="
	models := []Model{}
	var data []map[string]interface{}

	// Vérifier si des attributs ont été fournis
	if len(attr) == 2 {
		// Si des attributs ont été fournis, effectuer la récupération avec les attributs
		data = database.FetchData(model.Table, attr[0], attr[1])
	} else if len(attr) == 0 {
		// Sinon, effectuer la récupération sans condition
		data = database.FetchData(model.Table)
	} else {
		helpers.ErrorSimple(errors.New("not enough arguments"))
	}

	if len(data) > 0 {
		for _, row := range data {
			newModel := *model
			newModel.Values = make(map[string]interface{})

			for i, col := range row {
				newModel.Values[i] = col
			}
			models = append(models, newModel)
		}
	}

	return models
}

func (model *Model) FetchLike(attr ...[]interface{}) []Model {
	// Cette fonction peut être appelée avec 2 arguments les éléments
	// d'une condition ([]columns, []values)
	// Exemple : model.Fetch([]string{"typeLiked", "creationDate"}, []interface{}{"post","2024-01-03"})
	// Lorsqu'elle est appelée sans argument, toutes les entrées de la
	// table du model sont récupérées
	// L'opérateur pour la récupération est "="
	models := []Model{}
	var data []map[string]interface{}

	// Vérifier si des attributs ont été fournis
	if len(attr) == 2 {
		// Si des attributs ont été fournis, effectuer la récupération avec les attributs
		data = database.FetchLike(model.Table, attr[0], attr[1])
	} else if len(attr) == 0 {
		// Sinon, effectuer la récupération sans condition
		data = database.FetchLike(model.Table)
	} else {
		helpers.ErrorSimple(errors.New("not enough arguments"))
	}

	if len(data) > 0 {
		for _, row := range data {
			newModel := *model
			newModel.Values = make(map[string]interface{})

			for i, col := range row {
				newModel.Values[i] = col
			}
			models = append(models, newModel)
		}
	}

	return models
}

// Création d'un nouvel utilisateur
func (model *Model) Create() Model {
	model.Values["id"] = helpers.UUID()
	model.Values["creationDate"] = helpers.Now()
	database.InsertData(model.Table, model.Values)
	return *model
}

// Mise à jour d'un utilisateur
func (model *Model) Update(data map[string]interface{}, condAttr string, condVal interface{}) {
	var cols = make([]string, 0)
	var values = make([]interface{}, 0)
	for col, val := range data {
		cols = append(cols, col)
		values = append(values, val)
	}
	database.UpdateData(model.Table, cols, values, condAttr, condVal)
}

// Suppression d'un utilisateur
func (model *Model) Delete(attr ...string) {
	// Cette fonction peut être appelée avec 3 arguments les éléments
	// d'une condition (column, operator, value)
	// Exemple : model.Delete("nickname", "=", "hulk")
	// Lorsqu'elle est appelée sans argument, toutes les entrées de la
	// table du model sont supprimées
	// Vérifier si des attributs ont été fournis
	if len(attr) == 3 {
		// Si des attributs ont été fournis, effectuer la récupération avec les attributs
		database.DeleteData(model.Table, string(attr[0]), string(attr[1]), string(attr[2]))
	} else if len(attr) == 0 && model.Values["id"] != nil {
		// Sinon, effectuer la suppression sans condition
		database.DeleteData(model.Table, "id", "=", fmt.Sprintf("%s", model.Values["id"]))
	} else if len(attr) == 0 && len(model.Values) == 0 {
		// Sinon, effectuer la suppression sans condition
		database.DeleteData(model.Table)
	} else {
		helpers.ErrorSimple(errors.New("not enough arguments"))
	}
}

// Réinitialisation de la table
func (model *Model) Reset() {
	// Suppression de la table categories
	database.DropTable(model.Table)

	// Creation de la table categories
	database.CreateTable(model.Table, model.Columns)
}
