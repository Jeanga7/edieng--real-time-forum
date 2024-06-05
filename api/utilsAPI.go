package api

import (
	"encoding/json"
	"forum/helpers"
	"forum/models"
	"net/http"
)

func TransformModels(models []models.Model, fields []string) []map[string]interface{} {
	allValues := []map[string]interface{}{}
	for _, model := range models {
		data := make(map[string]interface{})
		for _, field := range fields {
			if value, ok := model.Values[field]; ok {
				data[field] = value
			}
		}
		allValues = append(allValues, data)
	}
	return allValues
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	helpers.ErrorSimple(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
