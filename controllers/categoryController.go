package controllers

import (
	"forum/global"
	"forum/models"
	"net/http"
)

func GetAllCategories(w http.ResponseWriter, r *http.Request) []map[string]interface{} {
	global.W = w
	data := models.Category.Fetch()
	var newData []map[string]interface{}
	for _, cat := range data {
		newData = append(newData, cat.Values)
	}
	return newData
}
