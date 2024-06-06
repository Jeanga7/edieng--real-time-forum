package controllers

import (
	"encoding/json"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"
)

var (
	IdUser string
)

type ActiveUsersData struct {
	UserName string
	Actions  int
}

func Index(w http.ResponseWriter, r *http.Request) {
	global.W = w
	// IsLogged = false
	// return
	Usernames, Sessions, shouldReturn, _ := FindUsersOnline(r)
	if shouldReturn {
		return
	}

	All_users, Users := GetUsers(Sessions, strings.Title(Usernames))

	SortedUsers := SortUsers(All_users, Sessions, strings.Title(Usernames))

	Posts := models.Post.Fetch()
	Likes := models.Like.Fetch()
	Dislikes := models.Dislike.Fetch()
	Comments := models.Comment.Fetch()

	UserAllData := make(map[string]map[string]interface{})

	stats := GetUsersStats(Posts, Likes, Dislikes, Comments, Users, UserAllData)

	// Récupérer les posts depuis la base de données
	isLogged, res := Auth(r, w)
	posts := CheckPostsToShow(r, w)

	var userData []map[string]interface{}
	for _, u := range res {
		userData = append(userData, u.Values)
	}

	var data = make(map[string]interface{})

	if len(userData) > 0 {
		data = userData[0]
		if data["picture"] == "" {
			data["picture"] = "default-user.png"
		}
		data["numPosts"] = len(MyPosts)
		data["numLikes"] = len(LikedPost) + NumCommentLike
		data["numDislikes"] = len(DisLikedPost) + NumCommentDisLike
		data["numComments"] = CountComment
	}

	categoriesStats := GetCategoriesStats(posts)
	jsonCategoriesStats, err := json.Marshal(categoriesStats)
	helpers.ErrorSimple(err)

	usersStats := MostActiveUser(stats, UserAllData)

	data["usersStats"] = usersStats
	data["categoriesStats"] = string(jsonCategoriesStats)

	allData := helpers.FinalData{
		Posts:      posts,
		IsLogged:   isLogged,
		UserData:   data,
		AllUser:    SortedUsers,
		Categories: GetAllCategories(w, r),
	}

	if global.Status != 200 {
		ErrorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		tmpl, err := template.ParseFiles("./views/index.html", "./views/header.html")
		helpers.ErrorSimple(err)

		tmpl.Execute(w, allData)
		global.Status = 200
	}
}

func SortUsers(All_users []helpers.User, Sessions []models.Model, username string) []helpers.User {
	SortedUsers := make([]helpers.User, len(All_users))
	for i, user := range All_users {
		SortedUsers[i].Username = strings.Title(strings.Join(strings.Split(user.Username, " "), "_"))
		for _, s := range Sessions {
			var id = GetIdByUsername(user.Username)
			if id == s.Values["idUser"].(string) {
				SortedUsers[i].Status = "Online"
			}
		}
	}

	sort.Slice(SortedUsers, func(i, j int) bool {
		return SortedUsers[i].Username < SortedUsers[j].Username
	})

	filteredUsers := []helpers.User{}
	for _, user := range SortedUsers {
		if user.Username != strings.Title(strings.Join(strings.Split(username, " "), "_")) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func GetIdByUsername(username string) string {
	conn := models.User.Fetch([]interface{}{"name"}, []interface{}{username})
	if len(conn) > 0 {
		return conn[0].Values["id"].(string)
	}
	return ""
}

func FindUsersOnline(r *http.Request) (string, []models.Model, bool, interface{}) {
	var Usernames string
	var userID interface{}
	Sessions := models.Session.Fetch()
	Cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Println(err)
	} else {
		sessions := models.Session.Fetch([]interface{}{"Id"}, []interface{}{Cookie.Value})
		if len(sessions) > 0 {
			var ok bool
			userID, ok = sessions[0].Values["idUser"]
			if !ok {
				log.Println("seesion introuvable")
				return "", nil, true, "noId"
			}
			Usernames = GetUsernameById(userID)
		}
	}
	return Usernames, Sessions, false, userID
}

func GetUsernameById(userID interface{}) string {
	var Usernames string
	conn := models.User.Fetch([]interface{}{"Id"}, []interface{}{userID})
	if len(conn) > 0 {
		Usernames = conn[0].Values["name"].(string)
	}
	return Usernames
}

func GetUsers(Sessions []models.Model, Usernames string) ([]helpers.User, []models.Model) {
	var All_users []helpers.User
	Users := models.User.Fetch()
	for _, v := range Users {
		var user helpers.User
		user.Username = v.Values["name"].(string)

		if user.Username != Usernames {
			All_users = append(All_users, user)
		}
	}
	return All_users, Users
}

func MostActiveUser(stats map[string]int, UserAllData map[string]map[string]interface{}) []ActiveUsersData {
	var usersData []ActiveUsersData

	for id, stat := range stats {
		if userData, ok := UserAllData[id]; ok && userData != nil {
			usersData = append(usersData, ActiveUsersData{
				UserName: userData["name"].(string),
				Actions:  stat,
			})
		}
	}

	for i := 0; i < len(usersData); i++ {
		for j := i + 1; j < len(usersData)-1; j++ {
			if usersData[i].Actions < usersData[j].Actions {
				usersData[i], usersData[j] = usersData[j], usersData[i]
			}
		}
	}

	// if len(usersData) > 5 {
	// 	usersData = usersData[:5]
	// }
	return usersData
}

func GetCategoriesStats(posts []helpers.ForumPost) map[string]int {
	var categoriesStats = make(map[string]int)
	for _, post := range posts {
		categories := strings.Split(post.Categories, ",")
		for _, v := range categories {
			categoriesStats[v]++
		}
	}
	return categoriesStats
}

func CheckPostsToShow(r *http.Request, w http.ResponseWriter) []helpers.ForumPost {
	global.W = w
	var posts []helpers.ForumPost

	if params := r.URL.Query(); len(params.Get("filter")) > 0 {
		if params.Get("filter") == "liked" {
			posts = LikedPost
		} else if params.Get("filter") == "disliked" {
			posts = DisLikedPost
		} else if params.Get("filter") == "MyPosts" {
			posts = MyPosts
		} else if params.Get("filter") == "MyComments" {
			posts = MyComments
		}
	} else {
		posts = ReadPost(w, r)
	}
	return posts
}

func GetUsersStats(Posts []models.Model, Likes []models.Model, Dislikes []models.Model, Comments []models.Model, Users []models.Model, UserAllData map[string]map[string]interface{}) map[string]int {
	var stats = make(map[string]int)

	for _, v := range Posts {
		stats[v.Values["authorId"].(string)]++
	}
	for _, v := range Likes {
		stats[v.Values["authorId"].(string)]++
	}
	for _, v := range Dislikes {
		stats[v.Values["authorId"].(string)]++
	}
	for _, v := range Comments {
		stats[v.Values["authorId"].(string)]++
	}
	for _, v := range Users {
		// fmt.Println("this is my name", v.Values["name"])
		UserAllData[v.Values["id"].(string)] = v.Values
	}

	return stats
}

func ErrorPage(w http.ResponseWriter, name string, code int) {
	tmpl, err := template.ParseFiles("./views/error.html", "./views/header.html")
	helpers.ErrorSimple(err)

	data := map[string]interface{}{
		"errorName": name,
		"errorCode": code,
	}
	tmpl.Execute(w, data)
	global.Status = 200
}
